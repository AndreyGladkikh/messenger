package bootstrap

import (
	"messenger/messenger/internal/directory/app/command"
	"messenger/messenger/internal/directory/app/query"
	"messenger/messenger/internal/directory/domain/entry"
	"messenger/messenger/internal/directory/infra/integration/identity"
	"messenger/messenger/internal/directory/infra/persistence"
	"messenger/messenger/internal/platform/command_middleware"
	"messenger/messenger/internal/platform/query_middleware"
	appCmd "messenger/messenger/internal/shared/app/command"
	appQuery "messenger/messenger/internal/shared/app/query"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	PersistenceSet,
	CommandHandlersSet,
	QueryHandlersSet,

	identity.NewUserLookup,
	command.NewTargetLookup,

	wire.Bind(new(command.UserLookup), new(*identity.UserLookup)),
)

var PersistenceSet = wire.NewSet(
	persistence.NewQueryService,
	persistence.NewEntryRepository,

	wire.Bind(new(entry.Repository), new(*persistence.EntryRepository)),
	wire.Bind(new(query.Searcher), new(*persistence.QueryService)),
	wire.Bind(new(query.Reader), new(*persistence.QueryService)),
)

var CommandHandlersSet = wire.NewSet(
	command.NewRegisterEntryHandler,
	ProvideRegisterEntryHandler,
)

func ProvideRegisterEntryHandler(
	h *command.RegisterEntryHandler,
	m *command_middleware.Registry,
) appCmd.Handler[command.RegisterEntry, any] {
	return command_middleware.Wrap(h, m)
}

var QueryHandlersSet = wire.NewSet(
	query.NewSearchHandler,
	ProvideSearchHandler,
	query.NewGetEntryHandler,
	ProvideGetEntryHandler,
	query.NewGetEntriesHandler,
	ProvideGetEntriesHandler,
)

func ProvideSearchHandler(
	h *query.SearchHandler,
	m *query_middleware.Registry,
) appQuery.Handler[query.Search, []query.SearchResultItem] {
	return query_middleware.Wrap(h, m)
}

func ProvideGetEntryHandler(
	h *query.GetEntryHandler,
	m *query_middleware.Registry,
) appQuery.Handler[query.GetEntry, query.EntryView] {
	return query_middleware.Wrap(h, m)
}

func ProvideGetEntriesHandler(
	h *query.GetEntriesHandler,
	m *query_middleware.Registry,
) appQuery.Handler[query.GetEntries, []query.EntryView] {
	return query_middleware.Wrap(h, m)
}
