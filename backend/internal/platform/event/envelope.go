package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Envelope struct {
	ID uuid.UUID `json:"id"`
	Type string `json:"type"`
	OccurredAt time.Time `json:"occurredAt"`
	Payload json.RawMessage `json:"payload"`
}