package di

import (
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/platform/event"
)

func BuildEventBusForProcessor(
	handlers *event.HandlerRegistry,
) (*event.Bus, error) {
	bus := event.NewBus()

	event.RegisterHandler(bus, message.MessageSent{}, handlers.NotifyChatParticipantsHandler)

	return bus, nil
}