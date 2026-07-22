package di

import (
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/platform/event"
)

func BuildEventBusForProcessor() (*event.Bus, error) {
	bus := event.NewBus()

	bus.Register(message.MessageSent{}, nil)

	return bus, nil
}