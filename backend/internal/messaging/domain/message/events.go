package message

import "github.com/google/uuid"

const MessageSentEventName = "MessageSent"

type MessageSent struct {
	MessageID uuid.UUID
	ChatID    uuid.UUID
}

func (m MessageSent) IsEvent() {}

func (m MessageSent) Name() string {
	return MessageSentEventName
}
