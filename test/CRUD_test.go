package test

import (
	"context"
	"encoding/json"
	bootstrap "parallel/init"
	"parallel/internal/dto"
	"parallel/internal/event"
	"parallel/internal/manager"
	counter "parallel/internal/worker/CounterWorker"
	logging "parallel/internal/worker/LoggingWorker"
	register "parallel/internal/worker/RegisterWorker"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	db     *mongo.Client
	router *manager.IngressRouter
)

// TestLoggingWorker tests create and update event types through the LoggingBus.
func TestLoggingWorker(t *testing.T) {
	col := db.Database("class-registration").Collection("success")
	ctx := context.Background()
	studentID := "stu-test-001"

	col.DeleteMany(ctx, bson.M{"student_id": studentID})
	defer col.DeleteMany(ctx, bson.M{"student_id": studentID})

	publish := func(p dto.LogginPayload) {
		raw, _ := json.Marshal(p)
		router.Publish("logging", event.Event{
			Request: event.Request{OriginID: p.StudentID, Payload: raw},
		})
		time.Sleep(200 * time.Millisecond)
	}

	// Action 0: create
	publish(dto.LogginPayload{Action: 0, StudentID: studentID, SuccessClassID: []int{101, 202}})

	var result dto.LogginPayload
	if err := col.FindOne(ctx, bson.M{"student_id": studentID}).Decode(&result); err != nil {
		t.Fatalf("create: document not found: %v", err)
	}
	if len(result.SuccessClassID) != 2 {
		t.Errorf("create: SuccessClassID length: got %d, want 2", len(result.SuccessClassID))
	}

	// Action 1: update — append new class IDs
	publish(dto.LogginPayload{Action: 1, StudentID: studentID, SuccessClassID: []int{303}})

	if err := col.FindOne(ctx, bson.M{"student_id": studentID}).Decode(&result); err != nil {
		t.Fatalf("update: document not found: %v", err)
	}
	if len(result.SuccessClassID) != 3 {
		t.Errorf("update: SuccessClassID length: got %d, want 3", len(result.SuccessClassID))
	}
	if result.SuccessClassID[2] != 303 {
		t.Errorf("update: last class ID: got %d, want 303", result.SuccessClassID[2])
	}
}

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, file, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(file), "..")

	rd, mongoClient := bootstrap.RunFrom(root)
	db = mongoClient
	router = manager.NewIngressRouter(256)

	register.NewRegisterBus(router.Response(), rd).Start(ctx, 4)
	logging.NewLoggingBus(router.Logging(), db).Start(ctx, 4)
	counter.NewCounterBus(router.Counter(), rd).Start(ctx, 4)

	m.Run()
}
