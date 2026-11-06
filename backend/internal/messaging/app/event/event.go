package event

import (
	"messenger/messenger/internal/shared/app/event"
	"time"
	"uuid"
)

const DirectConversationStartedEventName = "conversations.DirectConversationStarted.v1"

var _ event.Event = DirectConversationStarted{}

type DirectConversationStarted struct {
	ConversationID uuid.UUID `json:"conversationId"`
	SenderID       uuid.UUID `json:"senderId"`
	RecipientID    uuid.UUID `json:"recipientId"`
	MessageID      uuid.UUID `json:"messageId"`
	MessageBody    string    `json:"messageBody"`
	MessageSentAt    time.Time    `json:"messageSentAt"`
}

func (d DirectConversationStarted) Name() string {
	return DirectConversationStartedEventName
}
