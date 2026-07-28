package message

import (
	"messenger/messenger/internal/domain"
)

type Message struct {
	domain.BaseAggregate

	id               string
	body             string
	senderID         string
	chatID           string
	replyToMessageID string
}

func Send(
	messageID string,
	chatID string,
	senderID string,
	body string,
	replyToMessageID string,
) *Message {
	m := new(Message)
	m.id = messageID
	m.SetBody(body)
	m.SetSenderID(senderID)
	m.SetChatID(chatID)
	m.SetReplyToMessageID(replyToMessageID)

	m.AddEvent(MessageSent{
		MessageID: messageID,
		ChatID:    chatID,
	})

	return m
}

func Rehydrate(
	id string,
	chatID string,
	senderID string,
	body string,
	replyToMessageID string,
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

func (m *Message) SetSenderID(senderID string) error {
	m.senderID = senderID

	return nil
}

func (m *Message) SetChatID(chatID string) error {
	m.chatID = chatID

	return nil
}

func (m *Message) SetReplyToMessageID(messageID string) error {
	m.replyToMessageID = messageID

	return nil
}

func (m *Message) ID() string {
	return m.id
}

func (m *Message) Body() string {
	return m.body
}

func (m *Message) SenderID() string {
	return m.senderID
}

func (m *Message) ChatID() string {
	return m.chatID
}

func (m *Message) ReplyToMessageID() string {
	return m.replyToMessageID
}
