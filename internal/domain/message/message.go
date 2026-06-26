package message

import (
	"messenger/messenger/internal/domain/file"
	"messenger/messenger/internal/domain/event"
)

type Message struct {
	id string
	body string
	senderID string
	recipientID string
	attachments []*file.File
}

func NewMessage(
	id string,
	body string,
	senderID string,
	recipientID string,
) *Message {
	m := new(Message)
	m.id = id
	m.SetBody(body)
	m.SetSenderID(senderID)
	m.SetRecipientID(recipientID)

	event.Publisher().Publish(MessageCreated{
		MessageID: id,
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

func (m *Message) SetRecipientID(recipientID string) error {
	m.recipientID = recipientID
	
	return nil
}

func (m *Message) AddAttachment(file *file.File) error {
	m.attachments = append(m.attachments, file)
	
	return nil
}
