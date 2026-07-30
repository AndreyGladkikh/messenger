package event

import (
	"context"
	"fmt"
	"messenger/messenger/internal/messaging/application/event"
	"messenger/messenger/internal/messaging/domain"
	"slices"
)

type HandlerRegistry struct {
	handlers map[string][]Handler
}

func NewRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string][]Handler),
	}
}

func (r *HandlerRegistry) HandlersForEvent(event domain.Event) []Handler {
	return slices.Clone(r.handlers[event.Name()])
}

func RegisterEventHandler[E domain.Event](r *HandlerRegistry, eventName string, handler event.Handler[E]) {
	adaptedHandler := &ApEventHandlerAdapter[E]{handler}
	r.handlers[eventName] = append(r.handlers[eventName], adaptedHandler)
}

type Handler interface {
	Name() string
	Handle(context.Context, domain.Event) error
}

type ApEventHandlerAdapter[E domain.Event] struct {
	handler event.Handler[E]
}

func (h *ApEventHandlerAdapter[E]) Name() string {
	return h.handler.Name()
}

func (h *ApEventHandlerAdapter[E]) Handle(ctx context.Context, event domain.Event) error {
	typedEvent, ok := event.(E)
	if !ok {
		return fmt.Errorf("handler %s expected event of type %T, got %T", h.handler.Name(), new(E), event)
	}
	return h.handler.Handle(ctx, typedEvent)
}
