package message

import (
	"messenger/messenger/internal/conversation/domain/model/participant"
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

var _ domain.Aggregate = (*Message)(nil)

type Message struct {
	domain.AggregateRoot

	id             uuid.UUID
	Sender         participant.Participant
	ConversationID uuid.UUID
	Body           string
}

func (m *Message) ID() uuid.UUID {
	return m.id
}

func New(
	sender participant.Participant,
	conversationID uuid.UUID,
	body string,
) *Message {
	msg := &Message{
		id:             uuid.New(),
		Sender:         sender,
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
