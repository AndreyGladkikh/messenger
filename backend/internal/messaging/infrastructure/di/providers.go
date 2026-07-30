package di

import (
	"messenger/messenger/internal/messaging/application/command/create_private_chat"
	"messenger/messenger/internal/messaging/application/command/send_message"
	"messenger/messenger/internal/messaging/application/event/message_sent"
	"messenger/messenger/internal/messaging/domain/chat"
	"messenger/messenger/internal/messaging/domain/chat_participant"
	"messenger/messenger/internal/messaging/domain/message"
	"messenger/messenger/internal/messaging/infrastructure/event"
	"messenger/messenger/internal/messaging/infrastructure/outbox_relay"
	"messenger/messenger/internal/messaging/infrastructure/postgres"
	"messenger/messenger/internal/messaging/infrastructure/repositories"
	"messenger/messenger/internal/messaging/infrastructure/sqlc"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/db"
	"messenger/messenger/internal/platform/db/transaction"
	"messenger/messenger/internal/platform/http_server"
	"messenger/messenger/internal/platform/logger"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ApiSet = wire.NewSet(
	NewApi,

	// bindings
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
	logger.New,
)

var OutboxRelaySet = wire.NewSet(
	outbox_relay.NewOutboxRelay,
	NewEventHandlerRegistry,
	CommonSet,
)

var CommonSet = wire.NewSet(
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
