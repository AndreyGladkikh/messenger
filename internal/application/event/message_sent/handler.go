package message_sent

import (
	"context"
	"messenger/messenger/internal/domain/message"
)

type Handler struct {

}

func (h *Handler) Handle(ctx context.Context, event message.MessageCreated) error {
	return nil
}