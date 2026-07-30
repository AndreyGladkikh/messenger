package create_chat

import (
	"context"
	"messenger/messenger/internal/messaging/domain/chat"
)

type Handler struct {
	chatRepository chat.Repository
}

func NewHandler(
	chatRepository chat.Repository,
) *Handler {
	return &Handler{
		chatRepository: chatRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, command *Command) (response any, err error) {
	chatType, err := chat.CreateType(command.Type)
	if err != nil {
		return nil, err
	}

	message := chat.Create(
		chatType,
	)

	err = h.chatRepository.Add(ctx, message)

	return nil, err
}
