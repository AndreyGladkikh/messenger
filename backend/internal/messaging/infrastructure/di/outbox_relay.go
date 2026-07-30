package di

import (
	"messenger/messenger/internal/messaging/application/event/message_sent"
	"messenger/messenger/internal/messaging/domain/message"
	"messenger/messenger/internal/messaging/infrastructure/event"
)

func NewEventHandlerRegistry(
	notifyChatParticipantsHandler *message_sent.NotifyChatParticipantsHandler,
	rebuildQueryModelHandler *message_sent.RebuildQueryModelHandler,
) *event.HandlerRegistry {
	r := event.NewRegistry()

	event.RegisterEventHandler(r, message.MessageSentEventName, notifyChatParticipantsHandler)
	event.RegisterEventHandler(r, message.MessageSentEventName, rebuildQueryModelHandler)

	return r
}
