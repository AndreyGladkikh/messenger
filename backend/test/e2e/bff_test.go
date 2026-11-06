package e2e

import (
	"messenger/messenger/internal/messaging/domain/model/conversation"
	appHttp "messenger/messenger/internal/platform/transport/http"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHomePage_Success(t *testing.T) {
	server, _ := initApi(t)

	alice := alice(t, server)
	bob := bob(t, server)

	resp, _ := alice.httpClient.getHomePage(t)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	aliceHomePage := decode[appHttp.HomePage](t, resp.Body)
	assert.Empty(t, aliceHomePage.Peers)

	resp, _ = bob.httpClient.getHomePage(t)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bobHomePage := decode[appHttp.HomePage](t, resp.Body)
	assert.Empty(t, bobHomePage.Peers)

	msgBody := "Hi!"
	resp, _ = alice.httpClient.startDirectConversation(t, appHttp.DirectConversation{
		ParticipantId: bob.id,
		FirstMessage: appHttp.FirstMessage{
			Body: msgBody,
		},
	})
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	cnv := decode[appHttp.Conversation](t, resp.Body)

	unreadC := 0

	resp, _ = alice.httpClient.getHomePage(t)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	aliceHomePage = decode[appHttp.HomePage](t, resp.Body)
	assert.Equal(t, appHttp.HomePage{
		Peers: []appHttp.PeerItem{
			{
				Id: cnv.ConversationId,
				Preview: &msgBody,
				PreviewTime: &cnv.MessageSentAt,
				Title: bob.profile.displayName,
				Type: conversation.TypeDirect.String(),
				AvatarUrl: bob.profile.avatarURL,
				UnreadCount: &unreadC,
			},
		},
	}, aliceHomePage)

	resp, _ = bob.httpClient.getHomePage(t)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bobHomePage = decode[appHttp.HomePage](t, resp.Body)
	assert.Equal(t, appHttp.HomePage{
		Peers: []appHttp.PeerItem{
			{
				Id: cnv.ConversationId,
				Preview: &msgBody,
				PreviewTime: &cnv.MessageSentAt,
				Title: alice.profile.displayName,
				Type: conversation.TypeDirect.String(),
				AvatarUrl: alice.profile.avatarURL,
				UnreadCount: &unreadC,
			},
		},
	}, bobHomePage)
}
