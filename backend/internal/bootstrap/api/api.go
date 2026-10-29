package api

import (
	"messenger/messenger/internal/bootstrap"
	appHttp "messenger/messenger/internal/platform/transport/http"
	"messenger/messenger/internal/platform/transport/websocket"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/wire"
)

type Api struct {
	HttpHandler http.Handler
	HttpServer *appHttp.Server
}

var ApiSet = wire.NewSet(
	wire.Struct(new(Api), "*"),
	appHttp.NewServer,
	bootstrap.AppSet,
	appHttp.NewRouter,
	appHttp.NewApiHandler,
	appHttp.NewErrorHandler,
	wire.Struct(new(websocket.Handler), "*"),
	appHttp.NewProblemTypesRegistry,

	wire.Bind(new(http.Handler), new(*chi.Mux)),
	wire.Bind(new(appHttp.StrictServerInterface), new(*appHttp.ApiHandler)),
)
