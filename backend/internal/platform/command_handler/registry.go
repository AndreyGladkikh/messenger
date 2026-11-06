package command_handler

import (
	conversationsCmd "messenger/messenger/internal/conversations/app/command"
	directoryCmd "messenger/messenger/internal/directory/app/command"
	identityCmd "messenger/messenger/internal/identity/app/command"
	appCmd "messenger/messenger/internal/shared/app/command"
	"uuid"
)

type Registry struct {
	RegisterUser      appCmd.Handler[identityCmd.RegisterUser, identityCmd.RegisterUserResponse]
	LoginUser         appCmd.Handler[identityCmd.LoginUser, identityCmd.LoginUserResponse]
	RefreshUserTokens appCmd.Handler[identityCmd.RefreshUserTokens, identityCmd.RefreshUserTokensResponse]
	LogoutUser        appCmd.Handler[identityCmd.LogoutUser, any]

	StartDirectConversation appCmd.Handler[conversationsCmd.StartDirectConversation, conversationsCmd.StartDirectConversationResponse]
	SendMessage             appCmd.Handler[conversationsCmd.SendMessage, any]
	CreateGroupConversation appCmd.Handler[conversationsCmd.CreateGroupConversation, uuid.UUID]

	RegisterDirectoryEntry appCmd.Handler[directoryCmd.RegisterEntry, any]
}
