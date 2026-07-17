package di

import (
	"database/sql"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/application/command/send_message"
	"messenger/messenger/internal/application/logger"
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/platform/command_bus"
	"messenger/messenger/internal/platform/event"
	"messenger/messenger/internal/platform/http_server"
)

type Container struct {
	DB        *sql.DB
	TxManager *transaction.Manager

	Logger logger.Logger

	HttpServer *http_server.Server

	CommandBus *command_bus.Bus
	EventBus   *event.Bus

	SendMessageHandler *send_message.Handler

	MessageRepository message.Repository
}

func NewContainer(
	db *sql.DB,
	txManager *transaction.Manager,
	commandBus *command_bus.Bus,
	eventBus *event.Bus,
	sendMessageHandler *send_message.Handler,
	messageRepository message.Repository,
	httpServer *http_server.Server,
	logger logger.Logger,
) *Container {
	return &Container{
		DB:                 db,
		TxManager:          txManager,
		CommandBus:         commandBus,
		EventBus:           eventBus,
		SendMessageHandler: sendMessageHandler,
		MessageRepository:  messageRepository,
		HttpServer: httpServer,
		Logger: logger,
	}
}

// func InitContainer(ctx context.Context, cfg *config.Config) *Container {
// 	db, cleanup, err := postgres.NewPool(ctx, cfg)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer cleanup()

// 	qs := queries.New(db)

// 	txManager := transaction.NewManager(db)

// 	idProvider := new(uuid.Provider)

// 	messageRepository := repositories.NewMessageRepository(db, qs)

// 	sendMessageHandler := send_message.NewHandler(messageRepository, idProvider)

// 	commandBus := command_bus.BuildCommandBus(txManager, sendMessageHandler)
// 	eventBus := event.NewBus()

// 	return &Container{
// 		DB:        db,
// 		TxManager: txManager,

// 		CommandBus: commandBus,
// 		EventBus:   eventBus,

// 		SendMessageHandler: sendMessageHandler,

// 		MessageRepository: messageRepository,
// 	}
// }
