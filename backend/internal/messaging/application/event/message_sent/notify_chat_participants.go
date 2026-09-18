package message_sent

import (
	"context"
	"messenger/messenger/internal/messaging/domain/message"
	"messenger/messenger/internal/shared/application/chateventspublisher"
	"messenger/messenger/internal/shared/application/event"
)

type NotifyChatParticipantsHandler struct {
	chatEventsPublisher chateventspublisher.ChatEventsPublisher
}

func NewNotifyChatParticipantsHandler(
	chatEventsPublisher chateventspublisher.ChatEventsPublisher,
) *NotifyChatParticipantsHandler {
	return &NotifyChatParticipantsHandler{
		chatEventsPublisher: chatEventsPublisher,
	}
}

func (h *NotifyChatParticipantsHandler) Handle(ctx context.Context, e event.Envelope[message.MessageSent]) error {
	err := h.chatEventsPublisher.Publish(ctx, e.Event.ChatID, event.AsEnvelopeWithDomainEvent(e))
	if err != nil {
		return err
	}
	return nil
}

func (h *NotifyChatParticipantsHandler) Name() string {
	return "event_handler.notify_chat_participants.v1"
}
