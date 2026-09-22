package get_chat_list

import (
	"time"

	"github.com/google/uuid"
)

type ChatListView struct {
	Chats []*ChatListItemView `json:"chats"`
}

type ChatListItemView struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	LastMessage *LastMessage `json:"lastMessage"`
	UnreadCount int          `json:"unreadCount"`
}

type LastMessage struct {
	Body   string    `json:"body"`
	SentAt time.Time `json:"sentAt"`
}
