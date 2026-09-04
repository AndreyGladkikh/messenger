package message_sent

import (
	"context"
	"messenger/messenger/internal/messaging/application/notifier"
	"messenger/messenger/internal/messaging/domain/message"
)

type NotifyChatParticipantsHandler struct {
	notifier notifier.MessageSentNotifier
}

func NewNotifyChatParticipantsHandler(
	notifier notifier.MessageSentNotifier,
) *NotifyChatParticipantsHandler {
	return &NotifyChatParticipantsHandler{
		notifier: notifier,
	}
}

func (h *NotifyChatParticipantsHandler) Handle(ctx context.Context, event message.MessageSent) error {
	err := h.notifier.Notify(ctx, event)
	if err != nil {
		return err
	}
	return nil
}

func (h *NotifyChatParticipantsHandler) Name() string {
	return "event_handler.notify_chat_participants.v1"
}
