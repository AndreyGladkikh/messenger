package event

import (
	"context"
	"fmt"
	"messenger/messenger/internal/shared/app/event"
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

func (r *HandlerRegistry) HandlersForEvent(event event.Event) []Handler {
	return slices.Clone(r.handlers[event.Name()])
}

func (r *HandlerRegistry) RegisterEventHandler[E event.Event](handler event.Handler[E]) {
	var e E
	adaptedHandler := &ApEventHandlerAdapter[E]{handler}
	r.handlers[e.Name()] = append(r.handlers[e.Name()], adaptedHandler)
}

type Handler interface {
	Name() string
	Handle(context.Context, event.Envelope[event.Event]) error
}

type ApEventHandlerAdapter[E event.Event] struct {
	handler event.Handler[E]
}

func (h *ApEventHandlerAdapter[E]) Name() string {
	return h.handler.Name()
}

func (h *ApEventHandlerAdapter[E]) Handle(ctx context.Context, e event.Envelope[event.Event]) error {
	var eTP E

	typedEvent, ok := e.Event.(E)
	if !ok {
		return fmt.Errorf("handler %s expected event of type %T, got %T", h.handler.Name(), eTP, e.Event)
	}

	typedEnvelope := event.Envelope[E]{
		ID:         e.ID,
		OccurredAt: e.OccurredAt,
		Event:      typedEvent,
	}
	return h.handler.Handle(ctx, typedEnvelope)
}
