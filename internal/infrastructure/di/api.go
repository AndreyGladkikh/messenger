package di

import (
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/application/command/create_private_chat"
	"messenger/messenger/internal/application/command/send_message"
	"messenger/messenger/internal/infrastructure/commandbus"
	"messenger/messenger/internal/infrastructure/commandbus/middlewares"
	"messenger/messenger/internal/infrastructure/event"
	"messenger/messenger/internal/infrastructure/http_server"
	"messenger/messenger/internal/infrastructure/logger"
)

type Api struct {
	Logger     *logger.Logger
	HttpServer *http_server.Server
}

func NewApi(
	httpServer *http_server.Server,
	logger *logger.Logger,
) *Api {
	return &Api{
		HttpServer: httpServer,
		Logger:     logger,
	}
}

func BuildCommandBusForApi(
	txManager *transaction.Manager,
	logger *logger.Logger,
	eventStorage *event.EventStorage,
	sendMessageHandler *send_message.Handler,
	createPrivateChatHandler *create_private_chat.Handler,
) *commandbus.Bus {
	bus := commandbus.NewBus()

	loggingMiddlewareContainer := middlewares.NewLoggerMiddlewareContainer(logger)
	txMiddlewareContainer := middlewares.NewTransactionMiddlewareContainer(txManager, eventStorage)

	bus.Use(txMiddlewareContainer.Middleware)
	bus.Use(middlewares.Recoverer)
	bus.Use(loggingMiddlewareContainer.Middleware)
	bus.Use(middlewares.ErrorTranslator)

	commandbus.RegisterHandler(bus, &send_message.Command{}, sendMessageHandler)
	commandbus.RegisterHandler(bus, &create_private_chat.Command{}, createPrivateChatHandler)

	return bus
}
