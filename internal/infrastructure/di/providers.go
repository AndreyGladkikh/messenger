package di

import (
	"messenger/messenger/internal/application/command/create_private_chat"
	"messenger/messenger/internal/application/command/send_message"
	"messenger/messenger/internal/application/event/message_sent"
	"messenger/messenger/internal/application/id"
	"messenger/messenger/internal/domain/chat"
	"messenger/messenger/internal/domain/chat_participant"
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/infrastructure/config"
	"messenger/messenger/internal/infrastructure/db"
	"messenger/messenger/internal/infrastructure/db/transaction"
	"messenger/messenger/internal/infrastructure/event"
	"messenger/messenger/internal/infrastructure/http_server"
	"messenger/messenger/internal/infrastructure/idprovider"
	"messenger/messenger/internal/infrastructure/logger"
	"messenger/messenger/internal/infrastructure/outbox_relay"
	"messenger/messenger/internal/infrastructure/postgres"
	"messenger/messenger/internal/infrastructure/repositories"
	"messenger/messenger/internal/infrastructure/sqlc"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ApiSet = wire.NewSet(
	NewApi,

	// bindings
	wire.Bind(new(id.Provider), new(*idprovider.Provider)),

	wire.Bind(new(sqlc.DBTX), new(*pgxpool.Pool)),

	wire.Bind(new(message.Repository), new(*repositories.MessageRepository)),
	wire.Bind(new(chat.Repository), new(*repositories.ChatRepository)),
	wire.Bind(new(chat_participant.Repository), new(*repositories.ChatParticipantRepository)),

	// wire.Bind(new(appLogger.Logger), new(*loggerAdapter.Logger)),

	config.Load,

	http_server.NewServer,
	http_server.NewController,

	// persistence
	postgres.NewPool,
	transaction.NewManager,
	sqlc.New,
	db.NewStorage,

	// buses
	BuildCommandBusForApi,
	// event.NewBus,

	// repositories
	repositories.NewMessageRepository,
	repositories.NewChatRepository,
	repositories.NewChatParticipantRepository,

	event.NewEventService,

	// command handlers
	create_private_chat.NewHandler,
	send_message.NewHandler,

	// other
	idprovider.NewProvider,
	logger.New,
)

var OutboxRelaySet = wire.NewSet(
	outbox_relay.NewOutboxRelay,
	NewEventHandlerRegistry,
	CommonSet,
)

var CommonSet = wire.NewSet(
	wire.Bind(new(id.Provider), new(*idprovider.Provider)),
	wire.Bind(new(sqlc.DBTX), new(*pgxpool.Pool)),

	config.Load,

	sqlc.New,
	db.NewStorage,
	Repositories,

	CommandHandlers,
	EventHandlers,

	// persistence
	postgres.NewPool,
	transaction.NewManager,

	event.NewEventService,

	// other
	idprovider.NewProvider,
	logger.New,
)

var Repositories = wire.NewSet(
	wire.Bind(new(message.Repository), new(*repositories.MessageRepository)),
	wire.Bind(new(chat.Repository), new(*repositories.ChatRepository)),
	wire.Bind(new(chat_participant.Repository), new(*repositories.ChatParticipantRepository)),

	repositories.NewMessageRepository,
	repositories.NewChatRepository,
	repositories.NewChatParticipantRepository,
)

var CommandHandlers = wire.NewSet(
	create_private_chat.NewHandler,
	send_message.NewHandler,
)

var EventHandlers = wire.NewSet(
	message_sent.NewNotifyChatParticipantsHandler,
	message_sent.NewRebuildQueryModelHandler,
)
