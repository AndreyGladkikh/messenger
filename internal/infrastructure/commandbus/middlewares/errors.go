package middlewares

import (
	"context"
	"messenger/messenger/internal/application/aperr"
	"messenger/messenger/internal/application/command"
	"messenger/messenger/internal/infrastructure/commandbus"
)

func ErrorTranslator(next commandbus.Handler) commandbus.Handler {
	f := func(ctx context.Context, command command.Command) (any, error) {
		response, err := next.Handle(ctx, command)
		if err != nil {
			return nil, aperr.Translate(err)
		}
		return response, nil
	}
	return commandbus.HandlerFunc(f)
}
