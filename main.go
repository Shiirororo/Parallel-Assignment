package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	bootstrap "parallel/init"
	luascripting "parallel/internal/lua-scripting"
	"parallel/internal/manager"
	"parallel/internal/service"
	counter "parallel/internal/worker/CounterWorker"
	logging "parallel/internal/worker/LoggingWorker"
	register "parallel/internal/worker/RegisterWorker"
	"strings"
)

func main() {
	//Generate context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rd, db := bootstrap.Run()
	_ = db

	//Define router, bus
	router := manager.NewIngressRouter(256)

	//insert queue spawn from router
	registerBus := register.NewRegisterBus(router.Response(), rd)
	loggingBus := logging.NewLoggingBus(router.Logging(), db)
	counterBus := counter.NewCounterBus(router.Counter(), rd)

	registerBus.Start(ctx, 4) //init 4 workers
	loggingBus.Start(ctx, 4)  //init 4 workers
	counterBus.Start(ctx, 4)  //init 4 workers

	fmt.Printf("Active business goroutines: %d\n", service.ActiveGoroutines())
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
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

	})
	http.HandleFunc("/api/class/unregister", func(w http.ResponseWriter, r *http.Request) {

	})
	http.HandleFunc("/api/scale-up", func(w http.ResponseWriter, r *http.Request) {
		worker := r.URL.Query().Get("worker")
		if worker == "" {
			http.Error(w, "missing worker id", http.StatusBadRequest)
			return
		}
		// 	select {
		// 	case worker := "register":
		// 		registerBus.spawnWorker()
		// 	default:
		// 		return
		// }
		if worker == "register" {
			// registerBus.spawnWorker(context.Context)
		}
	})
	http.ListenAndServe(":36789", nil)
}
