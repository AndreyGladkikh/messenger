package event

import "messenger/messenger/internal/domain"

type Envelope struct {
	ID    string
	Event domain.Event
}

func NewEnvelope(
	id string,
	event domain.Event,
) *Envelope {
	return &Envelope{
		ID:    id,
		Event: event,
	}
}
