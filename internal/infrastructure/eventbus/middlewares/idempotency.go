package middlewares

import (
	"context"
	"messenger/messenger/internal/infrastructure/event"
)

func Idempotency(cache any) event.Middleware {
	return func(next event.Handler) event.Handler {
		f := func(ctx context.Context, envelope *event.Envelope) error {
			return nil
		}
		return event.HandlerFunc(f)
	}
}
