package bootstrap

import (
	"messenger/messenger/internal/bff/web/application/query"
	"messenger/messenger/internal/platform/querymiddleware"
	appQry "messenger/messenger/internal/shared/app/query"
	"messenger/messenger/internal/bff/web/application/port"
	"messenger/messenger/internal/bff/web/infra/integration/directory"
	"messenger/messenger/internal/bff/web/infra/integration/messaging"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	QueryHandlersSet,
	messaging.NewService,
	directory.NewService,

	wire.Bind(new(port.MessagingReader), new(*messaging.Service)),
	wire.Bind(new(port.DirectoryReader), new(*directory.Service)),
)

var QueryHandlersSet = wire.NewSet(
	wire.Struct(new(QueryHandlers), "*"),
	query.NewGetHomePageHandler,
	ProvideGetHomePageHandler,
)

type QueryHandlers struct {
	GetHomePage appQry.Handler[query.GetHomePage, query.HomePage]
}

func ProvideGetHomePageHandler(
	h *query.GetHomePageHandler,
	m *querymiddleware.Registry,
) appQry.Handler[query.GetHomePage, query.HomePage] {
	return querymiddleware.Wrap(h, m)
}
