package manager

type Bus struct {
	queue chan Event
}

func NewBus(size int) *Bus {
	return &Bus{
		queue: make(chan Event, size),
	}
}

func (b *Bus) Publish(e Event) {
	b.queue <- e
}

func (b *Bus) Queue() <-chan Event {
	return b.queue
}
