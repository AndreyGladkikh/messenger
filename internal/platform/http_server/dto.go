package http_server

import (
	"bytes"
)

type Chat struct {
	Type string
	Name string
}

type CreateChatRequest struct {
	Chat
}

type Message struct {
	Body             string
	ChatID           string
	ReplyToMessageID string
	Attachments      []*bytes.Buffer
}

type SendMessageRequest struct {
	Message
}
