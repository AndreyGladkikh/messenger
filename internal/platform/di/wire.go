//+build wireinject

package di

import (
	loggerAdapter "messenger/messenger/internal/adapters/out/logger"
	"messenger/messenger/internal/adapters/out/postgres/queries"
	"messenger/messenger/internal/adapters/out/postgres/repositories"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/adapters/out/uuid"
	"messenger/messenger/internal/application/command/create_private_chat"
	"messenger/messenger/internal/application/command/send_message"
	"messenger/messenger/internal/application/id"
	// appLogger "messenger/messenger/internal/application/logger"
	"messenger/messenger/internal/domain/chat"
	"messenger/messenger/internal/domain/chat_participant"
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/event"
	"messenger/messenger/internal/platform/http_server"
	"messenger/messenger/internal/platform/postgres"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitializeApi() (*Api, func(), error) {
	wire.Build(
		NewApi,

		// bindings
		wire.Bind(new(id.Provider), new(*uuid.Provider)),

		wire.Bind(new(queries.DBTX), new(*pgxpool.Pool)),

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
		queries.New,

		// buses
		BuildCommandBusForApi,
		// event.NewBus,

		// repositories
		repositories.NewMessageRepository,
		repositories.NewChatRepository,
		repositories.NewChatParticipantRepository,

		// storages
		event.NewEventStorage,

		// command handlers
		create_private_chat.NewHandler,
		send_message.NewHandler,

		// other
		uuid.NewProvider,
		loggerAdapter.New,
	)

	return new(Api), func() {}, nil
}

// func InitializeEventProcessor() (*event.Processor, func(), error) {
// 	wire.Build(
// 		event.NewProcessor,

// 		// bindings
// 		wire.Bind(new(id.Provider), new(*uuid.Provider)),

// 		wire.Bind(new(queries.DBTX), new(*pgxpool.Pool)),

// 		wire.Bind(new(message.Repository), new(*repositories.MessageRepository)),
// 		wire.Bind(new(chat.Repository), new(*repositories.ChatRepository)),
// 		wire.Bind(new(chat_participant.Repository), new(*repositories.ChatParticipantRepository)),

// 		// wire.Bind(new(appLogger.Logger), new(*loggerAdapter.Logger)),

// 		config.Load,

// 		http_server.NewServer,
// 		http_server.NewController,

// 		// persistence
// 		postgres.NewPool,
// 		transaction.NewManager,
// 		queries.New,

// 		// buses
// 		BuildCommandBusForApi,
// 		// event.NewBus,

// 		// repositories
// 		repositories.NewMessageRepository,
// 		repositories.NewChatRepository,
// 		repositories.NewChatParticipantRepository,

// 		// storages
// 		event.NewEventStorage,

// 		// command handlers
// 		create_private_chat.NewHandler,
// 		send_message.NewHandler,

// 		// other
// 		uuid.NewProvider,
// 		loggerAdapter.New,
// 	)

// 	return new(event.Processor), func() {}, nil
// }