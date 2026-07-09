package message

import (
	"messenger/messenger/internal/domain"
	"messenger/messenger/internal/domain/file"
)

type Message struct {
	*domain.BaseAggregate

	id               string
	body             string
	senderID         string
	chatID           string
	recipientID      string
	replyToMessageID string
	attachments      []*file.File
}

func Send(
	messageID string,
	chatID string,
	senderID string,
	recipientID string,
	body string,
	replyToMessageID string,
) *Message {
	m := new(Message)
	m.id = messageID
	m.SetBody(body)
	m.SetSenderID(senderID)
	m.SetChatID(chatID)
	m.SetRecipientID(recipientID)
	m.SetReplyToMessageID(replyToMessageID)

	m.AddEvent(MessageSent{
		MessageID: messageID,
		ChatID:    chatID,
	})

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

func (m *Message) SetRecipientID(recipientID string) error {
	m.recipientID = recipientID

	return nil
}

func (m *Message) SetReplyToMessageID(messageID string) error {
	m.replyToMessageID = messageID

	return nil
}

func (m *Message) AddAttachment(file *file.File) error {
	m.attachments = append(m.attachments, file)

	return nil
}

func (m *Message) ID() string {
	return m.id
}
