package di

import (
	"messenger/messenger/internal/platform/event"
)

func NewEventHandlerRegistry(
) *event.HandlerRegistry {
	r := event.NewRegistry()

	return r
}
