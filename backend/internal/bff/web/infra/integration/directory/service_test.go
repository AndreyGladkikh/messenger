package directory

import (
	"context"
	"errors"
	"messenger/messenger/internal/bff/web/application/port"
	drctrQry "messenger/messenger/internal/directory/app/query"
	appQuery "messenger/messenger/internal/shared/app/query"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceGetEntries(t *testing.T) {
	t.Parallel()

	entryID := uuid.New()
	avatarURL := "https://example.com/avatar.png"
	handler := appQuery.HandlerFunc[drctrQry.GetEntries, []drctrQry.EntryView](
		func(_ context.Context, qry drctrQry.GetEntries) ([]drctrQry.EntryView, error) {
			assert.Equal(t, []drctrQry.EntryRef{{
				ID:   entryID,
				Type: port.DirectoryEntryTypeUser,
			}}, qry.Refs)

			return []drctrQry.EntryView{{
				ID:          entryID,
				Type:        port.DirectoryEntryTypeUser,
				DisplayName: "Alice",
				AvatarURL:   &avatarURL,
			}}, nil
		},
	)
	service := NewService(handler)

	entries, err := service.GetEntries(t.Context(), []port.EntryRef{{
		ID:   entryID,
		Type: port.DirectoryEntryTypeUser,
	}})

	require.NoError(t, err)
	assert.Equal(t, []port.EntryView{{
		ID:          entryID,
		Type:        port.DirectoryEntryTypeUser,
		DisplayName: "Alice",
		AvatarURL:   &avatarURL,
	}}, entries)
}

func TestServiceGetEntriesWrapsHandlerError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("read entries")
	handler := appQuery.HandlerFunc[drctrQry.GetEntries, []drctrQry.EntryView](
		func(context.Context, drctrQry.GetEntries) ([]drctrQry.EntryView, error) {
			return nil, expectedErr
		},
	)
	service := NewService(handler)

	_, err := service.GetEntries(t.Context(), nil)

	require.ErrorIs(t, err, expectedErr)
	assert.ErrorContains(t, err, "get direct entries")
}
