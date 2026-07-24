package di

import (
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/application/command/create_private_chat"
	"messenger/messenger/internal/application/command/send_message"
	"messenger/messenger/internal/application/event/message_sent"
	"messenger/messenger/internal/application/id"
	"messenger/messenger/internal/domain/chat"
	"messenger/messenger/internal/domain/chat_participant"
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/infrastructure/config"
	"messenger/messenger/internal/infrastructure/db"
	"messenger/messenger/internal/infrastructure/event"
	"messenger/messenger/internal/infrastructure/idprovider"
	"messenger/messenger/internal/infrastructure/logger"
	"messenger/messenger/internal/infrastructure/postgres"
	"messenger/messenger/internal/infrastructure/repositories"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var CommonSet = wire.NewSet(
	wire.Bind(new(id.Provider), new(*idprovider.Provider)),
	wire.Bind(new(db.DBTX), new(*pgxpool.Pool)),

	config.Load,

	db.New,
	db.NewStorage,
	Repositories,

	CommandHandlers,
	EventHandlers,

	// persistence
	postgres.NewPool,
	transaction.NewManager,

	// storages
	event.NewEventStorage,

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
