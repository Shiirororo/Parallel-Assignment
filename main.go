package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	bootstrap "parallel/init"
	"parallel/internal/dto"
	"parallel/internal/event"
	"parallel/internal/manager"
	"parallel/internal/service"
	counter "parallel/internal/worker/CounterWorker"
	logging "parallel/internal/worker/LoggingWorker"
	register "parallel/internal/worker/RegisterWorker"
	"strings"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rd, db := bootstrap.Run()

	router := manager.NewIngressRouter(256)

	registerBus := register.NewRegisterBus(router.Response(), rd)
	loggingBus := logging.NewLoggingBus(router.Logging(), db)
	counterBus := counter.NewCounterBus(router.Counter(), rd)

	registerBus.Start(ctx, 4)
	loggingBus.Start(ctx, 4)
	counterBus.Start(ctx, 4)

	cache := service.NewCache(rd)
	if err := cache.LoadCache(ctx); err != nil {
		panic(err)
	} //init load cache, if fail we get diddied

	//parallel local cache layer, reduce stress on redis
	go func() {
		for {
			if err := cache.LoadCache(ctx); err != nil {
				log.Println(err)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	fmt.Printf("Active business goroutines: %d\n", service.ActiveGoroutines())
	//go func() { http.ListenAndServe("localhost:6060", nil) }()

	http.HandleFunc("/api/class/getClassInfo", func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query().Get("ids")
		if ids == "" {
			http.Error(w, "missing ids", http.StatusBadRequest)
			return
		}
		keys := strings.Split(ids, ",")

		result := make(map[string]dto.ClassInfo, len(keys))
		for _, key := range keys {
			if v, ok := cache.Get(key); ok {
				result[key] = dto.ClassInfo{RemainSlot: v}
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	http.HandleFunc("/api/class/register", func(w http.ResponseWriter, r *http.Request) {
		var p dto.RegisterPayload
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil || p.ClassID == "" || p.StudentID == "" {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		raw, _ := json.Marshal(p)
		e := event.Event{
			Type:    "register",
			Request: event.Request{OriginID: p.ResponseCh, Payload: raw},
		}

		// RegisterWorker: decrement Redis slot via Lua
		router.Publish("response", e)
		// CounterWorker: track active registrations
		router.Publish("counter", e)
		// LoggingWorker: persist to MongoDB (action=0 create)
		logRaw, _ := json.Marshal(dto.LogginPayload{Action: 0, StudentID: p.StudentID})
		router.Publish("logging", event.Event{
			Request: event.Request{OriginID: p.StudentID, Payload: logRaw},
		})

		w.WriteHeader(http.StatusAccepted)
	})

	http.HandleFunc("/api/class/unregister", func(w http.ResponseWriter, r *http.Request) {

	})

	http.HandleFunc("/api/scale-up", func(w http.ResponseWriter, r *http.Request) {
		worker := r.URL.Query().Get("worker")
		if worker == "" {
			http.Error(w, "missing worker id", http.StatusBadRequest)
			return
		}
		if worker == "register" {
			// registerBus.spawnWorker(ctx)
		}
	})

	http.ListenAndServe(":36789", nil)
}
