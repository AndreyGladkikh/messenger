package get_chat_list

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, query *Query) (*ChatListView, error) {
	return &ChatListView{
		Chats: []*ChatListItemView{
			{
				ID:   uuid.New(),
				Name: "Chat 1",
				LastMessage: &LastMessage{
					Body:   "Hello, how are you?",
					SentAt: time.Now(),
				},
				UnreadCount: 2,
			},
			{
				ID:   uuid.New(),
				Name: "Chat 2",
				LastMessage: &LastMessage{
					Body:   "Hey, what's up?",
					SentAt: time.Now(),
				},
				UnreadCount: 0,
			},
		},
	}, nil
}
