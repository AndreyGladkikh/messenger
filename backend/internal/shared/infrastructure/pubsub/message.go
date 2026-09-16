package pubsub

import (
	"encoding/json"
	"messenger/messenger/internal/platform/event"
)

type MessageType string

const (
	MessageTypeEvent = "event"
)

type Message struct {
	Type MessageType `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ChatEventsChannnelMessage struct {
	event.Envelope
}