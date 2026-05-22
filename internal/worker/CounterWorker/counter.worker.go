package counter

import (
	"context"
	"parallel/internal/event"
	"parallel/internal/worker"
	"sync/atomic"

	"github.com/redis/go-redis/v9"
)

type CounterBus struct {
	w           worker.Worker
	workerCount atomic.Int32
	client      *redis.Client
	queue       <-chan event.Event
	minWorker   int
	maxWorker   int
}

func NewCounterBus(queue <-chan event.Event, client *redis.Client) *CounterBus {
	return &CounterBus{queue: queue, client: client, minWorker: 10, maxWorker: 1000}
}
func (b *CounterBus) Start(ctx context.Context, n int) {
	for i := 0; i < n; i++ {
		// b.spawnWorker(ctx)
	}
	// go func() {

	// }
}

// counter worker
type CounterWorker struct {
	client *redis.Client
	w      worker.Worker
}

func newCounterWorker(id int32, client *redis.Client) *CounterWorker {
	rw := &CounterWorker{client: client}
	rw.w.ID = id
	rw.w.Job = "Register"
	rw.w.State.Store(true)
	return rw
}
