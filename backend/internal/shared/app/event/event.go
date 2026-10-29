package event

import (
	"context"
	"time"
	"uuid"
)

type Event interface {
	Name() string
}

type Handler[E Event] interface {
	Handle(context.Context, Envelope[E]) error
	Name() string
}

type Envelope[E Event] struct {
	ID         uuid.UUID `json:"id"`
	OccurredAt time.Time `json:"occurredAt"`
	Event      E         `json:"event"`
}

type Publisher interface {
	Publish(ctx context.Context, event Event)
}
