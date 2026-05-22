package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"parallel/internal/dto"
	"parallel/internal/event"
	"parallel/internal/service"
	"parallel/internal/worker"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//WORKER DECLARATIOn

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
	var p dto.LogginPayload
	if err := json.Unmarshal(e.Request.Payload, &p); err != nil {
		log.Printf("[LoggingWorker %d] bad payload: %v", lw.w.ID, err)
		return
	}
	collection := lw.client.Database("class-registration").Collection("success")

	switch p.Action {
	case 0:
		res, err := collection.InsertOne(ctx, p)
		if err != nil {
			log.Printf("[LoggingWorker %d] insert error: %v", lw.w.ID, err)
			return
		}
		fmt.Printf("[LoggingWorker %d] created record: %s\n", lw.w.ID, res.InsertedID)

	case 1:
		filter := bson.M{"student_id": p.StudentID}
		update := bson.M{"$push": bson.M{"success_class_id": bson.M{"$each": p.SuccessClassID}}}
		res, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Printf("[LoggingWorker %d] update error: %v", lw.w.ID, err)
			return
		}
		fmt.Printf("[LoggingWorker %d] updated record matched=%d modified=%d\n", lw.w.ID, res.MatchedCount, res.ModifiedCount)

	default:
		log.Printf("[LoggingWorker %d] unknown action: %d", lw.w.ID, p.Action)
	}
}

//BUS DECLARATION

type LoggingBus struct {
	queue       <-chan event.Event
	client      *mongo.Client
	workerCount atomic.Int32
	minWorker   int
	maxWorker   int
	wg          sync.WaitGroup
}

func NewLoggingBus(queue <-chan event.Event, client *mongo.Client) *LoggingBus {
	return &LoggingBus{queue: queue, client: client, minWorker: 10, maxWorker: 1000}
}

func (b *LoggingBus) Start(ctx context.Context, n int) {
	for i := 0; i < n; i++ {
		b.spawnWorker(ctx)
	}
	// go func() {

	// }
}

func (b *LoggingBus) spawnWorker(ctx context.Context) {
	if b.workerCount.Load() >= int32(b.maxWorker) {
		return
	}
	id := b.workerCount.Add(1)
	lw := newLoggingWorker(id, b.client)
	b.wg.Add(1)
	service.TrackGo(+1)

	go func() {
		defer b.wg.Done()
		defer service.TrackGo(-1)
		defer b.workerCount.Add(-1)
		defer lw.w.State.Store(false)
		idleTimer := time.NewTimer(30 * time.Second)
		defer idleTimer.Stop()
		for {
			select {
			case e, ok := <-b.queue:
				if !ok {
					return
				}
				if !idleTimer.Stop() {
					<-idleTimer.C
				}
				idleTimer.Reset(30 * time.Second)

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
			case <-idleTimer.C:
				if b.workerCount.Load() > int32(b.minWorker) {
					return
				}
			}
		}
	}()
}
