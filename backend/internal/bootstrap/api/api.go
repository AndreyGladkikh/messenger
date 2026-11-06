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
	HttpServer  *appHttp.Server
}

var ApiSet = wire.NewSet(
	wire.Struct(new(Api), "*"),
	bootstrap.AppSet,
	appHttp.NewServer,
	appHttp.NewRouter,
	appHttp.NewApiHandler,
	appHttp.NewErrorHandler,
	appHttp.NewProblemTypesRegistry,
	websocket.NewHandler,

	wire.Bind(new(http.Handler), new(*chi.Mux)),
	wire.Bind(new(appHttp.StrictServerInterface), new(*appHttp.ApiHandler)),
)
