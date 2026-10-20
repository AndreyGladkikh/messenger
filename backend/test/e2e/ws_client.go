package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type wsClient struct {

}

func newWSClient(t *testing.T, serverDSN string) *wsClient {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://%s/ws", serverDSN), nil)
	if err != nil {
		t.Fatalf("failed to establish ws connection: %v", err)
	}
	defer c.CloseNow()

	err = wsjson.Write(ctx, c, "hi")
	if err != nil {
		t.Fatalf("failed to send ws message: %v", err)
	}

	c.Close(websocket.StatusNormalClosure, "")

	return nil
}
