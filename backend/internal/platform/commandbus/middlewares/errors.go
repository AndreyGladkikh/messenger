package middlewares

import (
	"context"
	"messenger/messenger/internal/messaging/application/aperr"
	"messenger/messenger/internal/platform/commandbus"
)

func ErrorTranslator(next commandbus.CommandHandler) commandbus.CommandHandler {
	f := func(ctx context.Context, command commandbus.Command) (any, error) {
		response, err := next.Handle(ctx, command)
		if err != nil {
			return nil, aperr.Translate(err)
		}
		return response, nil
	}
	return commandbus.HandlerFunc(f)
}
