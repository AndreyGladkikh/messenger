package query_handler

import (
	webBffQry "messenger/messenger/internal/bff/web/app/query"
	directoryQry "messenger/messenger/internal/directory/app/query"
	directoryView "messenger/messenger/internal/directory/app/view"
	identityQry "messenger/messenger/internal/identity/app/query"
	mediaQry "messenger/messenger/internal/media/app/query"
	conversationsQry "messenger/messenger/internal/messaging/app/query"
	appQry "messenger/messenger/internal/shared/app/query"
	"uuid"
)

type Registry struct {
	GetHomePage appQry.Handler[webBffQry.GetHomePage, webBffQry.HomePage]
	Search appQry.Handler[webBffQry.Search, webBffQry.SearchResponse]

	ListUserConversations appQry.Handler[conversationsQry.ListUserConversations, []conversationsQry.ConversationSummary]

	GetUserProfiles    appQry.Handler[directoryQry.GetUserProfiles, []directoryView.UserProfile]
	SearchDirectoryEntry appQry.Handler[directoryQry.Search, []directoryQry.SearchResultItem]

	GetAccount appQry.Handler[identityQry.GetAccount, identityQry.AccountView]

	GetDownloadMediaURL appQry.Handler[mediaQry.GetDownloadURL, string]
	GetMediaPublicURLsHandler appQry.Handler[mediaQry.GetPublicURLs, map[uuid.UUID]string]
	MediaExists appQry.Handler[mediaQry.MediaExists, bool]
}
