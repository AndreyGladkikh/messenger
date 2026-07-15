package di

import (
	"context"
	"database/sql"
	"messenger/messenger/internal/adapters/out/postgres/queries"
	"messenger/messenger/internal/adapters/out/postgres/repositories"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/adapters/out/uuid"
	"messenger/messenger/internal/application/command/send_message"
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/platform/command"
	"messenger/messenger/internal/platform/command_bus"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/event"
	"messenger/messenger/internal/platform/postgres"
)

type Container struct {
	DB        *sql.DB
	TxManager *transaction.Manager

	CommandBus *command.Bus
	EventBus   *event.Bus

	SendMessageHandler *send_message.Handler

	MessageRepository message.Repository
}

func NewContainer(
	db *sql.DB,
	txManager *transaction.Manager,
	commandBus *command.Bus,
	eventBus *event.Bus,
	sendMessageHandler *send_message.Handler,
	messageRepository message.Repository,
) *Container {
	return &Container{
		DB:                 db,
		TxManager:          txManager,
		CommandBus:         commandBus,
		EventBus:           eventBus,
		SendMessageHandler: sendMessageHandler,
		MessageRepository:  messageRepository,
	}
}

// var once sync.Once
// var instance *Container

func InitContainer(ctx context.Context, cfg *config.Config) *Container {
	db, cleanup, err := postgres.NewPool(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	qs := queries.New(db)

	txManager := transaction.NewManager(db)

	idProvider := new(uuid.Provider)

	messageRepository := repositories.NewMessageRepository(db, qs)

	sendMessageHandler := send_message.NewHandler(messageRepository, idProvider)

	commandBus := command_bus.BuildCommandBus(txManager, sendMessageHandler)
	eventBus := event.NewBus()

	return &Container{
		DB:        db,
		TxManager: txManager,

		CommandBus: commandBus,
		EventBus:   eventBus,

		SendMessageHandler: sendMessageHandler,

		MessageRepository: messageRepository,
	}

	// once.Do(func() {
	// 	instance = &Container{}
	// })
	// return instance
}
