package di

import (
	"messenger/messenger/internal/application/event/message_sent"
	"messenger/messenger/internal/domain"
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/platform/event"
)

func BuildEventBusForProcessor() (*event.Bus[domain.Event], error) {
	bus := event.NewBus[domain.Event]()

	bus.Register(message.MessageSent{}, &message_sent.NotifyChatParticipantsHandler{})

	return bus, nil
}