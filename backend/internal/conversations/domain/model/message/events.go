package message

import "uuid"

const MessageSentEventName = "MessageSent"

type MessageSent struct {
	MessageID      uuid.UUID
	ConversationID uuid.UUID
	Body           string
}

func (m MessageSent) IsEvent() {}

func (m MessageSent) Name() string {
	return MessageSentEventName
}
