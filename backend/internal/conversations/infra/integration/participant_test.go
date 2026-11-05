package integration

import (
	"context"
	"errors"
	idQry "messenger/messenger/internal/identity/app/query"
	"messenger/messenger/internal/shared/app/apperr"
	appQuery "messenger/messenger/internal/shared/app/query"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParticipantLookupExists(t *testing.T) {
	t.Parallel()

	participantID := uuid.New()
	handler := appQuery.HandlerFunc[idQry.GetAccount, idQry.AccountView](
		func(_ context.Context, qry idQry.GetAccount) (idQry.AccountView, error) {
			assert.Equal(t, participantID, qry.ID)
			return idQry.AccountView{ID: participantID, Type: "user"}, nil
		},
	)

	exists, err := NewParticipantLookup(handler).Exists(t.Context(), participantID)

	require.NoError(t, err)
	assert.True(t, exists)
}

func TestParticipantLookupTreatsNotFoundAsMissing(t *testing.T) {
	t.Parallel()

	handler := appQuery.HandlerFunc[idQry.GetAccount, idQry.AccountView](
		func(context.Context, idQry.GetAccount) (idQry.AccountView, error) {
			return idQry.AccountView{}, apperr.NewNotFound("account not found")
		},
	)

	exists, err := NewParticipantLookup(handler).Exists(t.Context(), uuid.New())

	require.NoError(t, err)
	assert.False(t, exists)
}

func TestParticipantLookupPreservesUnexpectedError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("identity unavailable")
	handler := appQuery.HandlerFunc[idQry.GetAccount, idQry.AccountView](
		func(context.Context, idQry.GetAccount) (idQry.AccountView, error) {
			return idQry.AccountView{}, expectedErr
		},
	)

	exists, err := NewParticipantLookup(handler).Exists(t.Context(), uuid.New())

	assert.False(t, exists)
	require.ErrorIs(t, err, expectedErr)
}
