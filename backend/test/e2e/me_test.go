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
		Handle:      handle,
		DisplayName: displayName,
		AvatarUrl:   &avatarRL,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, err = user.httpClient.getMyProfile(t)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, appHttp.Profile{
		Handle:      handle,
		DisplayName: displayName,
		AvatarUrl:   &avatarRL,
	}, decode[appHttp.Profile](t, resp.Body))
}

func TestCreateMyProfile_HandleAlreadyExists(t *testing.T) {
	server, _ := initApi(t)
	aliceA := newUser(t, server, "aliceA")
	aliceB := newUser(t, server, "aliceB")

	resp, _ := aliceA.httpClient.createMyProfile(t, appHttp.CreateMyProfileJSONRequestBody{
		Handle:      "alice",
		DisplayName: "Alice A",
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, _ = aliceB.httpClient.createMyProfile(t, appHttp.CreateMyProfileJSONRequestBody{
		Handle:      "alice",
		DisplayName: "Alice B",
	})
	require.Equal(t, http.StatusConflict, resp.StatusCode)
	problem := decode[appHttp.ProblemDetails](t, resp.Body)
	assert.Equal(t, problemTypesRegistry.AlreadyExists, *problem.Type)

	resp, _ = aliceB.httpClient.getMyProfile(t)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	resp, _ = aliceB.httpClient.createMyProfile(t, appHttp.CreateMyProfileJSONRequestBody{
		Handle:      "alice_b",
		DisplayName: "Alice B",
	})
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestCreateMyProfile_ProfileAlreadyExists(t *testing.T) {
	server, _ := initApi(t)
	user := newUser(t, server, "alice")

	handle := "alice"
	displayName := "Alice"

	request := appHttp.CreateMyProfileJSONRequestBody{
		Handle:      handle,
		DisplayName: displayName,
	}

	resp, _ := user.httpClient.createMyProfile(t, request)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, _ = user.httpClient.createMyProfile(t, request)
	require.Equal(t, http.StatusConflict, resp.StatusCode)
	problem := decode[appHttp.ProblemDetails](t, resp.Body)
	assert.Equal(t, problemTypesRegistry.AlreadyExists, *problem.Type)

	resp, _ = user.httpClient.getMyProfile(t)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	stored := decode[appHttp.Profile](t, resp.Body)
	assert.Equal(t, handle, stored.Handle)
	assert.Equal(t, displayName, stored.DisplayName)
}
