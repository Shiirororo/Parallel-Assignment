package manager

type IngressRouter struct {
	cBus CounterBus
	rBus ResponseBus
	lBus LoggingBus
}

type CounterBus struct {
	queue chan Event
}
type LoggingBus struct {
	queue chan Event
}
type ResponseBus struct {
	queue chan Event
}

func NewIngressRouter(size int) *IngressRouter {
	return &IngressRouter{
		cBus: CounterBus{queue: make(chan Event, size)},
		rBus: ResponseBus{queue: make(chan Event, size)},
		lBus: LoggingBus{queue: make(chan Event, size)},
	}
}

func (ig *IngressRouter) Publish(bus string, e Event) {
	switch bus {
	case "counter":
		ig.cBus.queue <- e
	case "response":
		ig.rBus.queue <- e
	case "logging":
		ig.lBus.queue <- e
	}
}

func (ig *IngressRouter) Counter() <-chan Event  { return ig.cBus.queue }
func (ig *IngressRouter) Response() <-chan Event { return ig.rBus.queue }
func (ig *IngressRouter) Logging() <-chan Event  { return ig.lBus.queue }
