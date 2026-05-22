package register

import (
	"context"
	"encoding/json"
	"log"
	"parallel/internal/dto"
	"parallel/internal/event"
	luascripting "parallel/internal/lua-scripting"
	"parallel/internal/service"
	"parallel/internal/worker"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

// RegisterPayload is the expected event payload for a register request.

// ResponseWorker processes a single register-response event.
type RegisterWorker struct {
	client *redis.Client
	w      worker.Worker
}

func newRegisterWorker(id int32, client *redis.Client) *RegisterWorker {
	rw := &RegisterWorker{client: client}
	rw.w.ID = id
	rw.w.Job = "Response"
	rw.w.State.Store(true)
	return rw
}

func (r *RegisterWorker) handle(ctx context.Context, e event.Event) {
	var p dto.RegisterPayload
	if err := json.Unmarshal(e.Request.Payload, &p); err != nil {
		log.Printf("[ResponseWorker %d] bad payload: %v", r.w.ID, err)
		return
	}

	result, err := luascripting.RegisterScript.Run(ctx, r.client, []string{p.ClassID}).Int()
	if err != nil {
		log.Printf("[ResponseWorker %d] lua error: %v", r.w.ID, err)
		return
	}

	log.Printf("[ResponseWorker %d] class=%s student=%s status=%d", r.w.ID, p.ClassID, p.StudentID, result)
}

// Bus
// ResponseBus manages a pool of ResponseWorkers consuming from a shared queue.
type RegisterBus struct {
	queue       <-chan event.Event
	client      *redis.Client
	workerCount atomic.Int32
	minWorker   int
	maxWorker   int
	wg          sync.WaitGroup
}

func NewRegisterBus(queue <-chan event.Event, client *redis.Client) *RegisterBus {
	return &RegisterBus{queue: queue, client: client, minWorker: 10, maxWorker: 1000}
}

// Start launches n workers and blocks until ctx is cancelled, then drains remaining jobs.
func (b *RegisterBus) Start(ctx context.Context, n int) {
	for i := 0; i < n; i++ {
		b.spawnWorker(ctx)
	}
	// //monitor function
	// go func() {
	// 	defer
	// }

	//monitor, scale worker is needed?
}

func (b *RegisterBus) spawnWorker(ctx context.Context) {
	if b.workerCount.Load() >= int32(b.maxWorker) {
		return
	}
	id := b.workerCount.Add(1)
	rw := newRegisterWorker(id, b.client)
	b.wg.Add(1)
	service.TrackGo(+1)

	go func() {
		defer b.wg.Done()
		defer service.TrackGo(-1)
		defer b.workerCount.Add(-1)
		defer rw.w.State.Store(false)
		idleTimer := time.NewTimer(30 * time.Second)
		defer idleTimer.Stop()
		for {
			select {
			case e, ok := <-b.queue:
				if !ok {
					return
				}
				if !idleTimer.Stop() {
					select {
					case <-idleTimer.C:
					default:
					}
				}
				idleTimer.Reset(30 * time.Second)
				rw.handle(ctx, e)
			case <-idleTimer.C:
				if b.workerCount.Load() > int32(b.minWorker) {
					log.Printf("[RegisterWorker %d] idle timeout, scaling down", rw.w.ID)
					return
				}
				idleTimer.Reset(30 * time.Second)
			case <-ctx.Done():
				// drain remaining jobs before exit
				for {
					select {
					case e, ok := <-b.queue:
						if !ok {
							return
						}
						rw.handle(context.Background(), e)
					default:
						return
					}
				}
			}
		}
	}()
}

// Wait blocks until all workers have finished.
func (b *RegisterBus) Wait() { b.wg.Wait() }
