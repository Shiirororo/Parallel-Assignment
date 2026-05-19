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
	logging "parallel/internal/worker/LoggingWorker"
	response "parallel/internal/worker/RegisterWorker"
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

	registerBus := response.NewRegisterBus(router.Response(), rd)
	loggingBus := logging.NewLoggingBus(router.Logging(), db)
	//counterBus := counter.NewCounterBus()

	registerBus.Start(ctx, 4) //init 4 workers
	loggingBus.Start(ctx, 4)  //init 4 workers

	fmt.Printf("Active business goroutines: %d\n", service.ActiveGoroutines())
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	// http.HandleFunc("api/class/id?", func(w http.ResponseWriter, r *http.Request) {
	// 	//to do: build later
	// 	query := r.URL.Query()
	// 	class_id := query.Get("id")
	// 	_ = class_id
	// })
	http.HandleFunc("api/class/getClassInfo", func(w http.ResponseWriter, r *http.Request) {
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

		result := make(map[string]string, len(raw)/2)
		for i := 0; i+1 < len(raw); i += 2 {
			result[raw[i]] = raw[i+1]
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
	http.HandleFunc("api/class/register", func(w http.ResponseWriter, r *http.Request) {

	})
	http.HandleFunc("api/class/unregister", func(w http.ResponseWriter, r *http.Request) {

	})
	http.ListenAndServe(":36789", nil)
}
