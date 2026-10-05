package di

import (
	"messenger/messenger/internal/identity/application/command/login"
	"messenger/messenger/internal/identity/application/command/refresh"
	"messenger/messenger/internal/identity/application/command/register"
	"messenger/messenger/internal/platform/commandbus"
	"messenger/messenger/internal/platform/commandbus/middlewares"
	"messenger/messenger/internal/platform/db/transaction"
	"messenger/messenger/internal/platform/event"
	"messenger/messenger/internal/platform/logger"
	"messenger/messenger/internal/platform/transport/http"
)

type Api struct {
	Logger     *logger.Logger
	HttpServer *http.Server
}

func NewApi(
	httpServer *http.Server,
	logger *logger.Logger,
) *Api {
	return &Api{
		HttpServer: httpServer,
		Logger:     logger,
	}
}

func BuildCommandBus(
	txManager *transaction.Manager,
	logger *logger.Logger,
	eventService *event.EventService,
	registerUserHandler *register.Handler,
	loginHandler *login.Handler,
	refreshSessionHandler *refresh.Handler,
) *commandbus.Bus {
	bus := commandbus.NewBus()

	loggingMiddlewareContainer := middlewares.NewLoggerMiddlewareContainer(logger)
	txMiddlewareContainer := middlewares.NewTransactionMiddlewareContainer(txManager, eventService)

	bus.Use(txMiddlewareContainer.Middleware)
	bus.Use(middlewares.Recoverer)
	bus.Use(loggingMiddlewareContainer.Middleware)
	bus.Use(middlewares.ErrorMapping)

	bus.RegisterHandler(registerUserHandler)
	bus.RegisterHandler(loginHandler)
	bus.RegisterHandler(refreshSessionHandler)

	return bus
}
