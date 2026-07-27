package di

import (
	"messenger/messenger/internal/application/event/message_sent"
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/infrastructure/event"
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