package message

import (
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

const MessageSentEventName = "MessageSent"

var _ domain.Event = MessageSent{}

type MessageSent struct {
	MessageID      uuid.UUID
	ConversationID uuid.UUID
	Body           string
}

func (m MessageSent) Name() string {
	return MessageSentEventName
}
