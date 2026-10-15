package event

import (
	"context"
	"time"
	"uuid"
)

type Handler[E any] interface {
	Handle(context.Context, Envelope[E]) error
	Name() string
}

type Envelope[E any] struct {
	ID         uuid.UUID `json:"id"`
	OccurredAt time.Time `json:"occurredAt"`
	Event      E         `json:"event"`
}
