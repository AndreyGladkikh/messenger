package create_private_chat

import (
	"context"
	"messenger/messenger/internal/messaging/domain/chat"
	"messenger/messenger/internal/messaging/domain/chat_participant"

	"github.com/google/uuid"
)

type Handler struct {
	chatRepository            chat.Repository
	chatParticipantRepository chat_participant.Repository
}

func NewHandler(
	chatRepository chat.Repository,
	chatParticipantRepository chat_participant.Repository,
) *Handler {
	return &Handler{
		chatRepository:            chatRepository,
		chatParticipantRepository: chatParticipantRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, command *Command) (response any, err error) {
	exists, err := h.chatRepository.PrivateChatExists(ctx, command.InitiatorID, command.ChatWithUserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, chat.ErrPrivateChatAlreadyExists
	}

	chat := chat.Create(
		chat.PrivateChat,
	)

	err = h.chatRepository.Add(ctx, chat)

	participants := make([]*chat_participant.ChatParticipant, 0, 2)
	for _, participantID := range []uuid.UUID{
		command.InitiatorID,
		command.ChatWithUserID,
	} {
		participants = append(participants, chat_participant.AddToChat(
			chat.ID(),
			participantID,
			chat_participant.RoleParticipant,
		))
	}

	h.chatParticipantRepository.Add(ctx, participants)

	return nil, err
}
