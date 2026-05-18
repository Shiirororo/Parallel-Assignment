package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"parallel/internal/event"
	"parallel/internal/service"
	"parallel/internal/worker"
	"sync"
	"sync/atomic"

	"go.mongodb.org/mongo-driver/mongo"
)

type LogginPayload struct {
	StudentID      string `json:"student_id"`
	SuccessClassID []int  `json:"success_class_id"`
}

type LoggingWorker struct {
	client *mongo.Client
	w      worker.Worker
}

func newLoggingWorker(id int32, client *mongo.Client) *LoggingWorker {
	lw := &LoggingWorker{client: client}
	lw.w.ID = id
	lw.w.Job = "Logging"
	lw.w.State.Store(true)
	return lw
}

func (lw *LoggingWorker) handle(ctx context.Context, e event.Event) {
	var p LogginPayload
	if err := json.Unmarshal(e.Request.Payload, &p); err != nil {
		log.Printf("[ResponseWorker %d] bad payload: %v", lw.w.ID, err)
	}
	collection := lw.client.Database("class-registration").Collection("success")
	res, err := collection.InsertOne(ctx, p)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully Insert into a table: %s", res.InsertedID)
}

type LoggingBus struct {
	queue       <-chan event.Event
	client      *mongo.Client
	workerCount atomic.Int32
	wg          sync.WaitGroup
}

func NewLoggingBus(queue <-chan event.Event, client *mongo.Client) *LoggingBus {
	return &LoggingBus{queue: queue, client: client}
}

func (b *LoggingBus) Start(ctx context.Context, n int) {
	for i := 0; i < n; i++ {
		b.spawnWorker(ctx)
	}
}

func (b *LoggingBus) spawnWorker(ctx context.Context) {
	id := b.workerCount.Add(1)
	lw := newLoggingWorker(id, b.client)
	b.wg.Add(1)
	service.TrackGo(+1)
	go func() {
		defer b.wg.Done()
		defer service.TrackGo(-1)
		defer lw.w.State.Store(false)
		for {
			select {
			case e, ok := <-b.queue:
				if !ok {
					return
				}
				lw.handle(ctx, e)
			case <-ctx.Done():
				for {
					select {
					case e, ok := <-b.queue:
						if !ok {
							return
						}
						lw.handle(context.Background(), e)
					default:
						return
					}
				}
			}
		}
	}()
}
