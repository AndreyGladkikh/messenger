package message_sent

import (
	"context"
	"fmt"
	"messenger/messenger/internal/domain/message"
)

type NotifyChatParticipantsHandler struct {
}

func NewNotifyChatParticipantsHandler() *NotifyChatParticipantsHandler {
	return &NotifyChatParticipantsHandler{}
}

func (h *NotifyChatParticipantsHandler) Handle(ctx context.Context, event message.MessageSent) error {
	fmt.Printf("event handler %s executed", h.Name())
	return nil
}

func (h *NotifyChatParticipantsHandler) Name() string {
	return "event_handler.notify_chat_participants.v1"
}
