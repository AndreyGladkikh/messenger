package command

import (
	"context"
	"messenger/messenger/internal/conversation/application/service"
	"messenger/messenger/internal/conversation/domain/model/participant"
)

type SendFirstDirectMessage struct {
	Sender    participant.Participant
	Recipient participant.Participant
	Body      string
}

type SendFirstDirectMessageHandler struct {
	directConversationResolver *service.DirectConversationResolver
	messageSender              *service.MessageSender
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
