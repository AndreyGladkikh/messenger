package http_server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"messenger/messenger/internal/messaging/application/command/send_message"
	"messenger/messenger/internal/messaging/domain/message"
	"messenger/messenger/internal/platform/apperr"
	"messenger/messenger/internal/platform/commandbus"
	"messenger/messenger/internal/platform/db"
	"messenger/messenger/internal/platform/http_server/auth"
	"messenger/messenger/internal/platform/logger"
	"messenger/messenger/internal/platform/pubsub"
	"messenger/messenger/internal/platform/querybus"
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
	userID, _ := auth.UserIDFromContext(r.Context())
	sessionID, _ := auth.SessionIDFromContext(r.Context())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = auth.NewContextWithUserID(ctx, userID)
	ctx = auth.NewContextWithSessionID(ctx, sessionID)

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

	// wg.Go(func() {
	// 	// todo rename listenWS
	// 	if err := h.listenWSInbound(ctx, c); err != nil {
	// 		h.logger.ErrorContext(ctx, "websocket handler: ws inbound listener failed", "error", err.Error())
	// 	}
	// 	cancel()
	// })

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

	// wg.Go(func() {
	// 	// todo rename listenEvents
	// 	if err := h.listenSubscribers(ctx, c); err != nil {
	// 		h.logger.ErrorContext(ctx, "websocket handler: pubsub listener failed", "error", err.Error())
	// 	}
	// 	cancel()
	// })

	wg.Wait()

	if errors.Is(err, errConnTooSlow) {
		c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
		return
	}

	c.Close(websocket.StatusInternalError, "internal error")

	// for {
	// 	// for _, msg := range h.pubSubHub.Subscriptions() {

	// 	// }

	// 	var inbound WSMessage

	// 	err = wsjson.Read(context.Background(), c, &inbound)
	// 	if err != nil {
	// 		write(c, newResponse(inbound.ID, nil, err))
	// 		continue
	// 	}

	// 	var command command.Command
	// 	var query query.Query

	// 	switch inbound.Type {
	// 	case wsMessageTypeSendMessage:
	// 		// todo handle data type
	// 		mData := inbound.Data.(map[string]any)

	// 		chatID := uuid.MustParse(mData["chatId"].(string))
	// 		var replyToMessageID uuid.UUID
	// 		if v, ok := mData["replyToMessageId"].(string); ok {
	// 			replyToMessageID = uuid.MustParse(v)
	// 		}

	// 		command = &send_message.Command{
	// 			SenderID:         userID,
	// 			ChatID:           chatID,
	// 			MessageBody:      mData["body"].(string),
	// 			ReplyToMessageID: replyToMessageID,
	// 		}
	// 	default:
	// 		write(c, newResponse(inbound.ID, nil, errors.New("unknown message type")))
	// 		continue
	// 	}

	// 	if command != nil {
	// 		response, err := h.commandBus.Dispatch(context.Background(), command)
	// 		write(c, newResponse(inbound.ID, response, err))
	// 		continue
	// 	}

	// 	if query != nil {
	// 		response, err := h.queryBus.Dispatch(context.Background(), query)
	// 		write(c, newResponse(inbound.ID, response, err))
	// 		continue
	// 	}
	// }
}

