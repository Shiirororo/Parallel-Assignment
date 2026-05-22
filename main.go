package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	bootstrap "parallel/init"
	"parallel/internal/dto"
	"parallel/internal/event"
	luascripting "parallel/internal/lua-scripting"
	"parallel/internal/manager"
	"parallel/internal/service"
	counter "parallel/internal/worker/CounterWorker"
	logging "parallel/internal/worker/LoggingWorker"
	register "parallel/internal/worker/RegisterWorker"
	"strings"
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

	fmt.Printf("Active business goroutines: %d\n", service.ActiveGoroutines())
	go func() { http.ListenAndServe("localhost:6060", nil) }()

	http.HandleFunc("/api/class/getClassInfo", func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query().Get("ids")
		if ids == "" {
			http.Error(w, "missing ids", http.StatusBadRequest)
			return
		}
		keys := strings.Split(ids, ",")
		raw, err := luascripting.GetClassScript.Run(r.Context(), rd, keys).StringSlice()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		type ClassInfo struct {
			RemainSlot string `json:"remain_slot"`
		}
		result := make(map[string]ClassInfo, len(raw)/2)
		for i := 0; i+1 < len(raw); i += 2 {
			result[raw[i]] = ClassInfo{RemainSlot: raw[i+1]}
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
