package message

import (
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

var _ domain.Aggregate = (*Message)(nil)

type Message struct {
	domain.AggregateRoot

	id             uuid.UUID
	SenderID         uuid.UUID
	ConversationID uuid.UUID
	Body           string
}

func (m *Message) ID() uuid.UUID {
	return m.id
}

func New(
	senderID uuid.UUID,
	conversationID uuid.UUID,
	body string,
) *Message {
	msg := &Message{
		id:             uuid.New(),
		SenderID:         senderID,
		ConversationID: conversationID,
		Body:           body,
	}
	msg.AddEvent(MessageSent{
		MessageID:      msg.id,
		ConversationID: msg.ConversationID,
		Body:           msg.Body,
	})
	return msg
}
