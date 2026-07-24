package middlewares

import (
	"context"
	"fmt"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/application/command"
	"messenger/messenger/internal/platform/command_bus"
	"messenger/messenger/internal/platform/event"
	"messenger/messenger/internal/platform/uow"
)

type TransactionMiddlewareContainer struct {
	txManager *transaction.Manager
	eventStorage *event.EventStorage
}

func NewTransactionMiddlewareContainer(
	txManager *transaction.Manager,
	eventStorage *event.EventStorage,
) *TransactionMiddlewareContainer {
	return &TransactionMiddlewareContainer{
		txManager: txManager,
		eventStorage: eventStorage,
	}
}

func (c *TransactionMiddlewareContainer) Middleware(next command_bus.Handler) command_bus.Handler {
	f := func(ctx context.Context, command command.Command) (response any, err error) {
		ctx = uow.NewContext(ctx, uow.New())

		err = c.txManager.WithTransaction(ctx, func(ctx context.Context) error {
			response, err = next.Handle(ctx, command)
			if err != nil {
				return err
			}
			
			if err := c.storeEvents(ctx); err != nil {
				return err
			}

			return nil
		})

		return
	}
	return command_bus.HandlerFunc(f)
}

func (c *TransactionMiddlewareContainer) storeEvents(ctx context.Context) error {
	if uow, ok := uow.FromContext(ctx); ok {
		for _, aggregate := range uow.Aggregates() {
			for _, event := range aggregate.PullEvents() {
				if err := c.eventStorage.Add(ctx, event); err != nil {
					return fmt.Errorf("failed to store event: %w", err)
				}
			}
		}
	}
	return nil
}