package websocket

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"messenger/messenger/internal/platform/auth"
	"messenger/messenger/internal/platform/commandbus"
	"messenger/messenger/internal/platform/logger"
	"messenger/messenger/internal/platform/pubsub"
	"messenger/messenger/internal/platform/querybus"
	"messenger/messenger/internal/platform/store"
	"net/http"
	"sync"
	"uuid"

	"github.com/coder/websocket"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	errReadWSMessage  = errors.New("websocket handler: failed to read message")
	errParseWSMessage = errors.New("websocket handler: failed to parse message")
	errConnTooSlow    = errors.New("websocket connection too slow to keep up with messages")
)

var _ http.Handler = (*Handler)(nil)

type Handler struct {
	Logger     *logger.Logger
	CommandBus *commandbus.Bus
	QueryBus   *querybus.Bus
	Store      *store.Store
	PubSub     *pubsub.PubSub
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a, err := auth.AuthFromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to accept websocket connection: unauthorized", http.StatusUnauthorized)
		return
	}

	reqID := middleware.GetReqID(r.Context())
	if reqID == "" {
		reqID = uuid.New().String()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = auth.NewContextWithAuth(ctx, a)
	ctx = context.WithValue(ctx, middleware.RequestIDKey, reqID)

	c, err := newWSConnection(w, r)
	if err != nil {
		// todo check if http server logs this error
		http.Error(w, "failed to accept websocket connection", http.StatusInternalServerError)
		return
	}
	defer c.CloseNow()

	eventsSubscription, cleanup, err := h.initEventsSubscription(ctx)
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to subscribe to events")
		c.Close(websocket.StatusInternalError, "internal error")
		return
	}
	defer cleanup()

	var (
		wg       sync.WaitGroup
		errCause error // an error that causes the ws connection to close
	)

	wg.Go(func() {
		if err := h.readLoop(ctx, c); err != nil {
			h.Logger.ErrorContext(ctx, "websocket handler: ws read loop failed", "error", err.Error())
			if errCause == nil {
				errCause = err
			}
		}
		cancel()
	})

	wg.Go(func() {
		if err := h.writeLoop(ctx, c); err != nil {
			h.Logger.ErrorContext(ctx, "websocket handler: ws write loop failed", "error", err.Error())
			if errCause == nil {
				errCause = err
			}
		}
		cancel()
	})

	wg.Go(func() {
		if err := h.readEventsLoop(ctx, eventsSubscription, c); err != nil {
			h.Logger.ErrorContext(ctx, "websocket handler: pub/sub read loop failed", "error", err.Error())
			if errCause == nil {
				errCause = err
			}
		}
		cancel()
	})

	wg.Wait()

	if errors.Is(err, errConnTooSlow) {
		c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
		return
	}

	c.Close(websocket.StatusInternalError, "internal error")
}

func (h *Handler) initEventsSubscription(ctx context.Context) (*pubsub.Subscription, func(), error) {
	a, _ := auth.AuthFromContext(ctx)
	reqID := middleware.GetReqID(ctx)

	subscription := h.PubSub.Subscription(fmt.Sprintf("ws:%s", reqID))
	cleanup := func() {
		subscription.Cancel(ctx)
	}

	pubsub.SubscribeToUserEvents(ctx, subscription, a.Principal.Subject.ID)

	// todo refactor to app query
	userCnvIDs, err := h.Store.Queries(ctx).ListConversationIDsForUser(ctx, a.Principal.Subject.ID)
	if err != nil {
		return nil, cleanup, err
	}

	if len(userCnvIDs) > 0 {
		pubsub.SubscribeToConversationEvents(ctx, subscription, userCnvIDs...)
	}

	return subscription, cleanup, nil
}

// func (h *Handler) subscribeToEvents(ctx context.Context) (*pubsub.ChatEventsSubscription, func(), error) {
// 	a, _ := auth.AuthFromContext(ctx)

// 	subscription, err := h.PeerEventsPubSub.Subscription(
// 		ctx,
// 		fmt.Sprintf("%s:chat_events", a.SessionID),
// 	)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("failed to subscribe to peer events: %w", err)
// 	}
// 	cleanup := func() {
// 		subscription.Cancel(ctx)
// 	}

// 	// todo refactor to app query
// 	userChatIDs, err := h.Store.Queries(ctx).ListConversationIDsForUser(ctx, a.Principal.Subject.ID)
// 	if err != nil {
// 		return nil, cleanup, err
// 	}

// 	if len(userChatIDs) > 0 {
// 		subscription.AddChats(ctx, userChatIDs...)
// 	}

// 	return subscription, cleanup, nil
// }

// todo renemae to processClientMessage
func (h *Handler) processWSMessage(ctx context.Context, m WSMessage) (any, error) {
	switch m.Type {
	case "Ping":
		return "Pong", nil
	default:
		return nil, fmt.Errorf("unknown message type: %s", m.Type)
	}
}

func decodeAs[T any](data jsontext.Value) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, fmt.Errorf("failed to decode ws message payload: %w", err)
	}
	return result, nil
}

type err struct {
}

// todo implement
func (e *err) Error() string {
	return "not implemented"
}

func mapErr(e error) *err {
	return &err{}
}

func withError(err error) wsMessageOption {
	return func(m *WSMessage) error {
		errData := mapErr(err)

		jsonErrorData, err := json.Marshal(errData)
		if err != nil {
			return err
		}

		m.Error = jsonErrorData
		return nil
	}
}

func (h *Handler) readLoop(ctx context.Context, c *wsConnn) error {
	for {
		m, err := c.read(ctx)
		if err != nil {
			if needLogWSReadError(err) {
				h.Logger.ErrorContext(ctx, "websocket handler: failed to read message", "message", m, "error", err.Error())
			}

			if errors.Is(err, errParseWSMessage) {
				continue
			}
			return err
		}

		response, err := h.processWSMessage(ctx, m)
		if err != nil {
			h.Logger.ErrorContext(ctx, "websocket handler: failed to process message", "error", err.Error())
		}

		m, err = newResponse(m.ID, response, err)
		if err != nil {
			return err
		}

		err = c.enqueueWrite(ctx, m)
		if err != nil {
			return err
		}
	}
}

func needLogWSReadError(err error) bool {
	return !errors.Is(err, context.Canceled) &&
		websocket.CloseStatus(err) != websocket.StatusNormalClosure &&
		websocket.CloseStatus(err) != websocket.StatusGoingAway
}

func (h *Handler) writeLoop(ctx context.Context, c *wsConnn) error {
	for {
		select {
		case m := <-c.send:
			err := c.write(ctx, m)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (h *Handler) readEventsLoop(ctx context.Context, subscription *pubsub.Subscription, c *wsConnn) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := subscription.Read(ctx)
			if err != nil {
				return err
			}

			payload, err := pubsub.DecodeMessage(msg)
			if err != nil {
				h.Logger.ErrorContext(ctx, "failed to read subscription message",
					"err", err,
					"msg", msg,
				)
				return err
			}

			switch p := payload.(type) {
			case pubsub.DirectConversationCreated:
				m, err := NewWSMessage(wsMessageTypeMessageSent, withData(DirectConversationStartedPayload{
					ConversationId:    p.ConversationID,
					SenderId:          p.SenderID,
					SenderDisplayName: p.SenderName,
					AvatarUrl:         p.SenderAvatarUrl,
					Message: &MessagePayload{
						Body:   p.MessageBody,
						SentAt: p.MessageSentAt,
					},
				}))
				if err != nil {
					return err
				}

				err = c.enqueueWrite(ctx, m)
				if err != nil {
					return err
				}
			}

			continue
		}
	}
}
