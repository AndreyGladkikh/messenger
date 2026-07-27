package event

import "messenger/messenger/internal/domain"

type Envelope struct {
	EventID    string
	Event domain.Event
}

func NewEnvelope(
	id string,
	event domain.Event,
) *Envelope {
	return &Envelope{
		EventID:    id,
		Event: event,
	}
}
