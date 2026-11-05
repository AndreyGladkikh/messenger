package bootstrap

import (
	"context"
	"fmt"
	"messenger/messenger/internal/bff/web/application/port"
	webBffQry "messenger/messenger/internal/bff/web/application/query"
	webBffMessaging "messenger/messenger/internal/bff/web/infra/integration/messaging"
	cnvQry "messenger/messenger/internal/conversations/app/query"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type emptyConversationsReader struct{}

func (emptyConversationsReader) ListConversationsForUser(
	context.Context,
	uuid.UUID,
) ([]cnvQry.ConversationSummary, error) {
	return nil, nil
}

type unexpectedDirectoryReader struct {
	t *testing.T
}

func (r unexpectedDirectoryReader) GetEntries(
	context.Context,
	[]port.EntryRef,
) ([]port.EntryView, error) {
	r.t.Fatal("directory reader must not be called without entry references")
	return nil, nil
}

type queryRecordingLogger struct {
	queries []string
}

func (l *queryRecordingLogger) InfoContext(_ context.Context, _ string, args ...any) {
	l.record(args...)
}

func (l *queryRecordingLogger) ErrorContext(_ context.Context, _ string, args ...any) {
	l.record(args...)
}

func (l *queryRecordingLogger) record(args ...any) {
	for i := 0; i+1 < len(args); i += 2 {
		if args[i] == "query" {
			l.queries = append(l.queries, fmt.Sprint(args[i+1]))
			return
		}
	}
}

func TestDecoratedHomePageHandlerLogsNestedQueriesOnce(t *testing.T) {
	logger := &queryRecordingLogger{}
	listRaw := cnvQry.NewListUserConversationsHandler(emptyConversationsReader{})
	listDecorated := ProvideListUserConversationsQueryHandler(listRaw, logger)
	messagingReader := webBffMessaging.NewService(listDecorated)
	homeRaw := webBffQry.NewGetHomePageHandler(messagingReader, unexpectedDirectoryReader{t: t})
	homeDecorated := ProvideGetWebHomePageQueryHandler(homeRaw, logger)

	homePage, err := homeDecorated.Handle(t.Context(), webBffQry.GetHomePage{UserID: uuid.New()})

	require.NoError(t, err)
	assert.Empty(t, homePage.Peers)
	assert.Equal(t, []string{
		cnvQry.ListUserConversationsQueryName,
		webBffQry.GetHomePageQueryName,
	}, logger.queries)
}
