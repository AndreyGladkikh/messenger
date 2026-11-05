package bootstrap

import (
	webBffPort "messenger/messenger/internal/bff/web/application/port"
	webBffQry "messenger/messenger/internal/bff/web/application/query"
	webBffDirectory "messenger/messenger/internal/bff/web/infra/integration/directory"
	webBffMessaging "messenger/messenger/internal/bff/web/infra/integration/messaging"
	cnvCmd "messenger/messenger/internal/conversations/app/command"
	cnvQry "messenger/messenger/internal/conversations/app/query"
	cnvService "messenger/messenger/internal/conversations/app/service"
	"messenger/messenger/internal/conversations/domain/model/conversation"
	"messenger/messenger/internal/conversations/domain/model/group_conversation"
	"messenger/messenger/internal/conversations/domain/model/membership"
	"messenger/messenger/internal/conversations/domain/model/message"
	cnvIdentity "messenger/messenger/internal/conversations/infra/integration/identity"
	cnvPersistence "messenger/messenger/internal/conversations/infra/persistence"
	cnvRepo "messenger/messenger/internal/conversations/infra/persistence/repository"
	drctrCmd "messenger/messenger/internal/directory/app/command"
	drctrQry "messenger/messenger/internal/directory/app/query"
	"messenger/messenger/internal/directory/domain/entry"
	drctrIdentity "messenger/messenger/internal/directory/infra/integration/identity"
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
	"messenger/messenger/internal/platform/rabbitmq"
	"messenger/messenger/internal/platform/rabbitmq/consumers"
	"messenger/messenger/internal/platform/rabbitmq/producers"
	"messenger/messenger/internal/platform/redis"
	baseRepo "messenger/messenger/internal/platform/repository"
	"messenger/messenger/internal/platform/sqlc"
	"messenger/messenger/internal/platform/store"
	"messenger/messenger/internal/platform/transactions"
	appEvnt "messenger/messenger/internal/shared/app/event"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var AppSet = wire.NewSet(
	logger.New,
	PersistenceSet,
	InitCommandBus,
	CommandHandlersSet,
	InitQueryBus,
	QueryHandlersSet,
	event.NewEventService,
	idPassImpl.NewHasher,
	idTokenImpl.NewService,
	pubsub.NewPubSub,
	pubsub.NewPersonalEventsPubSub,
	pubsub.NewChatEventsPubSub,
	cnvService.NewDirectConversationResolver,
	cnvService.NewMessageSender,
	cnvIdentity.NewParticipantLookup,
	drctrIdentity.NewUserLookup,
	drctrCmd.NewTargetLookup,
	rabbitmq.NewConnection,
	wire.Struct(new(event.Publisher), "*"),
	producers.NewConversationIntegrationEventsProducer,
	wire.Struct(new(producers.Producers), "*"),
	consumers.ConsumerRunners,
	consumers.NewRealtimePersonalDeliveryConsumer,
	wire.Struct(new(consumers.Consumers), "*"),
	webBffMessaging.NewService,
	webBffDirectory.NewService,

	wire.Bind(new(cnvService.ParticipantLookup), new(*cnvIdentity.ParticipantLookup)),
	wire.Bind(new(drctrCmd.UserLookup), new(*drctrIdentity.UserLookup)),
	wire.Bind(new(password.Hasher), new(*idPassImpl.Hasher)),
	wire.Bind(new(token.Service), new(*idTokenImpl.Service)),
	wire.Bind(new(appEvnt.Publisher), new(*event.Publisher)),
	wire.Bind(new(webBffPort.MessagingReader), new(*webBffMessaging.Service)),
	wire.Bind(new(webBffPort.DirectoryReader), new(*webBffDirectory.Service)),
)

var PersistenceSet = wire.NewSet(
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
	cnvRepo.NewGroupConversationRepository,
	cnvPersistence.NewQueryService,

	wire.Bind(new(sqlc.DBTX), new(*pgxpool.Pool)),
	wire.Bind(new(drctrQry.Searcher), new(*drctrPrst.QueryService)),
	wire.Bind(new(drctrQry.Reader), new(*drctrPrst.QueryService)),
	wire.Bind(new(idQry.Reader), new(*idPrst.QueryService)),
	wire.Bind(new(conversation.Repository), new(*cnvRepo.ConversationRepository)),
	wire.Bind(new(membership.Repository), new(*cnvRepo.MembershipRepository)),
	wire.Bind(new(message.Repository), new(*cnvRepo.MessageRepository)),
	wire.Bind(new(user.Repository), new(*idRepo.UserRepository)),
	wire.Bind(new(session.Repository), new(*idRepo.SessionRepository)),
	wire.Bind(new(entry.Repository), new(*drctrPrst.EntryRepository)),
	wire.Bind(new(group_conversation.Repository), new(*cnvRepo.GroupConversationRepository)),
	wire.Bind(new(cnvQry.Reader), new(*cnvPersistence.QueryService)),
)

var CommandHandlersSet = wire.NewSet(
	idCmd.NewRegisterUserHandler,
	idCmd.NewLoginUserHandler,
	idCmd.NewRefreshUserTokensHandler,
	cnvCmd.NewSendMessageHandler,
	wire.Struct(new(cnvCmd.StartDirectConversationHandler), "*"),
	drctrCmd.NewRegisterEntryHandler,
	cnvCmd.NewCreateGroupConversationHandler,
)

var QueryHandlersSet = wire.NewSet(
	idQry.NewGetAccountHandler,
	wire.Struct(new(idQry.GetUserIDByLoginHandler), "*"),
	drctrQry.NewSearchHandler,
	drctrQry.NewGetEntriesHandler,
	webBffQry.NewGetHomePageHandler,
	wire.Struct(new(drctrQry.GetEntryHandler), "*"),
	cnvQry.NewListUserConversationsHandler,
)

func InitCommandBus(
	logger *logger.Logger,
	txManager *transactions.Manager,
	eventService *event.EventService,
	registerUserHandler *idCmd.RegisterUserHandler,
	loginUserHandler *idCmd.LoginUserHandler,
	refreshSessionHandler *idCmd.RefreshUserTokensHandler,
	sendFirstDirectMsgHandler *cnvCmd.StartDirectConversationHandler,
	sendMsgHandler *cnvCmd.SendMessageHandler,
	registerEntryHandler *drctrCmd.RegisterEntryHandler,
) *commandbus.Bus {
	bus := commandbus.NewBus()

	loggingMiddlewareContainer := middlewares.NewLoggerMiddlewareContainer(logger)
	txMiddlewareContainer := middlewares.NewTransactionMiddlewareContainer(txManager, eventService)

	bus.Use(txMiddlewareContainer.Middleware)
	bus.Use(middlewares.Recoverer)
	bus.Use(loggingMiddlewareContainer.Middleware)

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
	getDirectoryEntryHandler *drctrQry.GetEntryHandler,
	getHomePageHandler *webBffQry.GetHomePageHandler,
	listUserConversationsHandler *cnvQry.ListUserConversationsHandler,
) *querybus.Bus {
	b := querybus.NewBus()

	b.RegisterHandler(getAccountHandler)
	b.RegisterHandler(searchDirectoryEntryHandler)
	b.RegisterHandler(getDirectoryEntryHandler)
	b.RegisterHandler(getHomePageHandler)
	b.RegisterHandler(listUserConversationsHandler)

	return b
}
