package middlewares

import (
	"context"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/platform/command_bus"
	"messenger/messenger/internal/platform/uow"
)

// type txManager interface {
// 	WithTransaction(ctx context.Context, fn func(context.Context) error) error
// }

type TransactionMiddlewareContainer struct {
	txManager *transaction.Manager
}

func NewTransactionMiddlewareContainer(txManager *transaction.Manager) *TransactionMiddlewareContainer {
	return &TransactionMiddlewareContainer{
		txManager: txManager,
	}
}

func (c *TransactionMiddlewareContainer) Middleware(next command_bus.Handler) command_bus.Handler {
	return command_bus.HandlerFunc(func(ctx context.Context, command command_bus.Command) (response any, err error) {
		uowo := uow.New()
		ctx = uow.NewContext(ctx, uowo)

		err = c.txManager.WithTransaction(ctx, func(ctx context.Context) error {
			response, err = next.Handle(ctx, command)
			return err
		})

		for _, aggregate := range uowo.Aggregates() {
			_ = aggregate.PullEvents()
		}

		return
	})
}
