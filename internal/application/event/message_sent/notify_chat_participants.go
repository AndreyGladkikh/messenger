package message_sent

import (
	"context"
	"messenger/messenger/internal/domain/message"
)

type NotifyChatParticipantsHandler struct {

}

func (h *NotifyChatParticipantsHandler) Handle(ctx context.Context, event message.MessageSent) error {
	return nil
}

func (h *NotifyChatParticipantsHandler) Name() string {
	return "NotifyChatParticipantsHandler"
}