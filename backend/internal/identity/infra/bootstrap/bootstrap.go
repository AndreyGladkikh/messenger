package bootstrap

import (
	"messenger/messenger/internal/identity/app/command"
	passwordPort "messenger/messenger/internal/identity/app/password"
	"messenger/messenger/internal/identity/app/query"
	tokenPort "messenger/messenger/internal/identity/app/token"
	"messenger/messenger/internal/identity/domain/model/session"
	"messenger/messenger/internal/identity/domain/model/user"
	"messenger/messenger/internal/identity/infra/password"
	"messenger/messenger/internal/identity/infra/persistence"
	"messenger/messenger/internal/identity/infra/persistence/repository"
	"messenger/messenger/internal/identity/infra/token"
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
	password.NewHasher,
	token.NewService,

	wire.Bind(new(passwordPort.Hasher), new(*password.Hasher)),
	wire.Bind(new(tokenPort.Service), new(*token.Service)),
)

var PersistenceSet = wire.NewSet(
	persistence.NewQueryService,
	repository.NewUserRepository,
	repository.NewSessionRepository,
	wire.Bind(new(query.Reader), new(*persistence.QueryService)),
	wire.Bind(new(user.Repository), new(*repository.UserRepository)),
	wire.Bind(new(session.Repository), new(*repository.SessionRepository)),
)

var CommandHandlersSet = wire.NewSet(
	command.NewRegisterUserHandler,
	ProvideRegisterUserHandler,
	command.NewLoginUserHandler,
	ProvideLoginUserHandler,
	command.NewRefreshUserTokensHandler,
	ProvideRefreshUserTokensHandler,
	command.NewLogoutUserHandler,
	ProvideLogoutUserHandler,
)

func ProvideRegisterUserHandler(
	h *command.RegisterUserHandler,
	m *command_middleware.Registry,
) appCmd.Handler[command.RegisterUser, command.RegisterUserResponse] {
	return command_middleware.Wrap(h, m)
}

func ProvideLoginUserHandler(
	h *command.LoginUserHandler,
	m *command_middleware.Registry,
) appCmd.Handler[command.LoginUser, command.LoginUserResponse] {
	return command_middleware.Wrap(h, m)
}

func ProvideRefreshUserTokensHandler(
	h *command.RefreshUserTokensHandler,
	m *command_middleware.Registry,
) appCmd.Handler[command.RefreshUserTokens, command.RefreshUserTokensResponse] {
	return command_middleware.Wrap(h, m)
}

func ProvideLogoutUserHandler(
	h *command.LogoutUserHandler,
	m *command_middleware.Registry,
) appCmd.Handler[command.LogoutUser, any] {
	return command_middleware.Wrap(h, m)
}

var QueryHandlersSet = wire.NewSet(
	query.NewGetAccountHandler,
	ProvideGetAccountHandler,
)

func ProvideGetAccountHandler(
	h *query.GetAccountHandler,
	m *query_middleware.Registry,
) appQuery.Handler[query.GetAccount, query.AccountView] {
	return query_middleware.Wrap(h, m)
}
