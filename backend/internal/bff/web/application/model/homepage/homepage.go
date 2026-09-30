package homepage

import (
	"time"

	"github.com/google/uuid"
)

type HomePage struct {
	Chats []ChatItem
}

type ChatItem struct {
	ID          uuid.UUID
	Name        string
	LastMessage *LastMessage
	UnreadCount int
}

type LastMessage struct {
	Body   string
	SentAt time.Time
}
