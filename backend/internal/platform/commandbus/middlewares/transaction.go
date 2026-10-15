package middlewares

import (
	"context"
	"fmt"
	"messenger/messenger/internal/platform/commandbus"
	"messenger/messenger/internal/platform/event"
	"messenger/messenger/internal/platform/transactions"
	"messenger/messenger/internal/platform/uow"
	"messenger/messenger/internal/shared/app/command"
)

type TransactionMiddlewareContainer struct {
	txManager    *transactions.Manager
	eventService *event.EventService
}

func NewTransactionMiddlewareContainer(
	txManager *transactions.Manager,
	eventService *event.EventService,
) *TransactionMiddlewareContainer {
	return &TransactionMiddlewareContainer{
		txManager:    txManager,
		eventService: eventService,
	}
}

func (c *TransactionMiddlewareContainer) Middleware(next commandbus.CommandHandler) commandbus.CommandHandler {
	return commandbus.HandlerFunc(func(ctx context.Context, command command.Command) (response any, err error) {
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
	})
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
