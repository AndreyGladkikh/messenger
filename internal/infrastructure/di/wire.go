// +build wireinject

package di

import (
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/application/command/create_private_chat"
	"messenger/messenger/internal/application/command/send_message"
	"messenger/messenger/internal/application/id"
	"messenger/messenger/internal/domain/chat"
	"messenger/messenger/internal/domain/chat_participant"
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/infrastructure/config"
	"messenger/messenger/internal/infrastructure/db"
	"messenger/messenger/internal/infrastructure/event"
	"messenger/messenger/internal/infrastructure/http_server"
	"messenger/messenger/internal/infrastructure/idprovider"
	"messenger/messenger/internal/infrastructure/logger"
	"messenger/messenger/internal/infrastructure/postgres"
	"messenger/messenger/internal/infrastructure/repositories"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitializeApi() (*Api, func(), error) {
	wire.Build(
		NewApi,

		// bindings
		wire.Bind(new(id.Provider), new(*idprovider.Provider)),

		wire.Bind(new(db.DBTX), new(*pgxpool.Pool)),

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
		db.New,
		db.NewStorage,

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
		idprovider.NewProvider,
		logger.New,
	)

	return new(Api), func() {}, nil
}

func InitializeEventProcessor() (*event.Processor, func(), error) {
	wire.Build(
		event.NewProcessor,
		BuildEventBusForProcessor,
		CommonSet,
	)

	return new(event.Processor), func() {}, nil
}
