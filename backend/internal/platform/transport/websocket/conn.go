package websocket

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type wsConnn struct {
	*websocket.Conn
	send chan WSMessage
}

func newWSConnection(w http.ResponseWriter, r *http.Request) (*wsConnn, error) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:5173"},
	})
	if err != nil {
		return nil, err
	}

	return &wsConnn{
		Conn: c,
		send: make(chan WSMessage, 256),
	}, nil
}

func (c *wsConnn) read(ctx context.Context) (WSMessage, error) {
	var wsm WSMessage

	_, m, err := c.Conn.Read(ctx)
	if err != nil {
		return wsm, err
	}

	err = json.Unmarshal(m, &wsm)
	if err != nil {
		return wsm, fmt.Errorf("%w: failed to unmarshal JSON: %w", errParseWSMessage, err)
	}

	// todo validate message
	if wsm.ID == "" || wsm.Type == "" {
		return wsm, fmt.Errorf("%w: message id or type missing", errParseWSMessage)
	}

	return wsm, nil
}

func (c *wsConnn) write(ctx context.Context, msg any) error {
	// todo test timeout
	// todo set timeout 10 sec
	ctx, cancel := context.WithTimeout(ctx, time.Second*10000000)
	defer cancel()

	return wsjson.Write(ctx, c.Conn, msg)
}

func (c *wsConnn) enqueueWrite(ctx context.Context, m WSMessage) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c.send <- m:
		return nil
	default:
		close(c.send)
		// c.closeSlow()
		return errConnTooSlow
	}
}
