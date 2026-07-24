package di

import (
	"messenger/messenger/internal/application/event/message_sent"
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/infrastructure/event"
)

func BuildEventBusForProcessor(
	notifyChatParticipantsHandler *message_sent.NotifyChatParticipantsHandler,
) (*event.Bus, error) {
	bus := event.NewBus()

	event.RegisterHandler(bus, message.MessageSent{}, notifyChatParticipantsHandler)

	return bus, nil
}
