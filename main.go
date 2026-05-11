package main

import (
	"net/http"
	bootstrap "parallel/init"
	"parallel/internal/manager"
)

func main() {
	rd, db := bootstrap.Run()
	_, _ = rd, db

	router := manager.NewIngressRouter(256)
	_ = router

	http.HandleFunc("/class/id=?", func(w http.ResponseWriter, r *http.Request) {

	})
	http.ListenAndServe(":36789", nil)
}
