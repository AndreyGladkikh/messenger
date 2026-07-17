package di

import (
	"context"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/application/command/send_message"
	"messenger/messenger/internal/platform/command_bus"
	"messenger/messenger/internal/platform/command_bus/middlewares"
)

func BuildCommandBusForApi(
	txManager *transaction.Manager,
	sendMessageHandler *send_message.Handler,
) *command_bus.Bus {
	bus := command_bus.NewBus()

	txMiddlewareContainer := middlewares.NewTransactionMiddlewareContainer(txManager)
	bus.Use(txMiddlewareContainer.Middleware)
	bus.Use(middlewares.ErrorTranslator)

	bus.Register(&send_message.Command{}, command_bus.HandlerFunc(func(ctx context.Context, cmd command_bus.Command) (any, error) {
		return sendMessageHandler.Handle(ctx, cmd.(*send_message.Command))
	}))

	return bus
}