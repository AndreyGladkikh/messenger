package event

import (
	"context"
	"fmt"
	"maps"
	"messenger/messenger/internal/application/event"
	"messenger/messenger/internal/domain"
)

type Bus struct {
	handlers    map[string]map[string]Handler
	middlewares []Middleware
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string]map[string]Handler),
	}
}

func (b *Bus) Dispatch(ctx context.Context, e domain.Event) error {
	if handlers, ok := b.handlers[e.Name()]; ok {
		for _, h := range handlers {
			go h.Handle(ctx, e)
		}
	}
	return nil
}

func (b *Bus) DispatchForHandler(ctx context.Context, e domain.Event, h Handler) error {
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

func (b *Bus) HandlersForEvent(e domain.Event) map[string]Handler {
	return maps.Clone(b.handlers[e.Name()])
}

func RegisterHandler[E domain.Event](b *Bus, e E, h event.Handler[E]) {
	var handler Handler = EventHandlerAdapter[E]{handler: h}

	for _, m := range b.middlewares {
		handler = m(handler)
	}

	_, ok := b.handlers[e.Name()]
	if !ok {
		b.handlers[e.Name()] = make(map[string]Handler, 0)
	}
	b.handlers[e.Name()][h.Name()] = handler
}

type Middleware func(Handler) Handler

type Handler interface {
	Handle(context.Context, domain.Event) error
	Name() string
}

type EventHandlerAdapter[T domain.Event] struct {
	handler event.Handler[T]
}

func (h EventHandlerAdapter[T]) Handle(ctx context.Context, e domain.Event) error {
	typedEvent, ok := e.(T)
	if !ok {
		return fmt.Errorf("unexpected event type for handler %q", h.Name())
	}

	return h.handler.Handle(ctx, typedEvent)
}

func (h EventHandlerAdapter[T]) Name() string {
	return h.handler.Name()
}