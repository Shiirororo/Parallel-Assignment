package main

import (
	"net/http"
	bootstrap "parallel/init"
)

func main() {
	rd, db := bootstrap.Run()
	_, _ = rd, db
	http.HandleFunc("/class", func(w http.ResponseWriter, r *http.Request) {

	})
	http.ListenAndServe(":36789", nil)
}
