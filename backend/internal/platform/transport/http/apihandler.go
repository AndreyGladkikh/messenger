package http

import (
	"context"
	"encoding/json/v2"
	"fmt"
	cnvCommand "messenger/messenger/internal/conversations/app/command"
	idCommand "messenger/messenger/internal/identity/app/command"
	"messenger/messenger/internal/identity/app/query"
	"messenger/messenger/internal/identity/app/token"
	authDomain "messenger/messenger/internal/identity/domain"
	"messenger/messenger/internal/identity/domain/session"
	"messenger/messenger/internal/platform/auth"
	"messenger/messenger/internal/platform/commandbus"
	"messenger/messenger/internal/platform/querybus"
	"messenger/messenger/internal/platform/transport/http/cookies"
	"messenger/messenger/internal/shared/application/apperr"
	"net/http"
	"net/netip"
	"time"
	"uuid"

	"github.com/go-chi/chi/v5/middleware"
)

var _ StrictServerInterface = (*ApiHandler)(nil)

type ApiHandler struct {
	commandBus *commandbus.Bus
	queryBus   *querybus.Bus
}

func NewApiHandler(
	commandBus *commandbus.Bus,
	queryBus *querybus.Bus,
) *ApiHandler {
	return &ApiHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
	}
}

func (h *ApiHandler) GetMe(ctx context.Context, request GetMeRequestObject) (GetMeResponseObject, error) {
	a, err := getAuth(ctx)
	if err != nil {
		return nil, err
	}

	q := query.GetAccount{
		ID: a.Principal.Subject.ID,
	}
	response, err := h.queryBus.Dispatch[query.AccountView](ctx, q)
	if err != nil {
		return nil, err
	}

	return GetMe200JSONResponse{
		Id:   response.ID.String(),
		Type: response.Type,
	}, nil
}

func (h *ApiHandler) GetMyChatList(ctx context.Context, request GetMyChatListRequestObject) (GetMyChatListResponseObject, error) {
	_, err := getAuth(ctx)
	if err != nil {
		return nil, err
	}

	// query := &get_chat_list.Query{
	// 	UserID: a.Principal.Subject.ID,
	// }
	// response, err := h.queryBus.Dispatch[get_chat_list.ChatListView](ctx, query)
	// if err != nil {
	// 	return nil, err
	// }

	return GetMyChatList200JSONResponse{}, nil
}

func (h *ApiHandler) Login(ctx context.Context, request LoginRequestObject) (LoginResponseObject, error) {
	clientIP := middleware.GetClientIP(ctx)
	ip, err := netip.ParseAddr(clientIP)
	if err != nil {
		return nil, err
	}

	command := idCommand.LoginUser{
		Login:     request.Body.Login,
		Password:  request.Body.Password,
		IP:        ip,
		UserAgent: request.Params.UserAgent,
	}
	response, err := h.commandBus.Dispatch[idCommand.LoginUserResponse](ctx, command)
	if err != nil {
		return nil, err
	}

	return Login200JSONResponseCustom{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}, nil
}

func (h *ApiHandler) RefreshTokens(ctx context.Context, request RefreshTokensRequestObject) (RefreshTokensResponseObject, error) {
	command := idCommand.RefreshUserTokens{
		RefreshToken: request.Params.RefreshToken,
	}
	response, err := h.commandBus.Dispatch[idCommand.RefreshUserTokensResponse](ctx, command)
	if err != nil {
		return nil, err
	}

	return RefreshTokens200JSONResponseCustom{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}, nil
}

func (h *ApiHandler) Register(ctx context.Context, request RegisterRequestObject) (RegisterResponseObject, error) {
	clientIP := middleware.GetClientIP(ctx)
	ip, err := netip.ParseAddr(clientIP)
	if err != nil {
		return nil, err
	}

	command := idCommand.RegisterUser{
		Login:     request.Body.Login,
		Password:  request.Body.Password,
		IP:        ip,
		UserAgent: request.Params.UserAgent,
	}
	response, err := h.commandBus.Dispatch[idCommand.RegisterUserResponse](ctx, command)
	if err != nil {
		return nil, err
	}

	return Register200JSONResponseCustom{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}, nil
}

func getAuth(ctx context.Context) (auth.Auth, error) {
	a, err := auth.AuthFromContext(ctx)
	if err != nil {
		err = fmt.Errorf("%w: %w", authDomain.ErrUnauthorized, err)
		return auth.Auth{}, apperr.Map(err)
	}
	return a, nil
}

type Register200JSONResponseCustom struct {
	Register200JSONResponse
}

func (response Register200JSONResponseCustom) VisitRegisterResponse(w http.ResponseWriter) error {
	setAuthCookies(w, response.AccessToken, response.RefreshToken)
	return response.Register200JSONResponse.VisitRegisterResponse(w)
}

type Login200JSONResponseCustom struct {
	Login200JSONResponse
}

func (response Login200JSONResponseCustom) VisitLoginResponse(w http.ResponseWriter) error {
	setAuthCookies(w, response.AccessToken, response.RefreshToken)
	return response.Login200JSONResponse.VisitLoginResponse(w)
}

type RefreshTokens200JSONResponseCustom struct {
	RefreshTokens200JSONResponse
}

func (response RefreshTokens200JSONResponseCustom) VisitRefreshTokensResponse(w http.ResponseWriter) error {
	setAuthCookies(w, response.AccessToken, response.RefreshToken)
	return response.RefreshTokens200JSONResponse.VisitRefreshTokensResponse(w)
}

func setAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	accessTokenCookie := &http.Cookie{
		Name:     cookies.AccessToken,
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(token.AccessTokenTTL),
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	refreshTokenCookie := &http.Cookie{
		Name:     cookies.RefreshToken,
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(session.TTL),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)
}

func (h *ApiHandler) SendMsg(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	a, err := auth.AuthFromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to accept websocket connection: unauthorized", http.StatusUnauthorized)
		return
	}

	var rqst struct {
		Body        string    `json:"body"`
		RecipientID uuid.UUID `json:"recipietId"`
	}
	err = json.UnmarshalRead(r.Body, &rqst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cmd := cnvCommand.SendFirstDirectMessage{
		Sender:    a.Principal.Subject.ID,
		Recipient: rqst.RecipientID,
		Body:      rqst.Body,
	}

	rsp, err := h.commandBus.Dispatch[any](r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.MarshalWrite(w, rsp)
	// w.Write([]byte("foo"))
}
