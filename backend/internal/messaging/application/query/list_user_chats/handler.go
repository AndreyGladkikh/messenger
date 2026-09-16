package list_user_chats

import (
	"context"
)

type Handler struct {

}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, query *Query) (ChatListView, error) {
	return ChatListView{
		"332db6c9-d436-450f-a80e-04b64b1cb004",
		"ae16a9b4-a932-4783-b929-b7148b71ebab",
	}, nil
}