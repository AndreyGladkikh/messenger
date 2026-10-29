package e2e

import (
	appHttp "messenger/messenger/internal/platform/transport/http"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/stretchr/testify/require"
)

type user struct {
	id         uuid.UUID
	httpClient *httpClient
	wsClient   *wsClient
}

func newUser(t *testing.T, server *httptest.Server, login string) *user {
	t.Helper()

	httpClient := newHttpClient(t, server)

	resp, err := httpClient.registerUser(login, "TestPass123!")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	registerBody := decode[appHttp.Tokens](t, resp.Body)

	httpClient.JWT = &registerBody.AccessToken

	resp, err = httpClient.me()
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	meBody := decode[appHttp.GetMe200JSONResponse](t, resp.Body)
	id, err := uuid.Parse(meBody.Id)
	require.NoError(t, err)

	wsClient := newWSClient(t, httpClient)

	return &user{
		id:         id,
		httpClient: httpClient,
		wsClient:   wsClient,
	}
}
