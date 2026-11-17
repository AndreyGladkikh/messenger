package query_handler

import (
	webBffQry "messenger/messenger/internal/bff/web/app/query"
	directoryQry "messenger/messenger/internal/directory/app/query"
	identityQry "messenger/messenger/internal/identity/app/query"
	mediaQry "messenger/messenger/internal/media/app/query"
	conversationsQry "messenger/messenger/internal/messaging/app/query"
	appQuery "messenger/messenger/internal/shared/app/query"
	"uuid"
)

type Registry struct {
	GetHomePage appQuery.Handler[webBffQry.GetHomePage, webBffQry.HomePage]

	ListUserConversations appQuery.Handler[conversationsQry.ListUserConversations, []conversationsQry.ConversationSummary]

	GetDirectoryEntries  appQuery.Handler[directoryQry.GetEntries, []directoryQry.EntryView]
	GetDirectoryEntry    appQuery.Handler[directoryQry.GetEntry, directoryQry.EntryView]
	SearchDirectoryEntry appQuery.Handler[directoryQry.Search, []directoryQry.SearchResultItem]

	GetAccount appQuery.Handler[identityQry.GetAccount, identityQry.AccountView]

	GetDownloadMediaURL appQuery.Handler[mediaQry.GetDownloadURL, string]
	GetMediaPublicURLsHandler appQuery.Handler[mediaQry.GetPublicURLs, map[uuid.UUID]string]
}
