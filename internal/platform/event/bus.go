package event

import (
	"context"
	"fmt"
	"messenger/messenger/internal/application/event"
	"messenger/messenger/internal/domain"
)

type Bus struct {
	handlers    map[string]map[string]event.Handler
	middlewares []Middleware
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string]map[string]event.Handler),
	}
}

func (b *Bus) Register(e domain.Event, h event.Handler) {
	for _, m := range b.middlewares {
		h = m(h)
	}

	_, ok := b.handlers[e.Name()]
	if !ok {
		b.handlers[e.Name()] = make(map[string]event.Handler, 0)
	}
	b.handlers[e.Name()][h.Name()] = h
}

func (b *Bus) Dispatch(ctx context.Context, e domain.Event) error {
	if handlers, ok := b.handlers[e.Name()]; ok {
		for _, h := range handlers {
			go h.Handle(ctx, e)
		}
	}
	return nil
}

func (b *Bus) DispatchForHandler(ctx context.Context, e domain.Event, h event.Handler) error {
	if _, ok := b.handlers[e.Name()]; !ok {
		return fmt.Errorf("there is no handlers for %q event", e.Name())
	}
	if _, ok := b.handlers[e.Name()][h.Name()]; !ok {
		return fmt.Errorf("there is no handler named %q for %q event", h.Name(), e.Name())
	}

	handler := b.handlers[e.Name()][h.Name()]
	return handler.Handle(ctx, e)
}

func (b *Bus) Use(m Middleware) {
	b.middlewares = append(b.middlewares, m)

	for eventName, eventHandlers := range b.handlers {
		for idx, h := range eventHandlers {
			b.handlers[eventName][idx] = m(h)
		}
	}
}

func (b *Bus) HandlersForEvent(e domain.Event) map[string]event.Handler {
	return b.handlers[e.Name()]
}

type Middleware func(event.Handler) event.Handler

// type HandlerFunc func(context.Context, domain.Event) error

// func (hf HandlerFunc) Handle(ctx context.Context, e domain.Event) error {
// 	return hf(ctx, e)
// }

// func (hf HandlerFunc) Name() string {

// }

