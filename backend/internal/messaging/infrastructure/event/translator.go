package event

import (
	"encoding/json"
	"fmt"
	sharedDomain "messenger/messenger/internal/shared/domain"
	"messenger/messenger/internal/messaging/domain/message"
	"messenger/messenger/internal/messaging/infrastructure/sqlc"
)

type Decoder func([]byte) (sharedDomain.Event, error)

var decoders = map[string]Decoder{
	message.MessageSentEventName: decodeAs[message.MessageSent],
}

func decodeAs[E sharedDomain.Event](payload []byte) (sharedDomain.Event, error) {
	var e E
	if err := json.Unmarshal(payload, &e); err != nil {
		return nil, fmt.Errorf("deserialize stored event payload: %w", err)
	}
	return e, nil
}

func translateStoredEventToDomainEvent(storedEvent sqlc.Outbox) (sharedDomain.Event, error) {
	decoder, ok := decoders[storedEvent.EventType]
	if !ok {
		return nil, fmt.Errorf("not found decoder for event type %s", storedEvent.EventType)
	}

	return decoder(storedEvent.EventPayload)
}
