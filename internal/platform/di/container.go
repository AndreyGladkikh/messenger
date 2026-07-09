package di

import (
	"sync"

	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/application/use_cases/send_message"
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/platform/config"
)

type Container struct {
	TxManager *transaction.Manager

	SendMessageHandler *send_message.Handler

	MessageRepository message.Repository
}

var once sync.Once
var instance *Container

func InitContainer(cfg *config.Config) *Container {
	once.Do(func() {
		instance = &Container{}
	})
	return instance
}
