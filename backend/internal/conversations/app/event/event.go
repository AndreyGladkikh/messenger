package event

import (
	"messenger/messenger/internal/shared/app/event"
	"uuid"
)

const DirectConversationWithFirstMessageCreatedEventName = "conversations.DirectConversationWithFirstMessageCreated"

var _ event.Event = DirectConversationWithFirstMessageCreated{}

type DirectConversationWithFirstMessageCreated struct {
	ConversationID uuid.UUID `json:"conversationId"`
	SenderID  uuid.UUID `json:"senderId"`
	RecipientID  uuid.UUID `json:"recipientId"`
	MessageID  uuid.UUID `json:"messageId"`
	MessageBody  string `json:"messageBody"`
}

func (d DirectConversationWithFirstMessageCreated) Name() string {
	return DirectConversationWithFirstMessageCreatedEventName
}
