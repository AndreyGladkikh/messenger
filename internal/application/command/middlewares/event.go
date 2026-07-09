package middlewares

import (
	"context"
	"messenger/messenger/internal/application/command"
	"messenger/messenger/internal/domain/event"
)

func EventListener(next command.Handler) command.Handler {
	return command.HandlerFunc(func(ctx context.Context, command command.Command) (response any, err error) {
		publisher := event.Publisher()
		publisher.Subscribe(eventHandler)

		response, err = next.Handle(ctx, command)

		return 
	})
}

func eventHandler(event event.DomainEvent) error {
	return nil
}