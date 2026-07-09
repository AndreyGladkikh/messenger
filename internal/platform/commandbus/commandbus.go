package commandbus

import (
	"context"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/application/command"
	"messenger/messenger/internal/application/command/middlewares"
	"messenger/messenger/internal/application/use_cases/send_message"
)

func BuildCommandBus(
	txManager *transaction.Manager,
	sendMessageHandler *send_message.Handler,
) *command.Bus {
	bus := new(command.Bus)

	txMiddlewareContainer := middlewares.NewTransactionMiddlewareContainer(txManager)

	bus.Use(txMiddlewareContainer.Middleware)

	bus.Register(&send_message.Command{}, command.HandlerFunc(func(ctx context.Context, cmd command.Command) (any, error) {
		return sendMessageHandler.Handle(ctx, cmd.(*send_message.Command))
	}))

	return bus
}
