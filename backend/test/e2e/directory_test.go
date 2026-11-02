package e2e

import (
	appHttp "messenger/messenger/internal/platform/transport/http"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterDirectoryEntry_Success(t *testing.T) {
	server, _ := initApi(t)

	user := newUser(t, server, "alice")
	defer user.wsClient.CloseNow()

	resp, err := user.httpClient.registerDirectoryEntry(appHttp.DirectoryEntry{
		DisplayName: "Alice DN",
		Handle:      "alice",
		TargetId:    user.id,
		TargetType:  "user",
	})
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
