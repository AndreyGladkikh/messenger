package http_server

import (
	"bytes"
	"net/http"
)

type Message struct {
	Body             string
	ChatID           string
	ReplyToMessageID string
	Attachments      []*bytes.Buffer
}

func (m *Message) Bind(r *http.Request) error {
	return nil
}

type SendMessageRequest struct {
	*Message
}
