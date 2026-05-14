package process

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"parallel/internal/event"
	"parallel/internal/worker"

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
