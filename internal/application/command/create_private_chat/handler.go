package create_private_chat

import (
	"context"
	"messenger/messenger/internal/application/id"
	"messenger/messenger/internal/domain/chat"
)

type Handler struct {
	chatRepository chat.Repository
	idProvider        id.Provider
}

func NewHandler(
	chatRepository chat.Repository,
	idProvider id.Provider,
) *Handler {
	return &Handler{
		chatRepository: chatRepository,
		idProvider:        idProvider,
	}
}

func (h *Handler) Handle(ctx context.Context, command *Command) (response any, err error) {
	exists, err := h.chatRepository.PrivateChatExists(ctx, command.InitiatorID, command.ChatWithUserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, chat.ErrPrivateChatExists
	}

	message := chat.Create(
		h.idProvider.ID(),
		chat.PrivateChat,
	)

	err = h.chatRepository.Add(ctx, message)

	participants := []string{
		command.InitiatorID,
		command.ChatWithUserID,
	}

	_ = participants

	return nil, err
}