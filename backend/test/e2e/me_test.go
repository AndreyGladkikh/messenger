package e2e

import (
	appHttp "messenger/messenger/internal/platform/transport/http"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateMyProfile_Success(t *testing.T) {
	server, _ := initApi(t)
	user := newUser(t, server, "alice")

	handle := "alice_handle"
	displayName := "Alice Display Name"
	avatarRL := "http://example.com/avatar"

	resp, err := user.httpClient.createMyProfile(t, appHttp.CreateMyProfileJSONRequestBody{
		Handle: handle,
		DisplayName: displayName,
		AvatarUrl: &avatarRL,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, err = user.httpClient.getMyProfile(t)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, appHttp.Profile{
		Handle: handle,
		DisplayName: displayName,
		AvatarUrl: &avatarRL,
	}, decode[appHttp.Profile](t, resp.Body))
}