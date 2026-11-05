package e2e

import (
	// appHttp "messenger/messenger/internal/platform/transport/http"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHomePage_Success(t *testing.T) {
	server, _ := initApi(t)

	user := newUser(t, server, "alice")
	defer user.wsClient.CloseNow()

	resp, err := user.httpClient.getHomePage()
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)	
	
	// body := decode[appHttp.GetWebHomePage200JSONResponse](t, resp.Body)
	// assert.NotEmpty(t, body.Peers)
}
