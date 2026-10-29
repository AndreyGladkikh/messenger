package e2e

import (
	"context"
	"encoding/json/v2"
	appHttp "messenger/messenger/internal/platform/transport/http"
	"messenger/messenger/internal/platform/transport/websocket"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendMessage_Success(t *testing.T) {
	server, _ := initApi(t)

	alice := newUser(t, server, "alice")
	defer alice.wsClient.CloseNow()

	bob := newUser(t, server, "bob")
	defer bob.wsClient.CloseNow()

	// err = wsjson.Write(ctx, c, "hi")
	// if err != nil {
	// 	t.Fatalf("failed to send ws message: %v", err)
	// }

	_, m, err := bob.wsClient.Read(context.Background())
	require.NoError(t, err)
	_ = m

	// alice.wsClient.sendMessage(websocket.SendFirstDirectMessagePayload{
	// 	Body: "Hi!",
	// 	RecipientID: bob.id,
	// })

	resp, err := alice.httpClient.createDirectConversationWithFirstMessage(appHttp.CreateDirectConversationJSONRequestBody{
		ParticipantId: bob.id,
		FirstMessage:  appHttp.FirstMessage{Body: "Hi!"},
	})
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	drctCnvBody := decode[appHttp.Conversation](t, resp.Body)
	assert.NotEmpty(t, drctCnvBody.Id)

	require.Eventually(t, func() bool {
		_, m, err := bob.wsClient.Read(context.Background())
		if err != nil {
			return false
		}

		var msg websocket.WSMessage
		err = json.Unmarshal(m, &msg)
		require.NoError(t, err)

		// var

		return true
	}, time.Second, 10*time.Millisecond)

	// alice.wsClient

	// bob.wsClient.sendMessage()
}
