package identity

import (
	"context"
	"errors"
	idQry "messenger/messenger/internal/identity/app/query"
	appQuery "messenger/messenger/internal/shared/app/query"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserLookupExists(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	handler := appQuery.HandlerFunc[idQry.GetAccount, idQry.AccountView](
		func(_ context.Context, qry idQry.GetAccount) (idQry.AccountView, error) {
			assert.Equal(t, userID, qry.ID)
			return idQry.AccountView{ID: userID, Type: "user"}, nil
		},
	)

	exists, err := NewUserLookup(handler).Exists(t.Context(), userID)

	require.NoError(t, err)
	assert.True(t, exists)
}

func TestUserLookupRejectsNonUserAccount(t *testing.T) {
	t.Parallel()

	handler := appQuery.HandlerFunc[idQry.GetAccount, idQry.AccountView](
		func(context.Context, idQry.GetAccount) (idQry.AccountView, error) {
			return idQry.AccountView{Type: "bot"}, nil
		},
	)

	exists, err := NewUserLookup(handler).Exists(t.Context(), uuid.New())

	require.NoError(t, err)
	assert.False(t, exists)
}

func TestUserLookupPreservesHandlerError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("identity unavailable")
	handler := appQuery.HandlerFunc[idQry.GetAccount, idQry.AccountView](
		func(context.Context, idQry.GetAccount) (idQry.AccountView, error) {
			return idQry.AccountView{}, expectedErr
		},
	)

	exists, err := NewUserLookup(handler).Exists(t.Context(), uuid.New())

	assert.False(t, exists)
	require.ErrorIs(t, err, expectedErr)
}
