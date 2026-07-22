package event

import (
	"context"
	"fmt"
	"messenger/messenger/internal/application/event"
	"messenger/messenger/internal/domain"
)

type Bus[T domain.Event] struct {
	handlers    map[string]map[string]event.Handler[T]
	middlewares []Middleware[T]
}

func NewBus[T domain.Event]() *Bus[T] {
	return &Bus[T]{
		handlers: make(map[string]map[string]event.Handler[T]),
	}
}

func (b *Bus[T]) Register(e T, h event.Handler[T]) {
	for _, m := range b.middlewares {
		h = m(h)
	}

	_, ok := b.handlers[e.Name()]
	if !ok {
		b.handlers[e.Name()] = make(map[string]event.Handler[T], 0)
	}
	b.handlers[e.Name()][h.Name()] = h
}

func (b *Bus[T]) Dispatch(ctx context.Context, e T) error {
	if handlers, ok := b.handlers[e.Name()]; ok {
		for _, h := range handlers {
			go h.Handle(ctx, e)
		}
	}
	return nil
}

func (b *Bus[T]) DispatchForHandler(ctx context.Context, e T, h event.Handler[T]) error {
	if _, ok := b.handlers[e.Name()]; !ok {
		return fmt.Errorf("there is no handlers for %q event", e.Name())
	}
	if _, ok := b.handlers[e.Name()][h.Name()]; !ok {
		return fmt.Errorf("there is no handler named %q for %q event", h.Name(), e.Name())
	}

	handler := b.handlers[e.Name()][h.Name()]
	return handler.Handle(ctx, e)
}

func (b *Bus[T]) Use(m Middleware[T]) {
	b.middlewares = append(b.middlewares, m)

	for eventName, eventHandlers := range b.handlers {
		for idx, h := range eventHandlers {
			b.handlers[eventName][idx] = m(h)
		}
	}
}

func (b *Bus[T]) HandlersForEvent(e T) map[string]event.Handler[T] {
	return b.handlers[e.Name()]
}

type Middleware[T domain.Event] func(event.Handler[T]) event.Handler[T]

type HandlerFunc func(context.Context, domain.Event) error

func (hf HandlerFunc) Handle(ctx context.Context, e domain.Event) error {
	return hf(ctx, e)
}

func (hf HandlerFunc) Name() string {

}

