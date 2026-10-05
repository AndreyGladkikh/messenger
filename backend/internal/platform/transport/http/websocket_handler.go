package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	authDomain "messenger/messenger/internal/auth/domain"
	"messenger/messenger/internal/conversations/app/command"
	"messenger/messenger/internal/conversations/domain/model/message"
	"messenger/messenger/internal/platform/auth"
	"messenger/messenger/internal/platform/commandbus"
	"messenger/messenger/internal/platform/db"
	"messenger/messenger/internal/platform/logger"
	"messenger/messenger/internal/platform/pubsub"
	"messenger/messenger/internal/platform/querybus"
	"messenger/messenger/internal/shared/application/apperr"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

var (
	errReadWSMessage  = errors.New("websocket handler: failed to read message")
	errParseWSMessage = errors.New("websocket handler: failed to parse message")
	errConnTooSlow    = errors.New("websocket connection too slow to keep up with messages")
)

type WebsocketHandler struct {
	logger     *logger.Logger
	commandBus *commandbus.Bus
	queryBus   *querybus.Bus
	storage    *db.Storage
	pubsub     *pubsub.ChatEventsPubSub
}

func NewWebsocketHandler(
	logger *logger.Logger,
	commandBus *commandbus.Bus,
	queryBus *querybus.Bus,
	storage *db.Storage,
	pubsub *pubsub.ChatEventsPubSub,
) *WebsocketHandler {
	return &WebsocketHandler{
		logger:     logger,
		commandBus: commandBus,
		queryBus:   queryBus,
		storage:    storage,
		pubsub:     pubsub,
	}
}

