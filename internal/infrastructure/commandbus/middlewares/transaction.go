package middlewares

import (
	"context"
	"fmt"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/application/command"
	"messenger/messenger/internal/infrastructure/commandbus"
	"messenger/messenger/internal/infrastructure/event"
	"messenger/messenger/internal/infrastructure/uow"
)

type TransactionMiddlewareContainer struct {
	txManager    *transaction.Manager
	eventService *event.EventService
}

func NewTransactionMiddlewareContainer(
	txManager *transaction.Manager,
	eventService *event.EventService,
) *TransactionMiddlewareContainer {
	return &TransactionMiddlewareContainer{
		txManager:    txManager,
		eventService: eventService,
	}
}

func (c *TransactionMiddlewareContainer) Middleware(next commandbus.Handler) commandbus.Handler {
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
	return commandbus.HandlerFunc(f)
}

func (c *TransactionMiddlewareContainer) storeEvents(ctx context.Context) error {
	if uow, ok := uow.FromContext(ctx); ok {
		for _, aggregate := range uow.Aggregates() {
			for _, event := range aggregate.PullEvents() {
				if err := c.eventService.PutEventToOutbox(ctx, event); err != nil {
					return fmt.Errorf("failed to store event: %w", err)
				}
			}
		}
	}
	return nil
}
