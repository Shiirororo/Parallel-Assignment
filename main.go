package main

import (
	"context"
	"fmt"
	"net/http"
	bootstrap "parallel/init"
	"parallel/internal/manager"
	"parallel/internal/service"
	logging "parallel/internal/worker/LoggingWorker"
	response "parallel/internal/worker/RegisterWorker"
)

func main() {
	//Generate context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rd, db := bootstrap.Run()
	_, _ = rd, db

	//Define router, bus
	router := manager.NewIngressRouter(256)

	registerBus := response.NewResponseBus(router.Response(), rd)
	loggingBus := logging.NewLoggingBus(router.Logging(), db)
	//counterBus := counter.NewCounterBus()

	registerBus.Start(ctx, 4) //init 4 workers
	loggingBus.Start(ctx, 4)  //init 4 workers

	fmt.Printf("Active business goroutines: %d\n", service.ActiveGoroutines())

	http.HandleFunc("/class/id=?", func(w http.ResponseWriter, r *http.Request) {
		//to do: build later
	})
	http.ListenAndServe(":36789", nil)
}
