package di

import (
	"messenger/messenger/internal/adapters/out/logger"
	"messenger/messenger/internal/adapters/out/postgres/queries"
	"messenger/messenger/internal/adapters/out/postgres/repositories"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/adapters/out/uuid"
	"messenger/messenger/internal/application/command/create_private_chat"
	"messenger/messenger/internal/application/command/send_message"
	"messenger/messenger/internal/application/event/message_sent"
	"messenger/messenger/internal/application/id"
	"messenger/messenger/internal/domain/chat"
	"messenger/messenger/internal/domain/chat_participant"
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/event"
	"messenger/messenger/internal/platform/postgres"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var CommonSet = wire.NewSet(
	wire.Bind(new(id.Provider), new(*uuid.Provider)),
	wire.Bind(new(queries.DBTX), new(*pgxpool.Pool)),
	
	config.Load,
	Repositories,
	CommandHandlers,
	EventHandlers,


	// persistence
	postgres.NewPool,
	transaction.NewManager,
	queries.New,

	// storages
	event.NewEventStorage,

	// other
	uuid.NewProvider,
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

// var Repositories = wire.NewSet(
// 	wire.Bind(new(message.Repository), new(*repositories.MessageRepository)),
// 	wire.Bind(new(chat.Repository), new(*repositories.ChatRepository)),
// 	wire.Bind(new(chat_participant.Repository), new(*repositories.ChatParticipantRepository)),

// 	repositories.NewMessageRepository,
// 	repositories.NewChatRepository,
// 	repositories.NewChatParticipantRepository,
// )

var CommandHandlers = wire.NewSet(
	create_private_chat.NewHandler,
	send_message.NewHandler,
)

var EventHandlers = wire.NewSet(
	message_sent.NewNotifyChatParticipantsHandler,
)