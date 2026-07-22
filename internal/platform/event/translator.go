package event

import (
	"encoding/json"
	"fmt"
	"messenger/messenger/internal/adapters/out/postgres/queries"
	"messenger/messenger/internal/domain"
	"messenger/messenger/internal/domain/message"
)

func translateStoredEventToDispatchedEvent(storedEvent queries.Event) (domain.Event, error) {
	var event domain.Event

	switch storedEvent.EventType {
	case message.MessageSent{}.Name():
		event = message.MessageSent{}
	default:
		return nil, fmt.Errorf("unidentified stored event type: %q", storedEvent.EventType)
	}

	err := rehydrateEvent(storedEvent, event)
	return event, err
}

func rehydrateEvent(fromStoredEvent queries.Event, toEvent domain.Event) error {
	err := json.Unmarshal(fromStoredEvent.EventPayload, &toEvent)
	if err != nil {
		return fmt.Errorf("failed to deserialize stored event payload: %w", err)
	}
	return nil
}

// func rehydrateEvent(fromStoredEvent queries.Event, toEvent Event) error {
// 	if !fromStoredEvent.EventPayload.Valid {
// 		return nil
// 	}

// 	var data []byte

// 	a := fromStoredEvent.EventPayload.RawMessage.UnmarshalJSON(data)
// 	b, _ := fromStoredEvent.EventPayload.RawMessage.MarshalJSON()
// 	_ = a
// 	_ = b

// 	err := json.Unmarshal([]byte(fromStoredEvent.EventPayload.RawMessage), &toEvent)
// 	if err != nil {
// 		return fmt.Errorf("failed to deserialize stored event payload: %w", err)
// 	}
// 	return nil
// }
