package http_server

import (
	"bytes"
)

type Chat struct {
	Type string
	Name string
}

type CreatePrivateChatRequest struct {
	ChatWithUserID string
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
