package middlewares

import (
	"context"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/platform/command_bus"
	"messenger/messenger/internal/platform/event"
	"messenger/messenger/internal/platform/uow"
)

type TransactionMiddlewareContainer struct {
	txManager *transaction.Manager
	eventStorage *event.EventStorage
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
			if err != nil {
				return err
			}
			
			c.storeEvents(ctx, uowo)

			return nil
		})

		return
	})
}

func (c *TransactionMiddlewareContainer) storeEvents(ctx context.Context, uow *uow.UnitOfWork) error {
	for _, aggregate := range uow.Aggregates() {
		for _, event := range aggregate.PullEvents() {
			if err := c.eventStorage.Add(ctx, event); err != nil {
				return err
			}
		}
	}
	return nil
}