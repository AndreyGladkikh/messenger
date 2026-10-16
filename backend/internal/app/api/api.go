package api

import (
	"messenger/messenger/internal/app"
	platformHttp "messenger/messenger/internal/platform/transport/http"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/wire"
)

type App struct {
	HttpHandler http.Handler
}

func NewApp(
	httpHandler http.Handler,
) *App {
	return &App{
		HttpHandler: httpHandler,
	}
}

var AppSet = wire.NewSet(
	NewApp,
	app.AppSet,
	platformHttp.NewRouter,
	platformHttp.NewApiHandler,
	platformHttp.NewErrorHandler,
	platformHttp.NewWebsocketHandler,

	wire.Bind(new(http.Handler), new(*chi.Mux)),
	wire.Bind(new(platformHttp.StrictServerInterface), new(*platformHttp.ApiHandler)),
)
