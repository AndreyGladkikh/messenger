package create_private_chat

import (
	"context"
	"messenger/messenger/internal/application/id"
	"messenger/messenger/internal/domain/chat"
	"messenger/messenger/internal/domain/chat_participant"
)

type Handler struct {
	chatRepository chat.Repository
	chatParticipantRepository chat_participant.Repository
	idProvider        id.Provider
}

func NewHandler(
	chatRepository chat.Repository,
	chatParticipantRepository chat_participant.Repository,
	idProvider id.Provider,
) *Handler {
	return &Handler{
		chatRepository: chatRepository,
		chatParticipantRepository: chatParticipantRepository,
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

	chatID := h.idProvider.ID()
	chat := chat.Create(
		chatID,
		chat.PrivateChat,
	)

	err = h.chatRepository.Add(ctx, chat)

	participants := make([]*chat_participant.ChatParticipant, 0, 2)
	for _, participantID := range []string{
		command.InitiatorID,
		command.ChatWithUserID,
	} {
		participants = append(participants, chat_participant.AddToChat(
			h.idProvider.ID(),
			chatID,
			participantID,
			chat_participant.RoleParticipant,
		))
	}

	h.chatParticipantRepository.Add(ctx, participants)

	return nil, err
}