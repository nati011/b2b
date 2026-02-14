package event

import (
	"sync"
)

type Broker struct {
	mu          sync.RWMutex
	subscribers map[string][]chan any
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string][]chan any),
	}
}

// Subscribe returns a channel for receiving messages on a topic
func (b *Broker) Subscribe(topic string) <-chan any {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan any, 16) // buffered async queue
	b.subscribers[topic] = append(b.subscribers[topic], ch)
	return ch
}

// Publish sends a message to all subscribers of a topic
func (b *Broker) Publish(topic string, msg any) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.subscribers[topic] {
		select {
		case ch <- msg:
		default: // drop if queue is full (avoid deadlock)
		}
	}
}

// CloseTopic closes all subscriber channels for a topic
func (b *Broker) CloseTopic(topic string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, ch := range b.subscribers[topic] {
		close(ch)
	}
	delete(b.subscribers, topic)
}
