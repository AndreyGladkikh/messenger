package di

import (
	"messenger/messenger/internal/messaging/application/event/message_sent"
	"messenger/messenger/internal/platform/event"
)

func NewEventHandlerRegistry(
	notifyChatParticipantsHandler *message_sent.NotifyChatParticipantsHandler,
	rebuildQueryModelHandler *message_sent.RebuildQueryModelHandler,
) *event.HandlerRegistry {
	r := event.NewRegistry()

	event.RegisterEventHandler(r, notifyChatParticipantsHandler)
	event.RegisterEventHandler(r, rebuildQueryModelHandler)

	return r
}
