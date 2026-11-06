package bootstrap

import (
	"messenger/messenger/internal/messaging/app/command"
	"messenger/messenger/internal/messaging/app/query"
	"messenger/messenger/internal/messaging/app/service"
	"messenger/messenger/internal/messaging/domain/model/conversation"
	"messenger/messenger/internal/messaging/domain/model/group_conversation"
	"messenger/messenger/internal/messaging/domain/model/membership"
	"messenger/messenger/internal/messaging/domain/model/message"
	"messenger/messenger/internal/messaging/infra/integration"
	"messenger/messenger/internal/messaging/infra/persistence"
	"messenger/messenger/internal/messaging/infra/persistence/repository"
	"messenger/messenger/internal/platform/command_middleware"
	"messenger/messenger/internal/platform/query_middleware"
	appCmd "messenger/messenger/internal/shared/app/command"
	appQuery "messenger/messenger/internal/shared/app/query"
	"uuid"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	PersistenceSet,
	CommandHandlersSet,
	QueryHandlersSet,
	service.NewDirectConversationResolver,
	service.NewMessageSender,
	integration.NewParticipantLookup,

	wire.Bind(new(service.ParticipantLookup), new(*integration.ParticipantLookup)),
)

var PersistenceSet = wire.NewSet(
	repository.NewConversationRepository,
	repository.NewMembershipRepository,
	repository.NewMessageRepository,
	repository.NewGroupConversationRepository,
	persistence.NewQueryService,

	wire.Bind(new(conversation.Repository), new(*repository.ConversationRepository)),
	wire.Bind(new(membership.Repository), new(*repository.MembershipRepository)),
	wire.Bind(new(message.Repository), new(*repository.MessageRepository)),
	wire.Bind(new(group_conversation.Repository), new(*repository.GroupConversationRepository)),
	wire.Bind(new(query.Reader), new(*persistence.QueryService)),
)

var CommandHandlersSet = wire.NewSet(
	command.NewSendMessageHandler,
	ProvideSendMessageHandler,
	command.NewStartDirectConversationHandler,
	ProvideStartDirectConversationHandler,
	command.NewCreateGroupConversationHandler,
	ProvideCreateGroupConversationHandler,
)

func ProvideSendMessageHandler(
	h *command.SendMessageHandler,
	m *command_middleware.Registry,
) appCmd.Handler[command.SendMessage, any] {
	return command_middleware.Wrap(h, m)
}

func ProvideStartDirectConversationHandler(
	h *command.StartDirectConversationHandler,
	m *command_middleware.Registry,
) appCmd.Handler[command.StartDirectConversation, command.StartDirectConversationResponse] {
	return command_middleware.Wrap(h, m)
}

func ProvideCreateGroupConversationHandler(
	h *command.CreateGroupConversationHandler,
	m *command_middleware.Registry,
) appCmd.Handler[command.CreateGroupConversation, uuid.UUID] {
	return command_middleware.Wrap(h, m)
}

var QueryHandlersSet = wire.NewSet(
	query.NewListUserConversationsHandler,
	ProvideGetHomePageHandler,
)

func ProvideGetHomePageHandler(
	h *query.ListUserConversationsHandler,
	m *query_middleware.Registry,
) appQuery.Handler[query.ListUserConversations, []query.ConversationSummary] {
	return query_middleware.Wrap(h, m)
}
