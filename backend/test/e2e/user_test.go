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
	profile    profile
	httpClient *httpClient
	wsClient   *wsClient
}

type profile struct {
	handle      string
	displayName string
	avatarURL   *string
}

func newUser(t *testing.T, server *httptest.Server, login string, opts ...userOption) *user {
	t.Helper()

	httpClient := newHttpClient(t, server)

	resp, err := httpClient.registerUser(t, login, "TestPass123!")
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	registerBody := decode[appHttp.Tokens](t, resp.Body)

	httpClient.JWT = &registerBody.AccessToken

	resp, err = httpClient.me(t)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	meBody := decode[appHttp.GetMe200JSONResponse](t, resp.Body)
	id, err := uuid.Parse(meBody.Id)
	require.NoError(t, err)

	wsClient := newWSClient(t, httpClient)
	t.Cleanup(func() {
		wsClient.CloseNow()
	})

	u := &user{
		id:         id,
		httpClient: httpClient,
		wsClient:   wsClient,
	}

	for _, o := range opts {
		o(u)
	}

	return u
}

type userOption func(*user)

func userWithProfile(t *testing.T, handle string, displayName string, avatarURL *string) userOption {
	return func(u *user) {
		resp, err := u.httpClient.createMyProfile(t, appHttp.CreateMyProfileJSONRequestBody{
			Handle:      handle,
			DisplayName: displayName,
			AvatarUrl:   avatarURL,
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		u.profile = profile{
			handle: handle,
			displayName: displayName,
			avatarURL: avatarURL,
		}
	}
}
