package query

import (
	"context"
	"messenger/messenger/internal/bff/web/application/port"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type messagingReaderStub struct {
	listUserConversations func(context.Context, uuid.UUID) ([]port.ConversationSummary, error)
}

func (s messagingReaderStub) ListUserConvesations(
	ctx context.Context,
	userID uuid.UUID,
) ([]port.ConversationSummary, error) {
	return s.listUserConversations(ctx, userID)
}

type directoryReaderStub struct {
	getEntries func(context.Context, []port.EntryRef) ([]port.EntryView, error)
}

func (s directoryReaderStub) GetEntries(
	ctx context.Context,
	refs []port.EntryRef,
) ([]port.EntryView, error) {
	return s.getEntries(ctx, refs)
}

func TestGetHomePageHandlerDoesNotReadDirectoryWithoutConversations(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	messagingReader := messagingReaderStub{
		listUserConversations: func(_ context.Context, actualUserID uuid.UUID) ([]port.ConversationSummary, error) {
			assert.Equal(t, userID, actualUserID)
			return nil, nil
		},
	}
	directoryReader := directoryReaderStub{
		getEntries: func(context.Context, []port.EntryRef) ([]port.EntryView, error) {
			t.Fatal("directory reader must not be called without entry references")
			return nil, nil
		},
	}
	handler := NewGetHomePageHandler(messagingReader, directoryReader)

	homePage, err := handler.Handle(t.Context(), GetHomePage{UserID: userID})

	require.NoError(t, err)
	assert.Empty(t, homePage.Peers)
}
