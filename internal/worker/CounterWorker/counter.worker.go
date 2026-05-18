package counter

import "parallel/internal/worker"

type GetClassListPayload struct {
}

type CounterBus struct {
	w worker.Worker
}
