package command

import (
	"context"
	"messenger/messenger/internal/messaging/app/event"
	"messenger/messenger/internal/messaging/app/service"
	"messenger/messenger/internal/shared/app/command"
	appEvent "messenger/messenger/internal/shared/app/event"
	"time"
	"uuid"
)

const StartDirectConversationCommandName = "StartDirectConversation"

var _ command.Command = StartDirectConversation{}

type StartDirectConversation struct {
	Sender    uuid.UUID
	Recipient uuid.UUID
	Body      string
}

func (cmd StartDirectConversation) IsCommand() {}

func (cmd StartDirectConversation) Name() string {
	return StartDirectConversationCommandName
}

type StartDirectConversationResponse struct {
	ConversationID uuid.UUID
	MessageID      uuid.UUID
	MessageSentAt  time.Time
}

type StartDirectConversationHandler struct {
	DirectConversationResolver *service.DirectConversationResolver
	MessageSender              *service.MessageSender
	EventPublisher             appEvent.Publisher
}

func NewStartDirectConversationHandler(
	directConversationResolver *service.DirectConversationResolver,
	messageSender *service.MessageSender,
	eventPublisher appEvent.Publisher,
) *StartDirectConversationHandler {
	return &StartDirectConversationHandler{
		DirectConversationResolver: directConversationResolver,
		MessageSender:              messageSender,
		EventPublisher:             eventPublisher,
	}
}

func (h *StartDirectConversationHandler) Handle(ctx context.Context, cmd StartDirectConversation) (StartDirectConversationResponse, error) {
	cnv, found, err := h.DirectConversationResolver.FindOrCreate(ctx, cmd.Sender, cmd.Recipient)
	if err != nil {
		return StartDirectConversationResponse{}, err
	}

	m, err := h.MessageSender.Send(ctx, cmd.Sender, cnv, cmd.Body)
	if err != nil {
		return StartDirectConversationResponse{}, err
	}

	if !found {
		h.EventPublisher.Publish(ctx, event.DirectConversationStarted{
			ConversationID: cnv.ID(),
			SenderID:       cmd.Sender,
			RecipientID:    cmd.Recipient,
			MessageID:      m.ID(),
			MessageBody:    m.Body,
			MessageSentAt: m.CreatedAt,
		})
	}

	return StartDirectConversationResponse{
		ConversationID: cnv.ID(),
		MessageID:      m.ID(),
		MessageSentAt:  m.CreatedAt,
	}, nil
}
