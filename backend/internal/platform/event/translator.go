package event

import (
	"encoding/json"
	"fmt"
	cnvEvnt "messenger/messenger/internal/messaging/app/event"
	"messenger/messenger/internal/messaging/domain/model/message"
	"messenger/messenger/internal/platform/sqlc"
	"messenger/messenger/internal/shared/app/event"
)

type Decoder func([]byte) (event.Event, error)

var decoders = map[string]Decoder{
	cnvEvnt.DirectConversationStartedEventName: decodeAs[cnvEvnt.DirectConversationStarted],
	message.MessageSentEventName:               decodeAs[message.MessageSent],
}

func decodeAs[E event.Event](payload []byte) (event.Event, error) {
	var e E
	if err := json.Unmarshal(payload, &e); err != nil {
		return nil, fmt.Errorf("deserialize stored event payload: %w", err)
	}
	return e, nil
}

func TranslateStoredEventToDomainEvent(storedEvent sqlc.Outbox) (event.Event, error) {
	decoder, ok := decoders[storedEvent.Type]
	if !ok {
		return nil, fmt.Errorf("not found decoder for event type %s", storedEvent.Type)
	}

	return decoder(storedEvent.Payload)
}

func TranslateOutboxMsgToEventEnvelope(outboxMsg sqlc.Outbox) (event.Envelope[event.Event], error) {
	var envelope event.Envelope[event.Event]

	decoder, ok := decoders[outboxMsg.Type]
	if !ok {
		return envelope, fmt.Errorf("outbox to event envelope translator: not found decoder for event type %s", outboxMsg.Type)
	}

	e, err := decoder(outboxMsg.Payload)
	if err != nil {
		return envelope, fmt.Errorf("outbox to event envelope translator: failed to decode outbox payload: %w", err)
	}

	return event.Envelope[event.Event]{
		OccurredAt: outboxMsg.OccurredAt,
		Event:      e,
	}, nil
}
