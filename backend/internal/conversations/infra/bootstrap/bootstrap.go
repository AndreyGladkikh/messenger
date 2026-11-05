package bootstrap

import (
	"messenger/messenger/internal/conversations/app/command"
	"messenger/messenger/internal/conversations/app/query"
	"messenger/messenger/internal/conversations/app/service"
	"messenger/messenger/internal/conversations/domain/model/conversation"
	"messenger/messenger/internal/conversations/domain/model/group_conversation"
	"messenger/messenger/internal/conversations/domain/model/membership"
	"messenger/messenger/internal/conversations/domain/model/message"
	"messenger/messenger/internal/conversations/infra/integration"
	"messenger/messenger/internal/conversations/infra/persistence"
	"messenger/messenger/internal/conversations/infra/persistence/repository"
	"messenger/messenger/internal/platform/command_middleware"
	"messenger/messenger/internal/platform/querymiddleware"
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
	wire.Struct(new(command.StartDirectConversationHandler), "*"),
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
	m *querymiddleware.Registry,
) appQuery.Handler[query.ListUserConversations, []query.ConversationSummary] {
	return querymiddleware.Wrap(h, m)
}
