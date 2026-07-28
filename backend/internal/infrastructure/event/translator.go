package event

import (
	"encoding/json"
	"fmt"
	"messenger/messenger/internal/domain"
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/infrastructure/sqlc"
)

type Decoder func([]byte) (domain.Event, error)

var decoders = map[string]Decoder{
	message.MessageSentEventName: decodeAs[message.MessageSent],
}

func decodeAs[E domain.Event](payload []byte) (domain.Event, error) {
	var e E
	if err := json.Unmarshal(payload, &e); err != nil {
		return nil, fmt.Errorf("deserialize stored event payload: %w", err)
	}
	return e, nil
}

func translateStoredEventToDomainEvent(storedEvent sqlc.Outbox) (domain.Event, error) {
	decoder, ok := decoders[storedEvent.EventType]
	if !ok {
		return nil, fmt.Errorf("not found decoder for event type %s", storedEvent.EventType)
	}

	return decoder(storedEvent.EventPayload)
}
