package event

import (
	"encoding/json/jsontext"
	"messenger/messenger/internal/shared/app/event"
)

type Envelope struct {
	event.Envelope[event.Event]
	Type  string          `json:"type"`
	Event jsontext.Value `json:"event"`
}
