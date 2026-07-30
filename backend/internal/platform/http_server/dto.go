package http_server

import (
	"bytes"

	"github.com/google/uuid"
)

type Credentials struct {
	Login string `json:"login"`
	Password string `json:"password"`
}

type RegisterUserRequest struct {
	Credentials
}

type Login struct {
	Credentials
}

type Chat struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type CreatePrivateChatRequest struct {
	ChatWithUserID uuid.UUID `json:"chatWithUserId"`
}

type Message struct {
	Body             string `json:"body"`
	ChatID           uuid.UUID `json:"chatId"`
	ReplyToMessageID uuid.UUID `json:"replyToMessageId"`
	Attachments      []*bytes.Buffer `json:"attachments"`
}

type SendMessageRequest struct {
	Message
}
