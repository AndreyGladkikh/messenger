package command

import (
	"context"
	"messenger/messenger/internal/conversations/app/service"
	"uuid"
)

type SendFirstDirectMessage struct {
	Sender    uuid.UUID
	Recipient uuid.UUID
	Body      string
}

func (cmd SendFirstDirectMessage) IsCommand() {}

func (cmd SendFirstDirectMessage) Name() string {
	return "SendFirstDirectMessage"
}

type SendFirstDirectMessageHandler struct {
	directConversationResolver *service.DirectConversationResolver
	messageSender              *service.MessageSender
}

func NewSendFirstDirectMessageHandler(
	directConversationResolver *service.DirectConversationResolver,
	messageSender *service.MessageSender,
) *SendFirstDirectMessageHandler {
	return &SendFirstDirectMessageHandler{
		directConversationResolver: directConversationResolver,
		messageSender:              messageSender,
	}
}

func (h *SendFirstDirectMessageHandler) Handle(ctx context.Context, cmd SendFirstDirectMessage) (any, error) {
	cnv, err := h.directConversationResolver.FindOrCreate(ctx, cmd.Sender, cmd.Recipient)
	if err != nil {
		return nil, err
	}

	_, err = h.messageSender.Send(ctx, cmd.Sender, cnv, cmd.Body)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
