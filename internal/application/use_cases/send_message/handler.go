package send_message

import (
	"context"
	"messenger/messenger/internal/application/id"
	"messenger/messenger/internal/domain/message"
)

type Handler struct {
	messageRepository message.Repository
	idProvider id.Provider
}

func (h *Handler) Handle(ctx context.Context, command *Command) (response any, err error) {
	chatID := command.ChatID
	if chatID == "" {
		chatID = ""
	}

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
