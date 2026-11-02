package event

import (
	"encoding/json"
	"messenger/messenger/internal/shared/app/event"
	// "time"
	// "uuid"
)

type Envelope struct {
	event.Envelope[event.Event]
	// ID         uuid.UUID `json:"id"`
	// OccurredAt time.Time `json:"occurredAt"`
	Type  string          `json:"type"`
	Event json.RawMessage `json:"event"`
}
