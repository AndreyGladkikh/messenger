package middlewares

import (
	"context"
	"messenger/messenger/internal/platform/commandbus"
	"messenger/messenger/internal/shared/application/apperr"
)

func ErrorMapping(next commandbus.CommandHandler) commandbus.CommandHandler {
	return commandbus.HandlerFunc(func(ctx context.Context, command commandbus.Command) (any, error) {
		response, err := next.Handle(ctx, command)
		return response, apperr.Map(err)
	})
}
