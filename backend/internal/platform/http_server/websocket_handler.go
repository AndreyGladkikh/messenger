package http_server

import (
	"context"
	"messenger/messenger/internal/platform/redis"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

type WebsocketHandler struct {
	pubSubHub *redis.RedisPubSubHub
}

func NewWebsocketHandler(
	pubSubHub *redis.RedisPubSubHub,
) *WebsocketHandler {
	return &WebsocketHandler{
		pubSubHub: pubSubHub,
	}
}

func (h *WebsocketHandler) handlerFunc(w http.ResponseWriter, r *http.Request) {
	// userID, _ := auth.UserIDFromContext(r.Context())

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer c.CloseNow()

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// h.pubSubHub.Subscribe(ctx, "chats:%s:", []string{})

	for {
		// for _, msg := range h.pubSubHub.Subscriptions() {

		// }

		msg, err := readTimeout(context.Background(), time.Second*3, c)
		if err != nil {
			return
		}
		_ = msg

		// err := writeTimeout(ctx, time.Second*5, c, msg)
		// if err != nil {
		// 	return err
		// }
	}


	// if errors.Is(err, context.Canceled) {
	// 	return
	// }
	// if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
	// 	websocket.CloseStatus(err) == websocket.StatusGoingAway {
	// 	return
	// }
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
}

// func (co *WebsocketHandler) subscribe(w http.ResponseWriter, r *http.Request) error {
// 	c, err := websocket.Accept(w, r, nil)
// 	if err != nil {
// 		return err
// 	}
// 	defer c.CloseNow()

// 	ctx := c.CloseRead(context.Background())

// 	co.pubSubHub.Subscribe(ctx, "chats:%s:", []string{})

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		default:
// 			for _, msg := range co.pubSubHub.Subscriptions() {

// 			}
// 			err := writeTimeout(ctx, time.Second*5, c, msg)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// }

func readTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn) (any, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	mt, body, err := c.Read(ctx)
	if err != nil {
		return nil, err
	}
	_, _ = mt, body
	return body, nil
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}