package middlewares

import (
	"context"
	"messenger/messenger/internal/application/command"
)

type txManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type TransactionMiddlewareContainer struct {
	txManager txManager
}

func NewTransactionMiddlewareContainer(txManager txManager) *TransactionMiddlewareContainer {
	return &TransactionMiddlewareContainer{
		txManager: txManager,
	}
}

func (c *TransactionMiddlewareContainer) Middleware(next command.Handler) command.Handler {
	return command.HandlerFunc(func(ctx context.Context, command command.Command) (any, error) {
		var response any
		err := c.txManager.WithTransaction(ctx, func(ctx context.Context) error {
			var err error
			response, err = next.Handle(ctx, command)
			return err
		})
		return response, err
	})
}
