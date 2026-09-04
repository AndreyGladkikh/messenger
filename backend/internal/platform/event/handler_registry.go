package event

import (
	"context"
	"fmt"
	"messenger/messenger/internal/messaging/application/event"
	sharedDomain "messenger/messenger/internal/shared/domain"
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

func (r *HandlerRegistry) HandlersForEvent(event sharedDomain.Event) []Handler {
	return slices.Clone(r.handlers[event.Name()])
}

func RegisterEventHandler[E sharedDomain.Event](r *HandlerRegistry, eventName string, handler event.Handler[E]) {
	adaptedHandler := &ApEventHandlerAdapter[E]{handler}
	r.handlers[eventName] = append(r.handlers[eventName], adaptedHandler)
}

type Handler interface {
	Name() string
	Handle(context.Context, sharedDomain.Event) error
}

type ApEventHandlerAdapter[E sharedDomain.Event] struct {
	handler event.Handler[E]
}

func (h *ApEventHandlerAdapter[E]) Name() string {
	return h.handler.Name()
}

func (h *ApEventHandlerAdapter[E]) Handle(ctx context.Context, event sharedDomain.Event) error {
	typedEvent, ok := event.(E)
	if !ok {
		return fmt.Errorf("handler %s expected event of type %T, got %T", h.handler.Name(), new(E), event)
	}
	return h.handler.Handle(ctx, typedEvent)
}
