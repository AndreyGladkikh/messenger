package e2e

import (
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

	alice := newUser(t, server, "alice", userWithProfile(t, "alice_handle", "Alice DN", nil))
	bob := newUser(t, server, "bob", userWithProfile(t, "bob_handle", "Bob DN", nil))

	msgBody := "Hi!"

	resp, err := alice.httpClient.startDirectConversation(t, appHttp.StartDirectConversationJSONRequestBody{
		ParticipantId: bob.id,
		FirstMessage:  appHttp.FirstMessage{Body: msgBody},
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	drctCnvBody := decode[appHttp.Conversation](t, resp.Body)
	assert.NotEmpty(t, drctCnvBody.ConversationId)

	var msgRaw []byte
	select {
	case msgRaw = <-bob.wsClient.receive:
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for message")
	}

	var payload any
	err = json.Unmarshal(msgRaw, &payload)
	require.NoError(t, err)

	assert.Equal(t, websocket.DirectConversationStartedPayload{
		ConversationId:    drctCnvBody.ConversationId,
		SenderId:          alice.id,
		SenderDisplayName: alice.profile.displayName,
		AvatarUrl:         alice.profile.avatarURL,
		Message: &websocket.MessagePayload{
			Body: msgBody,
			SentAt: drctCnvBody.MessageSentAt,
		},
	}, decodeWSMsg[websocket.DirectConversationStartedPayload](t, msgRaw))
}
