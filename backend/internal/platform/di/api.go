package di

import (
	"messenger/messenger/internal/auth/application/command/register"
	"messenger/messenger/internal/messaging/application/command/create_private_chat"
	"messenger/messenger/internal/messaging/application/command/send_message"
	"messenger/messenger/internal/messaging/infrastructure/event"
	"messenger/messenger/internal/platform/commandbus"
	"messenger/messenger/internal/platform/commandbus/middlewares"
	"messenger/messenger/internal/platform/db/transaction"
	"messenger/messenger/internal/platform/http_server"
	"messenger/messenger/internal/platform/logger"
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
	eventService *event.EventService,
	registerUserHandler *register.Handler,
	sendMessageHandler *send_message.Handler,
	createPrivateChatHandler *create_private_chat.Handler,
) *commandbus.Bus {
	bus := commandbus.NewBus()

	loggingMiddlewareContainer := middlewares.NewLoggerMiddlewareContainer(logger)
	txMiddlewareContainer := middlewares.NewTransactionMiddlewareContainer(txManager, eventService)

	bus.Use(txMiddlewareContainer.Middleware)
	bus.Use(middlewares.Recoverer)
	bus.Use(loggingMiddlewareContainer.Middleware)

	commandbus.RegisterHandler(bus, registerUserHandler)
	commandbus.RegisterHandler(bus, sendMessageHandler)
	commandbus.RegisterHandler(bus, createPrivateChatHandler)

	return bus
}
