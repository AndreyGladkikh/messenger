package event

import (
	"context"
	"messenger/messenger/internal/shared/domain"
	"time"

	"github.com/google/uuid"
)

type Handler[E domain.Event] interface {
	Handle(context.Context, Envelope[E]) error
	Name() string
}

type Envelope[E domain.Event] struct {
	ID uuid.UUID `json:"id"`
	OccurredAt time.Time `json:"occurredAt"`
	Event E `json:"event"`
}

func AsEnvelopeWithDomainEvent[E domain.Event](e Envelope[E]) Envelope[domain.Event] {
	return Envelope[domain.Event]{
		ID: e.ID,
		OccurredAt: e.OccurredAt,
		Event: e.Event,
	}
}