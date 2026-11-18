package bootstrap

import (
	"messenger/messenger/internal/directory/app/command"
	"messenger/messenger/internal/directory/app/port"
	"messenger/messenger/internal/directory/app/query"
	"messenger/messenger/internal/directory/app/view"
	"messenger/messenger/internal/directory/domain/model/handle"
	"messenger/messenger/internal/directory/domain/model/user_profile"
	"messenger/messenger/internal/directory/infra/integration/identity"
	"messenger/messenger/internal/directory/infra/integration/media"
	"messenger/messenger/internal/directory/infra/persistence"
	"messenger/messenger/internal/directory/infra/persistence/repository"
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

	identity.NewService,
	media.NewService,

	wire.Bind(new(port.IdentityReader), new(*identity.Service)),
	wire.Bind(new(port.MediaReader), new(*media.Service)),
)

var PersistenceSet = wire.NewSet(
	persistence.NewQueryService,
	repository.NewUserProfile,
	repository.NewHandle,

	wire.Bind(new(user_profile.Repository), new(*repository.UserProfile)),
	wire.Bind(new(handle.Repository), new(*repository.Handle)),
	wire.Bind(new(query.Searcher), new(*persistence.QueryService)),
	wire.Bind(new(port.DirectoryReader), new(*persistence.QueryService)),
)

var CommandHandlersSet = wire.NewSet(
	command.NewSetUserProfileHandler,
	ProvideSetUserProfileHandler,
)

func ProvideSetUserProfileHandler(
	h *command.SetUserProfileHandler,
	m *command_middleware.Registry,
) appCmd.Handler[command.SetUserProfile, any] {
	return command_middleware.Wrap(h, m)
}

var QueryHandlersSet = wire.NewSet(
	query.NewSearchHandler,
	ProvideSearchHandler,
	query.NewGetUserProfilesHandler,
	ProvideGetUserProfilesHandler,
)

func ProvideSearchHandler(
	h *query.SearchHandler,
	m *query_middleware.Registry,
) appQuery.Handler[query.Search, []query.SearchResultItem] {
	return query_middleware.Wrap(h, m)
}

func ProvideGetUserProfilesHandler(
	h *query.GetUserProfilesHandler,
	m *query_middleware.Registry,
) appQuery.Handler[query.GetUserProfiles, []view.UserProfile] {
	return query_middleware.Wrap(h, m)
}
