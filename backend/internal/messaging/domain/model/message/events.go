package message

import (
	"messenger/messenger/internal/shared/domain"
	"time"
	"uuid"
)

const MessageSentEventName = "MessageSent"

var _ domain.Event = MessageSent{}

type MessageSent struct {
	MessageID      uuid.UUID
	ConversationID uuid.UUID
	Body           string
	SentAt         time.Time
}

func (m MessageSent) Name() string {
	return MessageSentEventName
}
