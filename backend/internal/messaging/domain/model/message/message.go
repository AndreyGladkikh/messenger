package message

import (
	"messenger/messenger/internal/shared/domain"
	"time"
	"uuid"
)

type Message struct {
	domain.AggregateRoot

	id             uuid.UUID
	SenderID       uuid.UUID
	ConversationID uuid.UUID
	Body           string
	CreatedAt      time.Time
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
		SenderID:       senderID,
		ConversationID: conversationID,
		Body:           body,
		CreatedAt:      time.Now(),
	}
	msg.AddEvent(MessageSent{
		MessageID:      msg.id,
		ConversationID: msg.ConversationID,
		Body:           msg.Body,
		SentAt: msg.CreatedAt,
	})
	return msg
}
