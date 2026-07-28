package create_chat

import (
	"context"
	"messenger/messenger/internal/application/id"
	"messenger/messenger/internal/domain/chat"
)

type Handler struct {
	chatRepository chat.Repository
	idProvider     id.Provider
}

func NewHandler(
	chatRepository chat.Repository,
	idProvider id.Provider,
) *Handler {
	return &Handler{
		chatRepository: chatRepository,
		idProvider:     idProvider,
	}
}

func (h *Handler) Handle(ctx context.Context, command *Command) (response any, err error) {
	chatType, err := chat.CreateType(command.Type)
	if err != nil {
		return nil, err
	}

	message := chat.Create(
		h.idProvider.ID(),
		chatType,
	)

	err = h.chatRepository.Add(ctx, message)

	return nil, err
}
