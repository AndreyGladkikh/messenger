package event

import (
	"encoding/json"
	"messenger/messenger/internal/messaging/domain/event"
	"time"

	"github.com/google/uuid"
)

type Envelope struct {
	ID uuid.UUID `json:"id"`
	Type string `json:"type"`
	OccurredAt time.Time `json:"occurredAt"`
	Event event.DomainEvent `json:"event"`
}

type RawEnvelope struct {
	ID uuid.UUID `json:"id"`
	Type string `json:"type"`
	OccurredAt time.Time `json:"occurredAt"`
	Payload json.RawMessage `json:"payload"`
}

// func (e Envelope) GetEvent() event {

// }