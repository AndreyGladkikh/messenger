package command

import (
	"context"
	"messenger/messenger/internal/messaging/domain/model/group_conversation"
	"messenger/messenger/internal/shared/app/command"
	"uuid"
)

const CreateGroupConversationCommandName = "CreateGroupConversation"

var _ command.Command = CreateGroupConversation{}

type CreateGroupConversation struct {
	Title     string
	AvatarUrl *string
}

func (c CreateGroupConversation) IsCommand() {}

func (c CreateGroupConversation) Name() string {
	return CreateGroupConversationCommandName
}

type CreateGroupConversationHandler struct {
	groupsRepo group_conversation.Repository
}

func NewCreateGroupConversationHandler(
	groupsRepo group_conversation.Repository,
) *CreateGroupConversationHandler {
	return &CreateGroupConversationHandler{
		groupsRepo: groupsRepo,
	}
}

func (h *CreateGroupConversationHandler) Handle(ctx context.Context, cmd CreateGroupConversation) (uuid.UUID, error) {
	cnv := group_conversation.New(
		cmd.Title,
		cmd.AvatarUrl,
	)

	err := h.groupsRepo.Save(ctx, cnv)
	if err != nil {
		return uuid.UUID{}, err
	}

	return cnv.ID(), nil
}
