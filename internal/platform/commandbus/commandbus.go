package commandbus

import (
	"messenger/messenger/internal/application/command"
	"messenger/messenger/internal/application/command/middlewares"
)

func InitCommandBus() {
	bus := new(command.Bus)

	middlewares.NewTransactionMiddlewareContainer(txManager)

	bus.Use(middlewares.TransactionMiddlewareContainer)
}
