package message

import (
	sharedDomain "messenger/messenger/internal/shared/domain"

	"github.com/google/uuid"
)

type Message struct {
	sharedDomain.BaseAggregate

	id               uuid.UUID
	body             string
	senderID         uuid.UUID
	chatID           uuid.UUID
	replyToMessageID uuid.UUID
}

func Send(
	chatID uuid.UUID,
	senderID uuid.UUID,
	body string,
	replyToMessageID uuid.UUID,
) *Message {
	m := new(Message)
	m.id = uuid.New()
	m.SetBody(body)
	m.SetSenderID(senderID)
	m.SetChatID(chatID)
	m.SetReplyToMessageID(replyToMessageID)

	m.AddEvent(MessageSent{
		MessageID: m.ID(),
		ChatID:    chatID,
	})

	return m
}

func Rehydrate(
	id uuid.UUID,
	chatID uuid.UUID,
	senderID uuid.UUID,
	body string,
	replyToMessageID uuid.UUID,
) *Message {
	m := new(Message)
	m.id = id
	m.body = body
	m.senderID = senderID
	m.chatID = chatID
	m.replyToMessageID = replyToMessageID

	return m
}

func (m *Message) SetBody(body string) error {
	m.body = body

	return nil
}

func (m *Message) SetSenderID(senderID uuid.UUID) error {
	m.senderID = senderID

	return nil
}

func (m *Message) SetChatID(chatID uuid.UUID) error {
	m.chatID = chatID

	return nil
}

func (m *Message) SetReplyToMessageID(messageID uuid.UUID) error {
	m.replyToMessageID = messageID

	return nil
}

func (m *Message) ID() uuid.UUID {
	return m.id
}

func (m *Message) Body() string {
	return m.body
}

func (m *Message) SenderID() uuid.UUID {
	return m.senderID
}

func (m *Message) ChatID() uuid.UUID {
	return m.chatID
}

func (m *Message) ReplyToMessageID() uuid.UUID {
	return m.replyToMessageID
}
