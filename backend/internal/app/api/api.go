package api

import (
	"messenger/messenger/internal/app"
	appHttp "messenger/messenger/internal/platform/transport/http"
	"messenger/messenger/internal/platform/transport/websocket"
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
	appHttp.NewRouter,
	appHttp.NewApiHandler,
	appHttp.NewErrorHandler,
	wire.Struct(new(websocket.Handler), "*"),
	appHttp.NewProblemTypesRegistry,

	wire.Bind(new(http.Handler), new(*chi.Mux)),
	wire.Bind(new(appHttp.StrictServerInterface), new(*appHttp.ApiHandler)),
)
