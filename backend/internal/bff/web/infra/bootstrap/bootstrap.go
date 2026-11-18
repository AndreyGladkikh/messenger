package bootstrap

import (
	"messenger/messenger/internal/bff/web/app/port"
	"messenger/messenger/internal/bff/web/app/query"
	"messenger/messenger/internal/bff/web/infra/integration/directory"
	"messenger/messenger/internal/bff/web/infra/integration/media"
	"messenger/messenger/internal/bff/web/infra/integration/messaging"
	"messenger/messenger/internal/platform/query_middleware"
	appQry "messenger/messenger/internal/shared/app/query"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	QueryHandlersSet,
	messaging.NewService,
	directory.NewService,
	media.NewService,

	wire.Bind(new(port.MessagingReader), new(*messaging.Service)),
	wire.Bind(new(port.DirectoryReader), new(*directory.Service)),
	wire.Bind(new(port.MediaReader), new(*media.Service)),
)

var QueryHandlersSet = wire.NewSet(
	query.NewGetHomePageHandler,
	ProvideGetHomePageHandler,
	query.NewSearchHandler,
	ProvideSearchHandler,
)

func ProvideGetHomePageHandler(
	h *query.GetHomePageHandler,
	m *query_middleware.Registry,
) appQry.Handler[query.GetHomePage, query.HomePage] {
	return query_middleware.Wrap(h, m)
}

func ProvideSearchHandler(
	h *query.SearchHandler,
	m *query_middleware.Registry,
) appQry.Handler[query.Search, query.SearchResponse] {
	return query_middleware.Wrap(h, m)
}
