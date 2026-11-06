package bootstrap

import (
	webBff "messenger/messenger/internal/bff/web/infra/bootstrap"
	conversations "messenger/messenger/internal/conversations/infra/bootstrap"
	directory "messenger/messenger/internal/directory/infra/bootstrap"
	identity "messenger/messenger/internal/identity/infra/bootstrap"
	"messenger/messenger/internal/platform/command_handler"
	"messenger/messenger/internal/platform/command_middleware"
	"messenger/messenger/internal/platform/event"
	"messenger/messenger/internal/platform/logger"
	"messenger/messenger/internal/platform/postgres"
	"messenger/messenger/internal/platform/pubsub"
	"messenger/messenger/internal/platform/query_handler"
	"messenger/messenger/internal/platform/query_middleware"
	"messenger/messenger/internal/platform/rabbitmq"
	"messenger/messenger/internal/platform/rabbitmq/consumers"
	"messenger/messenger/internal/platform/rabbitmq/producers"
	"messenger/messenger/internal/platform/redis"
	"messenger/messenger/internal/platform/repository"
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
	event.NewEventService,
	event.NewPublisher,
	pubsub.NewPubSub,
	pubsub.NewPersonalEventsPubSub,
	pubsub.NewChatEventsPubSub,
	rabbitmq.NewConnection,
	producers.NewConversationIntegrationEventsProducer,
	wire.Struct(new(producers.Producers), "*"),
	wire.Struct(new(consumers.Consumers), "*"),
	consumers.ConsumerRunners,
	consumers.NewRealtimePersonalDeliveryConsumer,
	query_middleware.NewLogging,
	query_middleware.NewValidation,
	command_middleware.NewLogging,
	command_middleware.NewRecovering,
	command_middleware.NewTransaction,
	wire.Struct(new(command_handler.Registry), "*"),
	wire.Struct(new(command_middleware.Registry), "*"),
	wire.Struct(new(query_handler.Registry), "*"),
	wire.Struct(new(query_middleware.Registry), "*"),

	webBff.Set,
	conversations.Set,
	directory.Set,
	identity.Set,

	wire.Bind(new(appEvnt.Publisher), new(*event.Publisher)),
)

var PersistenceSet = wire.NewSet(
	postgres.NewPool,
	transactions.NewManager,
	sqlc.New,
	store.NewStore,
	repository.NewRepository,
	redis.NewRedis,

	wire.Bind(new(sqlc.DBTX), new(*pgxpool.Pool)),
)
