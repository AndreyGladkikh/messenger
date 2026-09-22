package event

import (
	"context"
	"fmt"
	"messenger/messenger/internal/shared/application/event"
	"messenger/messenger/internal/shared/domain"
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

func RegisterEventHandler[E domain.Event](r *HandlerRegistry, handler event.Handler[E]) {
	e := *new(E)
	adaptedHandler := &ApEventHandlerAdapter[E]{handler}
	r.handlers[e.Name()] = append(r.handlers[e.Name()], adaptedHandler)
}

type Handler interface {
	Name() string
	Handle(context.Context, event.Envelope[domain.Event]) error
}

type ApEventHandlerAdapter[E domain.Event] struct {
	handler event.Handler[E]
}

func (h *ApEventHandlerAdapter[E]) Name() string {
	return h.handler.Name()
}

func (h *ApEventHandlerAdapter[E]) Handle(ctx context.Context, e event.Envelope[domain.Event]) error {
	typedEvent, ok := e.Event.(E)
	if !ok {
		return fmt.Errorf("handler %s expected event of type %T, got %T", h.handler.Name(), new(E), e.Event)
	}

	typedEnvelope := event.Envelope[E]{
		ID:         e.ID,
		OccurredAt: e.OccurredAt,
		Event:      typedEvent,
	}
	return h.handler.Handle(ctx, typedEnvelope)
}
