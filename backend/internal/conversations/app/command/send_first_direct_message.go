package command

import (
	"context"
	"messenger/messenger/internal/conversations/app/event"
	"messenger/messenger/internal/conversations/app/service"
	appEvent "messenger/messenger/internal/shared/app/event"
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

type SendFirstDirectMessageResponse struct {
	ConversationID uuid.UUID
	MessageID      uuid.UUID
}

type SendFirstDirectMessageHandler struct {
	DirectConversationResolver *service.DirectConversationResolver
	MessageSender              *service.MessageSender
	EventPublisher             appEvent.Publisher
}

func (h *SendFirstDirectMessageHandler) Handle(ctx context.Context, cmd SendFirstDirectMessage) (SendFirstDirectMessageResponse, error) {
	cnv, found, err := h.DirectConversationResolver.FindOrCreate(ctx, cmd.Sender, cmd.Recipient)
	if err != nil {
		return SendFirstDirectMessageResponse{}, err
	}

	m, err := h.MessageSender.Send(ctx, cmd.Sender, cnv, cmd.Body)
	if err != nil {
		return SendFirstDirectMessageResponse{}, err
	}

	if !found {
		h.EventPublisher.Publish(ctx, event.DirectConversationWithFirstMessageCreated{
			ConversationID: cnv.ID(),
			SenderID:       cmd.Sender,
			RecipientID:    cmd.Recipient,
			MessageID:      m.ID(),
			MessageBody:    m.Body,
		})
	}

	return SendFirstDirectMessageResponse{
		ConversationID: cnv.ID(),
		MessageID:      m.ID(),
	}, nil
}
