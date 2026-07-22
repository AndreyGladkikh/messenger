package middlewares

import (
	"context"
	"messenger/messenger/internal/application/aperr"
	"messenger/messenger/internal/application/command"
	"messenger/messenger/internal/platform/command_bus"
)

func ErrorTranslator(next command_bus.Handler) command_bus.Handler {
	return command_bus.HandlerFunc(func(ctx context.Context, command command.Command) (any, error) {
		response, err := next.Handle(ctx, command)
		if err != nil {
			return nil, aperr.Translate(err)
		}
		return response, nil
	})
}