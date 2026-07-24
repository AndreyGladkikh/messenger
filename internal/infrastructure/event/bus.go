package event

import (
	"context"
	"fmt"
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

func (b *Bus) Dispatch(ctx context.Context, envelope *Envelope) error {
	if handlers, ok := b.handlers[envelope.Event.Name()]; ok {
		for _, h := range handlers {
			if err := h.Handle(ctx, envelope); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Bus) Use(m Middleware) {
	b.middlewares = append(b.middlewares, m)

	for eventName, eventHandlers := range b.handlers {
		for idx, h := range eventHandlers {
			b.handlers[eventName][idx] = m(h)
		}
	}
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
	Handle(context.Context, *Envelope) error
}

type HandlerFunc func(context.Context, *Envelope) error

func (hf HandlerFunc) Handle(ctx context.Context, e *Envelope) error {
	return hf(ctx, e)
}

type EventHandlerAdapter[T domain.Event] struct {
	handler event.Handler[T]
}

func (h EventHandlerAdapter[E]) Handle(ctx context.Context, envelope *Envelope) error {
	typedEvent, ok := envelope.Event.(E)
	if !ok {
		return fmt.Errorf("handler %q: expected %T, got %T", h.handler.Name(), *new(E), envelope)
	}

	return h.handler.Handle(ctx, typedEvent)
}