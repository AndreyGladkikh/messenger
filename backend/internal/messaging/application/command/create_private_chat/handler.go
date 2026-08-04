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
	privateChatPair, err := chat.NewPrivateChatPair(command.InitiatorID, command.ChatWithUserID)
	if err != nil {
		return nil, err
	}

	chat := chat.Create(
		chat.PrivateChat,
	)

	err = h.chatRepository.AddPrivate(ctx, chat, privateChatPair)
	if err != nil {
		return nil, err
	}

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

	err = h.chatParticipantRepository.Add(ctx, participants)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
