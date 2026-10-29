package pubsub

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"time"
	"uuid"
)

// todo rename to MessageEnvelope
type Message struct {
	Type    string
	Payload jsontext.Value
}

const DirectConversationCreatedName = "DirectConversationCreated"

type DirectConversationCreated struct {
	ConversationID  uuid.UUID
	SenderAvatarUrl string
	SenderID        uuid.UUID
	SenderName      string
	MessageBody     string
	MessageSentAt   time.Time
}

var ErrDecodeMessage = errors.New("failed to decode pubsub message")

func DecodeMessage(msg Message) (any, error) {
	decode, ok := decoders[msg.Type]
	if !ok {
		return nil, fmt.Errorf("%w: no decoder for msg type %q", ErrDecodeMessage, msg.Type)
	}

	return decode(msg.Payload)
}

var decoders = map[string]decoder{
	DirectConversationCreatedName: decodeAs[DirectConversationCreated],
}

type decoder func(jsontext.Value) (any, error)

func decodeAs[T any](payload jsontext.Value) (any, error) {
	var t T

	err := json.Unmarshal(payload, &t)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecodeMessage, err)
	}

	return t, nil
}
