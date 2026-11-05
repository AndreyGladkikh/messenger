package command_middleware

import (
	"context"
	"fmt"
	"messenger/messenger/internal/platform/event"
	"messenger/messenger/internal/platform/transactions"
	"messenger/messenger/internal/platform/uow"
	"messenger/messenger/internal/shared/app/command"
)

type Transaction struct {
	txManager    *transactions.Manager
	eventService *event.EventService
}

func NewTransaction(
	txManager *transactions.Manager,
	eventService *event.EventService,
) *Transaction {
	return &Transaction{
		txManager:    txManager,
		eventService: eventService,
	}
}

func (c *Transaction) Wrap[C command.Command, R any](next command.Handler[C, R]) command.Handler[C ,R] {
	return command.HandlerFunc[C, R](func(ctx context.Context, command C) (response R, err error) {
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

// todo translate events to integration and domain events
func (c *Transaction) storeEvents(ctx context.Context) error {
	if uow, ok := uow.FromContext(ctx); ok {
		for _, aggregate := range uow.Aggregates() {
			for _, event := range aggregate.PullEvents() {
				if err := c.eventService.PutEventToOutbox(ctx, "integration", event); err != nil {
					return fmt.Errorf("failed to store event: %w", err)
				}
			}
		}
		for _, evnt := range uow.Events() {
			if err := c.eventService.PutEventToOutbox(ctx, "integration", evnt); err != nil {
				return fmt.Errorf("failed to store event: %w", err)
			}
		}
	}
	return nil
}
