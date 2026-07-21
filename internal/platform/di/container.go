package di

import (
	"messenger/messenger/internal/application/logger"
	"messenger/messenger/internal/platform/http_server"
)

type Api struct {
	Logger     logger.Logger
	HttpServer *http_server.Server
}

func NewApi(
	httpServer *http_server.Server,
	logger logger.Logger,
) *Api {
	return &Api{
		HttpServer: httpServer,
		Logger:     logger,
	}
}

// type Container struct {
// 	DB        *pgxpool.Pool
// 	TxManager *transaction.Manager

// 	Logger logger.Logger

// 	HttpServer *http_server.Server

// 	CommandBus *command_bus.Bus
// 	EventBus   *event.Bus

// 	SendMessageHandler *send_message.Handler
// 	CreatePrivateChatHandler *create_private_chat.Handler

// 	MessageRepository message.Repository
// }

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
