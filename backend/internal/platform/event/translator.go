package event

import (
	"encoding/json"
	"fmt"
	"messenger/messenger/internal/messaging/domain/message"
	"messenger/messenger/internal/messaging/infrastructure/sqlc"
	"messenger/messenger/internal/shared/application/event"
	"messenger/messenger/internal/shared/domain"
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

func TranslateStoredEventToDomainEvent(storedEvent sqlc.Outbox) (domain.Event, error) {
	decoder, ok := decoders[storedEvent.EventType]
	if !ok {
		return nil, fmt.Errorf("not found decoder for event type %s", storedEvent.EventType)
	}

	return decoder(storedEvent.EventPayload)
}

func TranslateOutboxMsgToEventEnvelope(outboxMsg sqlc.Outbox) (event.Envelope[domain.Event], error) {
	var envelope event.Envelope[domain.Event]

	decoder, ok := decoders[outboxMsg.EventType]
	if !ok {
		return envelope, fmt.Errorf("outbox to event envelope translator: not found decoder for event type %s", outboxMsg.EventType)
	}

	e, err := decoder(outboxMsg.EventPayload)
	if err != nil {
		return envelope, fmt.Errorf("outbox to event envelope translator: failed to decode outbox payload: %w", err)
	}

	return event.Envelope[domain.Event]{
		ID: outboxMsg.EventID,
		OccurredAt: outboxMsg.OccurredAt.Time,
		Event: e,
	}, nil
}
