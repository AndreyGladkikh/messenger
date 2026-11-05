package bootstrap

import (
	webBffQry "messenger/messenger/internal/bff/web/application/query"
	cnvQry "messenger/messenger/internal/conversations/app/query"
	drctrQry "messenger/messenger/internal/directory/app/query"
	idQry "messenger/messenger/internal/identity/app/query"
	"messenger/messenger/internal/platform/querymiddleware"
	appQuery "messenger/messenger/internal/shared/app/query"
)

func ProvideGetAccountQueryHandler(
	handler *idQry.GetAccountHandler,
	logger querymiddleware.Logger,
) appQuery.Handler[idQry.GetAccount, idQry.AccountView] {
	return querymiddleware.LoggingF(handler, logger)
}

func ProvideSearchDirectoryQueryHandler(
	handler *drctrQry.SearchHandler,
	logger querymiddleware.Logger,
) appQuery.Handler[drctrQry.Search, []drctrQry.SearchResultItem] {
	return querymiddleware.LoggingF(handler, logger)
}

func ProvideGetDirectoryEntryQueryHandler(
	handler *drctrQry.GetEntryHandler,
	logger querymiddleware.Logger,
) appQuery.Handler[drctrQry.GetEntry, drctrQry.EntryView] {
	return querymiddleware.LoggingF(handler, logger)
}

func ProvideGetDirectoryEntriesQueryHandler(
	handler *drctrQry.GetEntriesHandler,
	logger querymiddleware.Logger,
) appQuery.Handler[drctrQry.GetEntries, []drctrQry.EntryView] {
	return querymiddleware.LoggingF(handler, logger)
}

func ProvideListUserConversationsQueryHandler(
	handler *cnvQry.ListUserConversationsHandler,
	logger querymiddleware.Logger,
) appQuery.Handler[cnvQry.ListUserConversations, []cnvQry.ConversationSummary] {
	return querymiddleware.LoggingF(handler, logger)
}

func ProvideGetWebHomePageQueryHandler(
	handler *webBffQry.GetHomePageHandler,
	logger querymiddleware.Logger,
) appQuery.Handler[webBffQry.GetHomePage, webBffQry.HomePage] {
	return querymiddleware.LoggingF(handler, logger)
}
