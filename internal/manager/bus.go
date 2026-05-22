package manager

import (
	// "context"
	counter "parallel/internal/worker/CounterWorker"
	logging "parallel/internal/worker/LoggingWorker"
	register "parallel/internal/worker/RegisterWorker"
)

type Bus struct {
	counter  counter.CounterBus
	logging  logging.LoggingBus
	register register.RegisterBus
}


