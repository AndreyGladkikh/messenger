package di

import (
	"messenger/messenger/internal/auth/application/command/login"
	"messenger/messenger/internal/auth/application/command/refresh"
	"messenger/messenger/internal/auth/application/command/register"
	authPasswordPort "messenger/messenger/internal/auth/application/password"
	"messenger/messenger/internal/auth/application/query/get_current_user"
	authTokenPort "messenger/messenger/internal/auth/application/token"
	"messenger/messenger/internal/auth/domain/session"
	"messenger/messenger/internal/auth/domain/user"
	authPasswordAdapter "messenger/messenger/internal/auth/infrastructure/password"
	authRepository "messenger/messenger/internal/auth/infrastructure/repository"
	authTokenAdapter "messenger/messenger/internal/auth/infrastructure/token"
	"messenger/messenger/internal/messaging/application/command/create_private_chat"
	"messenger/messenger/internal/messaging/application/command/send_message"
	"messenger/messenger/internal/messaging/application/event/message_sent"
	"messenger/messenger/internal/messaging/application/query/get_chat_list"
	"messenger/messenger/internal/messaging/domain/chat"
	"messenger/messenger/internal/messaging/domain/chat_participant"
	"messenger/messenger/internal/messaging/domain/message"
	"messenger/messenger/internal/messaging/infrastructure/event"
	"messenger/messenger/internal/messaging/infrastructure/outbox_relay"
	messagingRepository "messenger/messenger/internal/messaging/infrastructure/repository"
	"messenger/messenger/internal/messaging/infrastructure/sqlc"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/db"
	"messenger/messenger/internal/platform/db/transaction"
	"messenger/messenger/internal/platform/http_server"
	"messenger/messenger/internal/platform/logger"
	"messenger/messenger/internal/platform/postgres"
	"messenger/messenger/internal/platform/querybus"
	"messenger/messenger/internal/shared/infrastructure/repository"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ApiSet = wire.NewSet(
	NewApi,

	CommonSet,
	http_server.NewServer,
	http_server.NewController,
	BuildCommandBus,
	querybus.InitBus,
	authPasswordAdapter.NewHasher,
	authTokenAdapter.NewService,

	wire.Bind(new(authPasswordPort.Hasher), new(*authPasswordAdapter.Hasher)),
	wire.Bind(new(authTokenPort.Service), new(*authTokenAdapter.Service)),
)

var OutboxRelaySet = wire.NewSet(
	outbox_relay.NewOutboxRelay,
	NewEventHandlerRegistry,
	CommonSet,
)

var CommonSet = wire.NewSet(
	config.Load,

	sqlc.New,
	db.NewStorage,
	Repositories,

	CommandHandlers,
	QueryHandlers,
	EventHandlers,

	// persistence
	postgres.NewPool,
	transaction.NewManager,

	event.NewEventService,

	// other
	logger.New,

	wire.Bind(new(sqlc.DBTX), new(*pgxpool.Pool)),
)

var Repositories = wire.NewSet(
	repository.NewRepository,
	messagingRepository.NewMessageRepository,
	messagingRepository.NewChatRepository,
	messagingRepository.NewChatParticipantRepository,
	authRepository.NewUserRepository,
	authRepository.NewSessionRepository,

	wire.Bind(new(message.Repository), new(*messagingRepository.MessageRepository)),
	wire.Bind(new(chat.Repository), new(*messagingRepository.ChatRepository)),
	wire.Bind(new(chat_participant.Repository), new(*messagingRepository.ChatParticipantRepository)),
	wire.Bind(new(user.Repository), new(*authRepository.UserRepository)),
	wire.Bind(new(session.Repository), new(*authRepository.SessionRepository)),
)

var CommandHandlers = wire.NewSet(
	register.NewHandler,
	login.NewHandler,
	refresh.NewHandler,
	create_private_chat.NewHandler,
	send_message.NewHandler,
)

var QueryHandlers = wire.NewSet(
	get_current_user.NewHandler,
	get_chat_list.NewHandler,
)

var EventHandlers = wire.NewSet(
	message_sent.NewNotifyChatParticipantsHandler,
	message_sent.NewRebuildQueryModelHandler,
)
