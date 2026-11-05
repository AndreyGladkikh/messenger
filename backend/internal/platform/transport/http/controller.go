package http

// import (
// 	"encoding/json"
// 	"errors"
// 	"messenger/messenger/internal/identity/app/command/login"
// 	"messenger/messenger/internal/identity/app/command/refresh"
// 	"messenger/messenger/internal/identity/app/command/register"
// 	"messenger/messenger/internal/identity/app/query/get_current_user"
// 	"messenger/messenger/internal/identity/app/token"
// 	authDomain "messenger/messenger/internal/identity/domain"
// 	"messenger/messenger/internal/identity/domain/session"
// 	"messenger/messenger/internal/messaging/application/command/create_private_chat"
// 	"messenger/messenger/internal/messaging/application/command/send_message"
// 	"messenger/messenger/internal/messaging/application/query/get_chat_list"
// 	"messenger/messenger/internal/platform/commandbus"
// 	"messenger/messenger/internal/platform/transport/http/auth"
// 	"messenger/messenger/internal/platform/transport/http/cookies"
// 	"net/http"
// 	"net/netip"
// 	"time"

// 	"github.com/go-chi/chi/v5/middleware"
// )

// type Controller struct {
// 	commandBus *commandbus.Bus
// }

// func NewController(
// 	commandBus *commandbus.Bus,
// ) *Controller {
// 	return &Controller{
// 		commandBus: commandBus,
// 	}
// }

// func (c *Controller) registerUser(w http.ResponseWriter, r *http.Request) {
// 	var request RegisterUserRequest
// 	json.NewDecoder(r.Body).Decode(&request)

// 	clientIP := middleware.GetClientIP(r.Context())
// 	ip, err := netip.ParseAddr(clientIP)
// 	if err != nil {
// 		NewResponse(nil, err).WriteTo(w)
// 		return
// 	}
// 	userAgend := r.Header.Get("User-Agent")

// 	command := &register.Command{
// 		Login:     request.Login,
// 		Password:  request.Password,
// 		IP:        ip,
// 		UserAgent: userAgend,
// 	}
// 	response, err := c.commandBus.Dispatch(r.Context(), command)

// 	if err == nil {
// 		setAuthCookies(w, response.(*register.Response).AccessToken, response.(*register.Response).RefreshToken)
// 	}

// 	NewResponse(response, err).WriteTo(w)
// }

// func (c *Controller) loginUser(w http.ResponseWriter, r *http.Request) {
// 	var request LoginRequest
// 	json.NewDecoder(r.Body).Decode(&request)

// 	clientIP := middleware.GetClientIP(r.Context())
// 	ip, err := netip.ParseAddr(clientIP)
// 	if err != nil {
// 		NewResponse(nil, err).WriteTo(w)
// 		return
// 	}
// 	userAgend := r.Header.Get("User-Agent")

// 	command := &login.Command{
// 		Login:     request.Login,
// 		Password:  request.Password,
// 		IP:        ip,
// 		UserAgent: userAgend,
// 	}
// 	response, err := c.commandBus.Dispatch(r.Context(), command)

// 	if err == nil {
// 		setAuthCookies(w, response.(*login.Response).AccessToken, response.(*login.Response).RefreshToken)
// 	}

// 	NewResponse(response, err).WriteTo(w)
// }

// func (c *Controller) refreshSession(w http.ResponseWriter, r *http.Request) {
// 	refreshToken, err := r.Cookie(cookies.RefreshToken)
// 	if err != nil {
// 		if errors.Is(err, http.ErrNoCookie) {
// 			err = authDomain.ErrUnauthorized
// 		}
// 		NewResponse(nil, err).WriteTo(w)
// 		return
// 	}

// 	command := &refresh.Command{
// 		RefreshToken: refreshToken.Value,
// 	}
// 	response, err := c.commandBus.Dispatch(r.Context(), command)

// 	if err == nil {
// 		setAuthCookies(w, response.(*refresh.Response).AccessToken, response.(*refresh.Response).RefreshToken)
// 	}

// 	NewResponse(response, err).WriteTo(w)
// }

// func setAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) {
// 	accessTokenCookie := &http.Cookie{
// 		Name:     cookies.AccessToken,
// 		Value:    accessToken,
// 		Path:     "/",
// 		Expires:  time.Now().Add(token.AccessTokenTTL),
// 		Secure:   true,
// 		SameSite: http.SameSiteLaxMode,
// 	}
// 	refreshTokenCookie := &http.Cookie{
// 		Name:     cookies.RefreshToken,
// 		Value:    refreshToken,
// 		Path:     "/",
// 		Expires:  time.Now().Add(session.TTL),
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteLaxMode,
// 	}

// 	http.SetCookie(w, accessTokenCookie)
// 	http.SetCookie(w, refreshTokenCookie)
// }

// func (c *Controller) sendMessage(w http.ResponseWriter, r *http.Request) {
// 	var request Message
// 	json.NewDecoder(r.Body).Decode(&request)

// 	senderID, _ := auth.UserIDFromContext(r.Context())

// 	command := &send_message.Command{
// 		SenderID:         senderID,
// 		ChatID:           request.ChatID,
// 		MessageBody:      request.Body,
// 		ReplyToMessageID: request.ReplyToMessageID,
// 		Attachments:      request.Attachments,
// 	}
// 	response, err := c.commandBus.Dispatch(r.Context(), command)

// 	NewResponse(response, err).WriteTo(w)
// }

// func (c *Controller) createPrivateChat(w http.ResponseWriter, r *http.Request) {
// 	var request CreatePrivateChatRequest
// 	json.NewDecoder(r.Body).Decode(&request)

// 	userID, _ := auth.UserIDFromContext(r.Context())

// 	command := &create_private_chat.Command{
// 		InitiatorID:    userID,
// 		ChatWithUserID: request.ChatWithUserID,
// 	}
// 	response, err := c.commandBus.Dispatch(r.Context(), command)

// 	NewResponse(response, err).WriteTo(w)
// }
