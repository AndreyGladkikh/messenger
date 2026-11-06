package e2e

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"messenger/messenger/internal/platform/transport/http/cookies"
	appWS "messenger/messenger/internal/platform/transport/websocket"
	"net/http"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/stretchr/testify/require"
)

type wsClient struct {
	*websocket.Conn
	receive chan []byte
}

func newWSClient(t *testing.T, httpClient *httpClient) *wsClient {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://%s/ws", "www.example.com"), &websocket.DialOptions{
		HTTPClient: httpClient.Client,
		HTTPHeader: http.Header{
			"Cookie": []string{fmt.Sprintf("%s=%s", cookies.AccessToken, *httpClient.JWT)},
		},
	})
	if err != nil {
		t.Fatalf("failed to establish ws connection: %v", err)
	}

	receive := make(chan []byte, 32)
	go func() {
		for {
			_, m, err := c.Read(context.Background())
			if err != nil {
				c.CloseNow()
				return
			}
			receive <- m
		}
	}()

	return &wsClient{
		Conn:    c,
		receive: receive,
	}
}

func decodeWSMsg[T any](t *testing.T, data []byte) T {
	t.Helper()
	
	var msg appWS.WSMessage
	err := json.Unmarshal(data, &msg)
	require.NoError(t, err)
	require.NotEmpty(t, msg.Type)
	require.NotEmpty(t, msg.Payload)

	var payload T
	err = json.Unmarshal(msg.Payload, &payload)
	require.NoError(t, err)
	
	return payload
}