func (h *WebsocketHandler) subscribeToChatEvents(ctx context.Context) (*pubsub.ChatEventsSubscription, error) {
	userID, _ := auth.UserIDFromContext(ctx)
	sessionID, _ := auth.SessionIDFromContext(ctx)

	subscription, err := h.pubsub.Subscription(
		ctx,
		fmt.Sprintf("%s:chat_events", sessionID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to chat events: %w", err)
	}
	defer subscription.Cancel(ctx)

	// todo refactor to app query
	userChatIDs, err := h.storage.Queries(ctx).ListChatIDsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(userChatIDs) > 0 {
		subscription.AddChats(ctx, userChatIDs...)
	}

	return subscription, nil
}

// func (h *WebsocketHandler) listenWSInbound(ctx context.Context, c *websocket.Conn) error {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return nil
// 		default:
// 			m, err := read(ctx, c)
// 			if err != nil {
// 				if errors.Is(err, context.Canceled) ||
// 					websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
// 					websocket.CloseStatus(err) == websocket.StatusGoingAway {
// 					return nil
// 				}
// 				h.logger.ErrorContext(ctx, "websocket handler: failed to read message", "error", err.Error())
// 				if errors.Is(err, errReadWSMessage) {
// 					return err
// 				}
// 				if errors.Is(err, errParseWSMessage) {
// 					write(ctx, c, newResponse(m.ID, nil, errors.New("invalid message format")))
// 					continue
// 				}
// 				write(ctx, c, newResponse(m.ID, nil, errors.New("internal error")))
// 				continue
// 			}
// 			response, err := h.processWSMessage(ctx, m)
// 			if err != nil {
// 				h.logger.ErrorContext(ctx, "websocket handler: failed to process message", "error", err.Error())
// 				write(ctx, c, newResponse(m.ID, nil, err))
// 				continue
// 			}
// 			write(ctx, c, newResponse(m.ID, response, nil))
// 		}
// 	}
// }

// todo renemae to processClientMessage
func (h *WebsocketHandler) processWSMessage(ctx context.Context, m WSMessage) (any, error) {
	userID, _ := auth.UserIDFromContext(ctx)

	switch m.Type {
	case wsMessageTypeSendMessage:
		data, err := decodeWSMessageDataAsType[SendMessageRequest](m.Data)
		if err != nil {
			return nil, err
		}

		command := &send_message.Command{
			SenderID:         userID,
			ChatID:           data.ChatID,
			MessageBody:      data.Body,
			ReplyToMessageID: data.ReplyToMessageID,
		}
		return h.commandBus.Dispatch(ctx, command)
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

// func (h *WebsocketHandler) listenWSInboundV1(ctx context.Context, c *websocket.Conn, userID, sessionID uuid.UUID) error {
// 	for {
// 		var inbound WSMessage

// 		err := wsjson.Read(ctx, c, &inbound)
// 		if err != nil {
// 			write(c, newResponse(inbound.ID, nil, err))
// 			continue
// 		}

// 		var command command.Command
// 		var query query.Query

// 		switch inbound.Type {
// 		case wsMessageTypeSendMessage:
// 			// todo handle data type
// 			mData := inbound.Data.(map[string]any)

// 			chatID := uuid.MustParse(mData["chatId"].(string))
// 			var replyToMessageID uuid.UUID
// 			if v, ok := mData["replyToMessageId"].(string); ok {
// 				replyToMessageID = uuid.MustParse(v)
// 			}

// 			command = &send_message.Command{
// 				SenderID:         userID,
// 				ChatID:           chatID,
// 				MessageBody:      mData["body"].(string),
// 				ReplyToMessageID: replyToMessageID,
// 			}
// 		case wsMessageTypeAddSubscription:
// 			// todo handle data type
// 			mData := inbound.Data.(map[string]any)
// 			sub := mData["sub"].(string)

// 			_, err = h.pubSubHub.Subscribe(ctx, fmt.Sprintf("%s:%s:%s", userID, sessionID, sub), []string{sub})
// 			if err != nil {
// 				return err
// 			}
// 		default:
// 			write(c, newResponse(inbound.ID, nil, errors.New("unknown message type")))
// 			continue
// 		}

// 		if command != nil {
// 			response, err := h.commandBus.Dispatch(ctx, command)
// 			write(c, newResponse(inbound.ID, response, err))
// 			continue
// 		}

// 		if query != nil {
// 			response, err := h.queryBus.Dispatch(ctx, query)
// 			write(c, newResponse(inbound.ID, response, err))
// 			continue
// 		}
// 	}
// }

// func (h *WebsocketHandler) listenSubscribers(ctx context.Context, c *websocket.Conn) error {
// 	for {
// 		for _, s := range h.pubSubHub.Subscriptions() {
// 			for _, msg := range s.Messages(0) {
// 				m, ok := msg.Payload.(map[string]any)
// 				if !ok {
// 					continue
// 				}
// 				mType, ok := m["type"]
// 				if !ok {
// 					continue
// 				}
// 				mData, ok := m["data"]
// 				if !ok {
// 					continue
// 				}

// 				switch mType {
// 				case pubsub.MessageTypeMessageSent:
// 					write(ctx, c, NewWSMessage(wsMessageTypeMessageSent, mData))
// 					continue
// 				default:
// 					continue
// 				}
// 			}
// 		}
// 	}
// }

func decodeWSMessageDataAsType[T any](data json.RawMessage) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, fmt.Errorf("failed to decode ws message data: %w", err)
	}
	return result, nil
}

// func read(ctx context.Context, c *wsConnn) (WSMessage, error) {
// 	var wsm WSMessage

// 	mt, m, err := c.Read(ctx)
// 	if err != nil {
// 		return wsm, fmt.Errorf("%w: %w", errReadWSMessage, err)
// 	}

// 	_=mt

// 	err = json.Unmarshal(m, &wsm)
// 	if err != nil {
// 		return wsm, fmt.Errorf("%w: failed to unmarshal JSON: %w", errParseWSMessage, err)
// 	}

// 	return wsm, nil
// }

// func write(ctx context.Context, c *wsConnn, msg any) error {
// 	// todo test timeout
// 	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
// 	defer cancel()

// 	return wsjson.Write(ctx, c.Conn, msg)
// }

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

// conn, loops

// type wsConnn struct {
// 	c *websocket.Conn
// 	send chan WSMessage
// }

// func newWSConnection(
// 	conn *websocket.Conn,
// ) *wsConnn {
// 	return &wsConnn{
// 		c: conn,
// 		send: make(chan WSMessage, 256),
// 	}
// }

type wsConnn struct {
	*websocket.Conn
	send chan WSMessage
}

func newWSConnection(
	w http.ResponseWriter,
	r *http.Request,
) (*wsConnn, error) {
	c, err := websocket.Accept(w, r, nil)
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

// func (c *wsConnn) closeSlow() {
// 	c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
// }

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
