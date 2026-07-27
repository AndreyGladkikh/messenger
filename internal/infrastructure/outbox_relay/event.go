package outbox_relay

import (
	"encoding/json"
	"fmt"
	"messenger/messenger/internal/domain"
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/infrastructure/db"
)

type eventStatus string

const (
	eventStatusPending    eventStatus = "pending"
	eventStatusProcessing eventStatus = "processing"
	eventStatusRetry      eventStatus = "retry"
	eventStatusSucceeded  eventStatus = "succeeded"
	eventStatusDead       eventStatus = "dead"
)

type Decoder func([]byte) (domain.Event, error)

var decoders = map[string]Decoder{
	message.MessageSentEventName: decodeAs[message.MessageSent],
}

func decodeAs[E domain.Event](payload []byte) (domain.Event, error) {
	var e E
	if err := json.Unmarshal(payload, &e); err != nil {
		return nil, fmt.Errorf("failed to deserialize stored event payload: %w", err)
	}
	return e, nil
}

func translateStoredEventToDomainEvent(storedEvent db.Outbox) (domain.Event, error) {
	decoder, ok := decoders[storedEvent.EventType]
	if !ok {
		return nil, fmt.Errorf("not found decoder for event type %s", storedEvent.EventType)
	}

	return decoder(storedEvent.EventPayload)
}