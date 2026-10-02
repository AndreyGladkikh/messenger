package command

import (
	"context"
	"messenger/messenger/internal/conversation/application/service"
	"messenger/messenger/internal/conversation/domain/model/conversation"
	"messenger/messenger/internal/conversation/domain/model/participant"
)

type SendFirstDirectMessage struct {
	Sender    participant.Participant
	Recipient participant.Participant
	Body      string
}

type SendFirstDirectMessageHandler struct {
	directConversationResolver *service.DirectConversationResolver
	messageSender              service.MessageSender
	// directConversationRepository direct_conversation.Repository
	// messageSender     *message.Sender
	// messageRepository message.Repository
}

func (h *SendFirstDirectMessageHandler) Handle(ctx context.Context, c SendFirstDirectMessage) error {
	participants, err := conversation.NewDirectConversationParticipantsPair(c.Sender, c.Recipient)
	if err != nil {
		return err
	}

	cnv, err := h.directConversationResolver.Create(ctx, participants)
	if err != nil {
		return err
	}

	// msg, err := h.messageSender.Send(ctx, c.Sender, conv, c.Body)
	if err != nil {
		return err
	}

	// err = h.messageRepository.Save(ctx, msg)
	if err != nil {
		return err
	}

	_ = cnv
	return nil
}
