package di

import (
	idCommand "messenger/messenger/internal/identity/app/command"
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
	registerUserHandler *idCommand.RegisterUserHandler,
	loginUserHandler *idCommand.LoginUserHandler,
	refreshSessionHandler *idCommand.RefreshUserTokensHandler,
) *commandbus.Bus {
	bus := commandbus.NewBus()

	loggingMiddlewareContainer := middlewares.NewLoggerMiddlewareContainer(logger)
	txMiddlewareContainer := middlewares.NewTransactionMiddlewareContainer(txManager, eventService)

	bus.Use(txMiddlewareContainer.Middleware)
	bus.Use(middlewares.Recoverer)
	bus.Use(loggingMiddlewareContainer.Middleware)
	bus.Use(middlewares.ErrorMapping)

	bus.RegisterHandler(registerUserHandler)
	bus.RegisterHandler(loginUserHandler)
	bus.RegisterHandler(refreshSessionHandler)

	return bus
}
