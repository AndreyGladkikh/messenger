//+build wireinject

package di

import (
	"database/sql"
	loggerAdapter "messenger/messenger/internal/adapters/out/logger"
	"messenger/messenger/internal/adapters/out/postgres/queries"
	"messenger/messenger/internal/adapters/out/postgres/repositories"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/adapters/out/uuid"
	"messenger/messenger/internal/application/command/send_message"
	"messenger/messenger/internal/application/id"
	appLogger "messenger/messenger/internal/application/logger"
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/event"
	"messenger/messenger/internal/platform/http_server"
	"messenger/messenger/internal/platform/postgres"

	"github.com/google/wire"
)

func InitializeApi() (*Container, func(), error) {
	wire.Build(
		NewContainer,

		// bindings
		wire.Bind(new(id.Provider), new(*uuid.Provider)),

		wire.Bind(new(queries.DBTX), new(*sql.DB)),

		wire.Bind(new(message.Repository), new(*repositories.MessageRepository)),
		
		wire.Bind(new(appLogger.Logger), new(*loggerAdapter.Logger)),

		config.Load,

		http_server.NewServer,
		http_server.NewController,

		// persistence
		postgres.NewPool,
		transaction.NewManager,
		queries.New,

		// buses
		BuildCommandBusForApi,
		event.NewBus,

		// repositories
		repositories.NewMessageRepository,
		// repositories.NewChatRepository,

		// storages
		event.NewEventStorage,

		// command handlers
		send_message.NewHandler,

		// other
		uuid.NewProvider,
		loggerAdapter.New,
	)
	return new(Container), func() {}, nil
}
