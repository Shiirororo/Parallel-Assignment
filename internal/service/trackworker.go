package service

import "sync/atomic"

var goroutineCount atomic.Int32

func TrackGo(delta int32)     { goroutineCount.Add(delta) }
func ActiveGoroutines() int32 { return goroutineCount.Load() }

//Basically, we have another selection using net/http/pprof
// to track number of Goroutine

//code sample
// import _ "net/http/pprof"
// import "net/http"

// func main() {
// 	go func() {
// 		http.ListenAndServe("localhost:6060", nil)
// 	}()

// 	// app code
// }