func (h *WebsocketHandler) handlerFunc(w http.ResponseWriter, r *http.Request) {
	a, err := auth.AuthFromContext(r.Context())
	if err != nil {
		err = fmt.Errorf("%w: %w", authDomain.ErrUnauthorized, err)
		http.Error(w, "failed to accept websocket connection: unauthorized", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = auth.NewContextWithAuth(ctx, a)

	c, err := newWSConnection(w, r)
	if err != nil {
		http.Error(w, "failed to accept websocket connection", http.StatusInternalServerError)
		return
	}
	defer c.CloseNow()

	chatEventsSubscription, err := h.subscribeToChatEvents(ctx)
	if err != nil {
		c.Close(websocket.StatusInternalError, "internal error")
		return
	}

	var (
		wg       sync.WaitGroup
		errCause error // an error that causes the ws connection to close
	)

	wg.Go(func() {
		if err := h.readLoop(ctx, c); err != nil {
			h.logger.ErrorContext(ctx, "websocket handler: ws read loop failed", "error", err.Error())
			if errCause == nil {
				errCause = err
			}
		}
		cancel()
	})

	wg.Go(func() {
		if err := h.writeLoop(ctx, c); err != nil {
			h.logger.ErrorContext(ctx, "websocket handler: ws write loop failed", "error", err.Error())
			if errCause == nil {
				errCause = err
			}
		}
		cancel()
	})

	wg.Go(func() {
		if err := h.readChatEventsLoop(ctx, chatEventsSubscription, c); err != nil {
			h.logger.ErrorContext(ctx, "websocket handler: pub/sub read loop failed", "error", err.Error())
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

func (h *WebsocketHandler) subscribeToChatEvents(ctx context.Context) (*pubsub.ChatEventsSubscription, error) {
	a, _ := auth.AuthFromContext(ctx)

	subscription, err := h.pubsub.Subscription(
		ctx,
		fmt.Sprintf("%s:chat_events", a.SessionID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to chat events: %w", err)
	}
	defer subscription.Cancel(ctx)

	// todo refactor to app query
	userChatIDs, err := h.storage.Queries(ctx).ListConversationIDsForUser(ctx, a.Principal.Subject.ID)
	if err != nil {
		return nil, err
	}

	if len(userChatIDs) > 0 {
		subscription.AddChats(ctx, userChatIDs...)
	}

	return subscription, nil
}

// todo renemae to processClientMessage
func (h *WebsocketHandler) processWSMessage(ctx context.Context, m WSMessage) (any, error) {
	a, _ := auth.AuthFromContext(ctx)
	// userID, _ := auth.UserIDFromContext(ctx)

	switch m.Type {
	case wsMessageTypeSendMessage:
		data, err := decodeWSMessageDataAsType[SendMessageRequest](m.Data)
		if err != nil {
			return nil, err
		}

		command := command.SendMessage{
			SenderID:         a.Principal.Subject.ID,
			ConversationID:           data.ChatID,
			Body:      data.Body,
		}
		return h.commandBus.Dispatch[any](ctx, command)
	// case wsMessageTypeAddSubscription:
	// 	data, err := decodeWSMessageDataAsType[map[string]any](m.Data)
	// 	if err != nil {
	// 		// todo handle
	// 		return nil, err
	// 	}
	// 	sub := data["sub"].(string)

	// 	_, err = h.pubSubHub.Subscribe(ctx, fmt.Sprintf("%s:%s:%s", userID, sessionID, sub), []string{sub})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	default:
		return nil, fmt.Errorf("unknown message type: %s", m.Type)
	}
}

func decodeWSMessageDataAsType[T any](data json.RawMessage) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, fmt.Errorf("failed to decode ws message data: %w", err)
	}
	return result, nil
}

type WSMessageType string

const (
	wsMessageTypeResponse WSMessageType = "response"

	wsMessageTypeSendMessage     WSMessageType = "sendMessage"
	wsMessageTypeAddSubscription WSMessageType = "addSubscription"

	wsMessageTypeMessageSent WSMessageType = "messageSent"
)

type WSMessage struct {
	ID    string          `json:"id,omitempty"`
	Type  WSMessageType   `json:"type"`
	Data  json.RawMessage `json:"data,omitempty"`
	Error json.RawMessage `json:"error,omitempty"`
}

func NewWSMessage(typ WSMessageType, opts ...wsMessageOption) (WSMessage, error) {
	m := WSMessage{
		Type: typ,
	}

	for _, opt := range opts {
		err := opt(&m)
		if err != nil {
			return WSMessage{}, err
		}
	}

	return m, nil
}

func newResponse(id string, data any, err error) (WSMessage, error) {
	opts := []wsMessageOption{withID(id)}
	if data != nil {
		opts = append(opts, withData(data))
	}
	if err != nil {
		opts = append(opts, withError(err))
	}

	return NewWSMessage(
		wsMessageTypeResponse,
		opts...,
	)
}

type wsMessageOption func(*WSMessage) error

func withID(id string) wsMessageOption {
	return func(m *WSMessage) error {
		m.ID = id
		return nil
	}
}

func withData(data any) wsMessageOption {
	return func(m *WSMessage) error {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}

		m.Data = jsonData
		return nil
	}
}

func withError(err error) wsMessageOption {
	return func(m *WSMessage) error {
		errData := apperr.Translate(err)

		jsonErrorData, err := json.Marshal(errData)
		if err != nil {
			return err
		}

		m.Error = jsonErrorData
		return nil
	}
}

type wsConnn struct {
	*websocket.Conn
	send chan WSMessage
}

func newWSConnection(
	w http.ResponseWriter,
	r *http.Request,
) (*wsConnn, error) {
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

	_, m, err := c.Read(ctx)
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
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
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

func (h *WebsocketHandler) readLoop(ctx context.Context, c *wsConnn) error {
	for {
		m, err := c.read(ctx)
		if err != nil {
			if needLogWSReadError(err) {
				h.logger.ErrorContext(ctx, "websocket handler: failed to read message", "message", m, "error", err.Error())
			}

			if errors.Is(err, errParseWSMessage) {
				continue
			}
			return err
		}

		response, err := h.processWSMessage(ctx, m)
		if err != nil {
			h.logger.ErrorContext(ctx, "websocket handler: failed to process message", "error", err.Error())
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

func (h *WebsocketHandler) writeLoop(ctx context.Context, c *wsConnn) error {
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

func (h *WebsocketHandler) readChatEventsLoop(ctx context.Context, chatEventsSubscription *pubsub.ChatEventsSubscription, c *wsConnn) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			eventEnvelope, err := chatEventsSubscription.Read(ctx)
			if err != nil {
				return err
			}

			switch eventEnvelope.Event.(type) {
			case message.MessageSent:
				m, err := NewWSMessage(wsMessageTypeMessageSent, withData(eventEnvelope))
				if err != nil {
					return err
				}
				err = c.enqueueWrite(ctx, m)
				if err != nil {
					return err
				}
				continue
			default:
				continue
			}
		}
	}
}
