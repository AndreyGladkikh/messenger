package e2e

import (
	"context"
	"fmt"
	"messenger/messenger/internal/platform/transport/http/cookies"
	// appWS "messenger/messenger/internal/platform/transport/websocket"
	"net/http"
	"testing"

	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type wsClient struct {
	*websocket.Conn
	receive chan []byte
}

func newWSClient(t *testing.T, httpClient *httpClient) *wsClient {
	ctx := context.Background()
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()

	c, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://%s/ws", "www.example.com"), &websocket.DialOptions{
		HTTPClient: httpClient.Client,
		HTTPHeader: http.Header{
			"Cookie": []string{fmt.Sprintf("%s=%s", cookies.AccessToken, *httpClient.JWT)},
		},
	})
	if err != nil {
		t.Fatalf("failed to establish ws connection: %v", err)
	}

	// err = wsjson.Write(ctx, c, "hi")
	// if err != nil {
	// 	t.Fatalf("failed to send ws message: %v", err)
	// }

	// c.Close(websocket.StatusNormalClosure, "")

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

func (c *wsClient) sendMessage(msg any) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return wsjson.Write(ctx, c.Conn, msg)
}
