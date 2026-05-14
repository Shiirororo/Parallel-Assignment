package responseworker

import (
	"context"
	"encoding/json"
	"log"
	"parallel/internal/event"
	luascripting "parallel/internal/lua-scripting"
	"parallel/internal/worker"
	"sync"
	"sync/atomic"

	"github.com/redis/go-redis/v9"
)

// RegisterPayload is the expected event payload for a register request.
type RegisterPayload struct {
	ClassID    string `json:"class_id"`
	StudentID  string `json:"student_id"`  //Need to be array
	ResponseCh string `json:"response_ch"` // identifier / correlation ID for the response
}


// ResponseWorker processes a single register-response event.
type RegisterWorker struct {
	client *redis.Client
	w      worker.Worker
}

func newResponseWorker(id int32, client *redis.Client) *RegisterWorker {
	rw := &RegisterWorker{client: client}
	rw.w.ID = id
	rw.w.Job = "Response"
	rw.w.State.Store(true)
	return rw
}

func (r *RegisterWorker) handle(ctx context.Context, e event.Event) {
	var p RegisterPayload
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
	wg          sync.WaitGroup
}

func NewResponseBus(queue <-chan event.Event, client *redis.Client) *RegisterBus {
	return &RegisterBus{queue: queue, client: client}
}

// Start launches n workers and blocks until ctx is cancelled, then drains remaining jobs.
func (b *RegisterBus) Start(ctx context.Context, n int) {
	for i := 0; i < n; i++ {
		b.spawnWorker(ctx)
	}
}

func (b *RegisterBus) spawnWorker(ctx context.Context) {
	id := b.workerCount.Add(1)
	rw := newResponseWorker(id, b.client)
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		defer rw.w.State.Store(false)
		for {
			select {
			case e, ok := <-b.queue:
				if !ok {
					return
				}
				rw.handle(ctx, e)
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
