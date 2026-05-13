package worker

import "sync/atomic"

type Worker struct {
	ID    int32
	State atomic.Bool
	Job   string
}
