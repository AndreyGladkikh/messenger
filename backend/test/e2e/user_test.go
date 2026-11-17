package e2e

import (
	"bytes"
	appHttp "messenger/messenger/internal/platform/transport/http"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func alice(t *testing.T, server *httptest.Server) *user {
	return newUser(t, server, "alice", userWithProfile(t, "alice_handle", "Alice DN", true))
}

func bob(t *testing.T, server *httptest.Server) *user {
	return newUser(t, server, "bob", userWithProfile(t, "bob_handle", "Bob DN", true))
}

type userOption func(*user)

func userWithProfile(t *testing.T, handle string, displayName string, withAvatar bool) userOption {
	return func(u *user) {
		var avatarURL *string
		var avatarMediaID *uuid.UUID
		if withAvatar {
			response, err := u.httpClient.initMediaUpload(t, appHttp.InitMediaUploadJSONRequestBody{
				FileName: "avatar.png",
				Purpose:  "avatar",
			})
			require.NoError(t, err)
			require.Equal(t, http.StatusCreated, response.StatusCode)
			upload := decode[appHttp.InitedMediaUpload](t, response.Body)

			request, err := http.NewRequestWithContext(t.Context(), http.MethodPut, upload.PresignedUrl, bytes.NewReader([]byte("avatar")))
			require.NoError(t, err)
			response, err = (&http.Client{Timeout: 10 * time.Second}).Do(request)
			require.NoError(t, err)
			response.Body.Close()
			require.Equal(t, http.StatusOK, response.StatusCode)

			response, err = u.httpClient.completeMediaUpload(t, appHttp.CompleteMediaUploadRequestObject{UploadId: upload.Id})
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, response.StatusCode)
			asset := decode[appHttp.Media](t, response.Body)

			avatarMediaID = &asset.MediaId
			url := testCDNHost + "/public/" + asset.MediaId.String() + "/original"
			avatarURL = &url
		}

		resp, err := u.httpClient.createMyProfile(t, appHttp.CreateMyProfileJSONRequestBody{
			Handle:        handle,
			DisplayName:   displayName,
			AvatarMediaId: avatarMediaID,
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
