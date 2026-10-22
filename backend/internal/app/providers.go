package app

import (
	cnvCmd "messenger/messenger/internal/conversations/app/command"
	cnvService "messenger/messenger/internal/conversations/app/service"
	"messenger/messenger/internal/conversations/domain/model/conversation"
	"messenger/messenger/internal/conversations/domain/model/membership"
	"messenger/messenger/internal/conversations/domain/model/message"
	cnvIdentity "messenger/messenger/internal/conversations/infra/identity"
	cnvRepo "messenger/messenger/internal/conversations/infra/persistence/repository"
	drctrCmd "messenger/messenger/internal/directory/app/command"
	drctrQry "messenger/messenger/internal/directory/app/query"
	"messenger/messenger/internal/directory/domain/entry"
	drctrIdentity "messenger/messenger/internal/directory/infra/identity"
	drctrPrst "messenger/messenger/internal/directory/infra/persistence"
	idCmd "messenger/messenger/internal/identity/app/command"
	"messenger/messenger/internal/identity/app/password"
	idQry "messenger/messenger/internal/identity/app/query"
	"messenger/messenger/internal/identity/app/token"
	"messenger/messenger/internal/identity/domain/model/session"
	"messenger/messenger/internal/identity/domain/model/user"
	idPassImpl "messenger/messenger/internal/identity/infra/password"
	idPrst "messenger/messenger/internal/identity/infra/persistence"
	idRepo "messenger/messenger/internal/identity/infra/persistence/repository"
	idTokenImpl "messenger/messenger/internal/identity/infra/token"
	"messenger/messenger/internal/platform/commandbus"
	"messenger/messenger/internal/platform/commandbus/middlewares"
	"messenger/messenger/internal/platform/event"
	"messenger/messenger/internal/platform/logger"
	pg "messenger/messenger/internal/platform/postgres"
	"messenger/messenger/internal/platform/pubsub"
	"messenger/messenger/internal/platform/querybus"
	"messenger/messenger/internal/platform/redis"
	baseRepo "messenger/messenger/internal/platform/repository"
	"messenger/messenger/internal/platform/sqlc"
	"messenger/messenger/internal/platform/store"
	"messenger/messenger/internal/platform/transactions"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var AppSet = wire.NewSet(
	logger.New,
	Persistence,
	InitCommandBus,
	CommandHandlers,
	InitQueryBus,
	QueryHandlers,
	event.NewEventService,
	idPassImpl.NewHasher,
	idTokenImpl.NewService,
	pubsub.NewPubSub,
	pubsub.NewChatEventsPubSub,
	cnvService.NewDirectConversationResolver,
	// wire.Struct(new(cnvService.DirectConversationResolver), "*"),
	cnvService.NewMessageSender,
	cnvIdentity.NewParticipantLookup,
	drctrIdentity.NewUserLookup,
	drctrCmd.NewTargetLookup,

	wire.Bind(new(cnvService.ParticipantLookup), new(*cnvIdentity.ParticipantLookup)),
	wire.Bind(new(drctrCmd.UserLookup), new(*drctrIdentity.UserLookup)),
	wire.Bind(new(password.Hasher), new(*idPassImpl.Hasher)),
	wire.Bind(new(token.Service), new(*idTokenImpl.Service)),
)

var Persistence = wire.NewSet(
	pg.NewPool,
	transactions.NewManager,
	sqlc.New,
	store.NewStore,
	redis.NewRedis,
	drctrPrst.NewQueryService,
	idPrst.NewQueryService,
	baseRepo.NewRepository,
	cnvRepo.NewConversationRepository,
	cnvRepo.NewMembershipRepository,
	cnvRepo.NewMessageRepository,
	idRepo.NewUserRepository,
	idRepo.NewSessionRepository,
	drctrPrst.NewEntryRepository,

	wire.Bind(new(sqlc.DBTX), new(*pgxpool.Pool)),
	wire.Bind(new(drctrQry.Searcher), new(*drctrPrst.QueryService)),
	wire.Bind(new(idQry.Reader), new(*idPrst.QueryService)),
	wire.Bind(new(conversation.Repository), new(*cnvRepo.ConversationRepository)),
	wire.Bind(new(membership.Repository), new(*cnvRepo.MembershipRepository)),
	wire.Bind(new(message.Repository), new(*cnvRepo.MessageRepository)),
	wire.Bind(new(user.Repository), new(*idRepo.UserRepository)),
	wire.Bind(new(session.Repository), new(*idRepo.SessionRepository)),
	wire.Bind(new(entry.Repository), new(*drctrPrst.EntryRepository)),
)

var CommandHandlers = wire.NewSet(
	idCmd.NewRegisterUserHandler,
	idCmd.NewLoginUserHandler,
	idCmd.NewRefreshUserTokensHandler,
	cnvCmd.NewSendMessageHandler,
	cnvCmd.NewSendFirstDirectMessageHandler,
	drctrCmd.NewRegisterEntryHandler,
)

var QueryHandlers = wire.NewSet(
	idQry.NewGetAccountHandler,
	drctrQry.NewSearchHandler,
)

func InitCommandBus(
	logger *logger.Logger,
	txManager *transactions.Manager,
	eventService *event.EventService,
	registerUserHandler *idCmd.RegisterUserHandler,
	loginUserHandler *idCmd.LoginUserHandler,
	refreshSessionHandler *idCmd.RefreshUserTokensHandler,
	sendFirstDirectMsgHandler *cnvCmd.SendFirstDirectMessageHandler,
	sendMsgHandler *cnvCmd.SendMessageHandler,
	registerEntryHandler *drctrCmd.RegisterEntryHandler,
) *commandbus.Bus {
	bus := commandbus.NewBus()

	loggingMiddlewareContainer := middlewares.NewLoggerMiddlewareContainer(logger)
	txMiddlewareContainer := middlewares.NewTransactionMiddlewareContainer(txManager, eventService)

	bus.Use(txMiddlewareContainer.Middleware)
	bus.Use(middlewares.Recoverer)
	bus.Use(loggingMiddlewareContainer.Middleware)
	// bus.Use(middlewares.ErrorMapping)

	bus.RegisterHandler(registerUserHandler)
	bus.RegisterHandler(loginUserHandler)
	bus.RegisterHandler(refreshSessionHandler)
	bus.RegisterHandler(sendFirstDirectMsgHandler)
	bus.RegisterHandler(sendMsgHandler)
	bus.RegisterHandler(registerEntryHandler)

	return bus
}

func InitQueryBus(
	getAccountHandler *idQry.GetAccountHandler,
	searchDirectoryEntryHandler *drctrQry.SearchHandler,
) *querybus.Bus {
	b := querybus.NewBus()

	b.RegisterHandler(getAccountHandler)
	b.RegisterHandler(searchDirectoryEntryHandler)

	return b
}
