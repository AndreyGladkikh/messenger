package di

import (
	"messenger/messenger/internal/adapters/out/logger"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/application/command/create_private_chat"
	"messenger/messenger/internal/application/command/send_message"
	"messenger/messenger/internal/platform/command_bus"
	"messenger/messenger/internal/platform/command_bus/middlewares"
	"messenger/messenger/internal/platform/event"
	"messenger/messenger/internal/platform/http_server"
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
) *command_bus.Bus {
	bus := command_bus.NewBus()

	loggingMiddlewareContainer := middlewares.NewLoggerMiddlewareContainer(logger)
	txMiddlewareContainer := middlewares.NewTransactionMiddlewareContainer(txManager, eventStorage)
	
	bus.Use(txMiddlewareContainer.Middleware)
	bus.Use(middlewares.Recoverer)
	bus.Use(loggingMiddlewareContainer.Middleware)
	bus.Use(middlewares.ErrorTranslator)

	command_bus.RegisterHandler(bus, &send_message.Command{}, sendMessageHandler)
	command_bus.RegisterHandler(bus, &create_private_chat.Command{}, createPrivateChatHandler)

	return bus
}