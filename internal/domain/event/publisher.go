package event

import (
	"context"
	"sync"
)

var once sync.Once

var publisher *EventPublisher

func Publisher() *EventPublisher {
	once.Do(func() {
		publisher = newEventPublisher()
	})

	return publisher
}

type EventSubscriber interface {
	supportsEvent(event DomainEvent) bool
	handleEvent(ctx context.Context, event DomainEvent) error
}

type EventHandler func(event DomainEvent) error

// type EventHandler interface {
// 	handle(event DomainEvent) error
// }

type EventPublisher struct {
	mu          sync.Mutex
	subscribers []EventSubscriber
	handlers    map[string][]EventHandler
}

func newEventPublisher() *EventPublisher {
	return &EventPublisher{
		handlers: make(map[string][]EventHandler),
	}
}

func (p *EventPublisher) Publish(e DomainEvent) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, h := range p.handlers[e.Name()] {
		h(e)
	}
}

func (p *EventPublisher) Subscribe(h EventHandler, events ...DomainEvent) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, e := range events {
		handlers := p.handlers[e.Name()]
		handlers = append(handlers, h)
		p.handlers[e.Name()] = handlers
	}
}

// func (p *EventPublisher) Publish(e DomainEvent) {
// 	for _, s := range p.subscribers {
// 		if s.supportsEvent(e) {
// 			s.handleEvent(e)
// 		}
// 	}
// }

// func (p *EventPublisher) Subscribe(s EventSubscriber, events ...DomainEvent) {
// 	for _, e := range events {
// 		subscribers := p.subscribers[e.Name()]
// 		subscribers = append(subscribers, s)
// 		p.subscribers[e.Name()] = subscribers
// 	}
// }
