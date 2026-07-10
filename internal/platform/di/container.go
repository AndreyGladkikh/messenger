package di

import (
	"context"
	"messenger/messenger/internal/adapters/out/postgres/queries"
	"messenger/messenger/internal/adapters/out/postgres/repositories"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/application/use_cases/send_message"
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/postgres"
)

type Container struct {
	TxManager *transaction.Manager

	SendMessageHandler *send_message.Handler

	MessageRepository message.Repository
}

// var once sync.Once
// var instance *Container

func InitContainer(ctx context.Context, cfg *config.Config) *Container {
	db, cleanup, err := postgres.NewPool(ctx)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	qs := queries.New(db)

	txManager := transaction.NewManager(db)

	messageRepository := repositories.NewMessageRepository(db, qs)

	return &Container{
		TxManager: txManager,

		SendMessageHandler: send_message.NewHandler(qs, txManager),

		MessageRepository: messageRepository,
	}

	// once.Do(func() {
	// 	instance = &Container{}
	// })
	// return instance
}
