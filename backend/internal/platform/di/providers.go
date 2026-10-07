package di

import (
	cnvCommand "messenger/messenger/internal/conversations/app/command"
	cnvService "messenger/messenger/internal/conversations/app/service"
	"messenger/messenger/internal/conversations/domain/model/conversation"
	"messenger/messenger/internal/conversations/domain/model/membership"
	"messenger/messenger/internal/conversations/domain/model/message"
	cnvLookup "messenger/messenger/internal/conversations/infra/lookup"
	cnvRepo "messenger/messenger/internal/conversations/infra/persistence/repository"
	idCommand "messenger/messenger/internal/identity/app/command"
	authPasswordPort "messenger/messenger/internal/identity/app/password"
	idQuery "messenger/messenger/internal/identity/app/query"
	authTokenPort "messenger/messenger/internal/identity/app/token"
	"messenger/messenger/internal/identity/domain/model/session"
	"messenger/messenger/internal/identity/domain/model/user"
	authPasswordAdapter "messenger/messenger/internal/identity/infra/password"
	idQueryReader "messenger/messenger/internal/identity/infra/persistence/reader"
	authRepository "messenger/messenger/internal/identity/infra/persistence/repository"
	authTokenAdapter "messenger/messenger/internal/identity/infra/token"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/event"
	"messenger/messenger/internal/platform/logger"
	"messenger/messenger/internal/platform/outbox_relay"
	"messenger/messenger/internal/platform/postgres"
	"messenger/messenger/internal/platform/pubsub"
	"messenger/messenger/internal/platform/querybus"
	"messenger/messenger/internal/platform/redis"
	"messenger/messenger/internal/platform/repository"
	"messenger/messenger/internal/platform/sqlc"
	"messenger/messenger/internal/platform/store"
	"messenger/messenger/internal/platform/transactions"
	"messenger/messenger/internal/platform/transport/http"
	"messenger/messenger/internal/shared/app/chateventspublisher"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ApiSet = wire.NewSet(
	NewApi,

	CommonSet,
	http.NewServer,
	http.NewWebsocketHandler,
	http.NewErrorHandler,
	http.NewRouter,
	http.NewApiHandler,
	BuildCommandBus,
	querybus.InitBus,
	authPasswordAdapter.NewHasher,
	authTokenAdapter.NewService,

	wire.Bind(new(authPasswordPort.Hasher), new(*authPasswordAdapter.Hasher)),
	wire.Bind(new(authTokenPort.Service), new(*authTokenAdapter.Service)),
	wire.Bind(new(http.StrictServerInterface), new(*http.ApiHandler)),
)

var OutboxRelaySet = wire.NewSet(
	outbox_relay.NewOutboxRelay,
	NewEventHandlerRegistry,
	CommonSet,
)

var CommonSet = wire.NewSet(
	config.Load,

	sqlc.New,
	store.NewStore,
	Repositories,

	CommandHandlers,
	QueryHandlers,
	EventHandlers,

	postgres.NewPool,
	transactions.NewManager,

	redis.NewRedis,

	event.NewEventService,

	logger.New,
	pubsub.NewPubSub,
	pubsub.NewChatEventsPubSub,
	cnvService.NewDirectConversationResolver,
	cnvService.NewMessageSender,
	cnvLookup.NewParticipantLookup,
	idQueryReader.NewReader,

	wire.Bind(new(sqlc.DBTX), new(*pgxpool.Pool)),
	wire.Bind(new(chateventspublisher.ChatEventsPublisher), new(*pubsub.ChatEventsPubSub)),
	wire.Bind(new(cnvService.ParticipantLookup), new(*cnvLookup.ParticipantLookup)),
	wire.Bind(new(idQuery.Reader), new(*idQueryReader.Reader)),
)

var Repositories = wire.NewSet(
	repository.NewRepository,
	cnvRepo.NewConversationRepository,
	cnvRepo.NewMembershipRepository,
	cnvRepo.NewMessageRepository,
	authRepository.NewUserRepository,
	authRepository.NewSessionRepository,

	wire.Bind(new(message.Repository), new(*cnvRepo.MessageRepository)),
	wire.Bind(new(conversation.Repository), new(*cnvRepo.ConversationRepository)),
	wire.Bind(new(membership.Repository), new(*cnvRepo.MembershipRepository)),
	wire.Bind(new(user.Repository), new(*authRepository.UserRepository)),
	wire.Bind(new(session.Repository), new(*authRepository.SessionRepository)),
)

var CommandHandlers = wire.NewSet(
	idCommand.NewRegisterUserHandler,
	idCommand.NewLoginUserHandler,
	idCommand.NewRefreshUserTokensHandler,
	cnvCommand.NewSendMessageHandler,
	cnvCommand.NewSendFirstDirectMessageHandler,
)

var QueryHandlers = wire.NewSet(
	idQuery.NewGetAccountHandler,
)

var EventHandlers = wire.NewSet()
