package send_message

import (
	"context"
	"messenger/messenger/internal/application/id"
	"messenger/messenger/internal/domain/message"
)

type Handler struct {
	messageRepository message.Repository
	idProvider        id.Provider
}

func NewHandler(
	messageRepository message.Repository,
	idProvider id.Provider,
) *Handler {
	return &Handler{
		messageRepository: messageRepository,
		idProvider:        idProvider,
	}
}

func (h *Handler) Handle(ctx context.Context, command *Command) (response any, err error) {
	message := message.Send(
		h.idProvider.ID(),
		command.ChatID,
		command.SenderID,
		command.MessageBody,
		command.ReplyToMessageID,
	)

	err = h.messageRepository.Add(ctx, message)

	return nil, err
}
