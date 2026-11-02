package e2e

import (
	"encoding/json/v2"
	appHttp "messenger/messenger/internal/platform/transport/http"
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

	resp, err := alice.httpClient.registerDirectoryEntry(appHttp.DirectoryEntry{
		DisplayName: "Alice DN",
		Handle:      "alice",
		TargetId:    alice.id,
		TargetType:  "user",
	})
	require.NoError(t, err)
	defer resp.Body.Close()

	bob := newUser(t, server, "bob")
	defer bob.wsClient.CloseNow()

	resp, err = bob.httpClient.registerDirectoryEntry(appHttp.DirectoryEntry{
		DisplayName: "Bob DN",
		Handle:      "bob",
		TargetId:    bob.id,
		TargetType:  "user",
	})
	require.NoError(t, err)
	defer resp.Body.Close()

	resp, err = alice.httpClient.startDirectConversation(appHttp.StartDirectConversationJSONRequestBody{
		ParticipantId: bob.id,
		FirstMessage:  appHttp.FirstMessage{Body: "Hi!"},
	})
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	drctCnvBody := decode[appHttp.Conversation](t, resp.Body)
	assert.NotEmpty(t, drctCnvBody.Id)

	var msg []byte

	select {
	case msg = <-bob.wsClient.receive:
		a := 1
		_=a
		_=msg
	case <-time.After(1*time.Second):
		t.Fatal("timeout waiting for message")
	}

	var payload any
	err = json.Unmarshal(msg, &payload)
	require.NoError(t, err)
}
