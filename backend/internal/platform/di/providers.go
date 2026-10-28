package di

import (
	cnvCommand "messenger/messenger/internal/conversations/app/command"
	cnvService "messenger/messenger/internal/conversations/app/service"
	"messenger/messenger/internal/conversations/domain/model/conversation"
	"messenger/messenger/internal/conversations/domain/model/membership"
	"messenger/messenger/internal/conversations/domain/model/message"
	cnvIdentity "messenger/messenger/internal/conversations/infra/identity"
	cnvRepo "messenger/messenger/internal/conversations/infra/persistence/repository"
	drctrCommand "messenger/messenger/internal/directory/app/command"
	drctrQuery "messenger/messenger/internal/directory/app/query"
	"messenger/messenger/internal/directory/domain/entry"
	drctrIdentity "messenger/messenger/internal/directory/infra/identity"
	drctrPrstnc "messenger/messenger/internal/directory/infra/persistence"
	idCommand "messenger/messenger/internal/identity/app/command"
	authPasswordPort "messenger/messenger/internal/identity/app/password"
	idQuery "messenger/messenger/internal/identity/app/query"
	authTokenPort "messenger/messenger/internal/identity/app/token"
	"messenger/messenger/internal/identity/domain/model/session"
	"messenger/messenger/internal/identity/domain/model/user"
	authPasswordAdapter "messenger/messenger/internal/identity/infra/password"
	idReader "messenger/messenger/internal/identity/infra/persistence"
	idRepository "messenger/messenger/internal/identity/infra/persistence/repository"
	authTokenAdapter "messenger/messenger/internal/identity/infra/token"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/event"
	"messenger/messenger/internal/platform/logger"
	"messenger/messenger/internal/platform/outbox_relay"
	pg "messenger/messenger/internal/platform/postgres"
	"messenger/messenger/internal/platform/pubsub"
	"messenger/messenger/internal/platform/querybus"
	"messenger/messenger/internal/platform/redis"
	baseRepository "messenger/messenger/internal/platform/repository"
	"messenger/messenger/internal/platform/sqlc"
	"messenger/messenger/internal/platform/store"
	"messenger/messenger/internal/platform/transactions"
	"messenger/messenger/internal/platform/transport/http"
	"messenger/messenger/internal/platform/transport/websocket"
	"messenger/messenger/internal/shared/app/chateventspublisher"
	netHttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ApiSet = wire.NewSet(
	NewApi,

	CommonSet,
	http.NewServer,
	wire.Struct(new(websocket.Handler), "*"),
	http.NewErrorHandler,
	http.NewRouter,
	http.NewApiHandler,
	BuildCommandBus,
	querybus.InitBus,
	authPasswordAdapter.NewHasher,
	authTokenAdapter.NewService,
	http.NewProblemTypesRegistry,

	wire.Bind(new(authPasswordPort.Hasher), new(*authPasswordAdapter.Hasher)),
	wire.Bind(new(authTokenPort.Service), new(*authTokenAdapter.Service)),
	wire.Bind(new(http.StrictServerInterface), new(*http.ApiHandler)),
	wire.Bind(new(netHttp.Handler), new(*chi.Mux)),
)

var OutboxRelaySet = wire.NewSet(
	outbox_relay.NewOutboxRelay,
	NewEventHandlerRegistry,
	CommonSet,
)

var CommonSet = wire.NewSet(
	config.Load,

	logger.New,

	// pg.NewPool,
	transactions.NewManager,
	sqlc.New,
	store.NewStore,
	// Repositories,
	Persistence,

	redis.NewRedis,

	CommandHandlers,
	QueryHandlers,
	EventHandlers,

	event.NewEventService,
	pubsub.NewPubSub,
	pubsub.NewChatEventsPubSub,
	cnvService.NewDirectConversationResolver,
	cnvService.NewMessageSender,
	cnvIdentity.NewParticipantLookup,
	idReader.NewQueryService,
	drctrIdentity.NewUserLookup,
	drctrCommand.NewTargetLookup,
	// dirQueryService.NewQueryService,

	// wire.Bind(new(sqlc.DBTX), new(*pgxpool.Pool)),
	wire.Bind(new(chateventspublisher.ChatEventsPublisher), new(*pubsub.ChatEventsPubSub)),
	wire.Bind(new(cnvService.ParticipantLookup), new(*cnvIdentity.ParticipantLookup)),
	wire.Bind(new(idQuery.Reader), new(*idReader.QueryService)),
	wire.Bind(new(drctrCommand.UserLookup), new(*drctrIdentity.UserLookup)),
	// wire.Bind(new(dirQuery.Searcher), new(*dirQueryService.QueryService)),
)

var Persistence = wire.NewSet(
	pg.NewPool,
	Repositories,
	drctrPrstnc.NewQueryService,

	wire.Bind(new(sqlc.DBTX), new(*pgxpool.Pool)),
	wire.Bind(new(drctrQuery.Searcher), new(*drctrPrstnc.QueryService)),
)

var Repositories = wire.NewSet(
	baseRepository.NewRepository,
	cnvRepo.NewConversationRepository,
	cnvRepo.NewMembershipRepository,
	cnvRepo.NewMessageRepository,
	idRepository.NewUserRepository,
	idRepository.NewSessionRepository,
	drctrPrstnc.NewEntryRepository,

	wire.Bind(new(conversation.Repository), new(*cnvRepo.ConversationRepository)),
	wire.Bind(new(membership.Repository), new(*cnvRepo.MembershipRepository)),
	wire.Bind(new(message.Repository), new(*cnvRepo.MessageRepository)),
	wire.Bind(new(user.Repository), new(*idRepository.UserRepository)),
	wire.Bind(new(session.Repository), new(*idRepository.SessionRepository)),
	wire.Bind(new(entry.Repository), new(*drctrPrstnc.EntryRepository)),
)

var CommandHandlers = wire.NewSet(
	idCommand.NewRegisterUserHandler,
	idCommand.NewLoginUserHandler,
	idCommand.NewRefreshUserTokensHandler,
	cnvCommand.NewSendMessageHandler,
	cnvCommand.NewSendFirstDirectMessageHandler,
	drctrCommand.NewRegisterEntryHandler,
)

var QueryHandlers = wire.NewSet(
	idQuery.NewGetAccountHandler,
	wire.Struct(new(idQuery.GetUserIDByLoginHandler), "*"),
	drctrQuery.NewSearchHandler,
)

var EventHandlers = wire.NewSet()
