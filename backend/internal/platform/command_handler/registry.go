package command_handler

import (
	directoryCmd "messenger/messenger/internal/directory/app/command"
	identityCmd "messenger/messenger/internal/identity/app/command"
	mediaCmd "messenger/messenger/internal/media/app/command"
	conversationsCmd "messenger/messenger/internal/messaging/app/command"
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

	SetUserProfile appCmd.Handler[directoryCmd.SetUserProfile, any]

	InitMediaUpload appCmd.Handler[mediaCmd.InitUpload, mediaCmd.InitedUpload]
	CompleteMediaUpload appCmd.Handler[mediaCmd.ComlpeteUpload, uuid.UUID]
}
