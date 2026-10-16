package http

import (
	"context"
	"encoding/json/v2"
	"fmt"

	// cnvCommand "messenger/messenger/internal/conversations/app/command"
	"messenger/messenger/internal/directory/app/command"
	"messenger/messenger/internal/directory/app/query"
	idCommand "messenger/messenger/internal/identity/app/command"
	idQuery "messenger/messenger/internal/identity/app/query"
	"messenger/messenger/internal/identity/app/token"
	authDomain "messenger/messenger/internal/identity/domain/model"
	"messenger/messenger/internal/identity/domain/model/session"
	"messenger/messenger/internal/platform/auth"
	"messenger/messenger/internal/platform/commandbus"
	"messenger/messenger/internal/platform/querybus"
	"messenger/messenger/internal/platform/transport/http/cookies"
	"messenger/messenger/internal/shared/app/apperr"
	"net/http"
	"net/netip"
	"time"

	// "uuid"

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

	q := idQuery.GetAccount{
		ID: a.Principal.Subject.ID,
	}
	response, err := h.queryBus.Dispatch[idQuery.AccountView](ctx, q)
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

func (h *ApiHandler) LoginUser(ctx context.Context, request LoginUserRequestObject) (LoginUserResponseObject, error) {
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

func (h *ApiHandler) RegisterDirectoryEntry(ctx context.Context, request RegisterDirectoryEntryRequestObject) (RegisterDirectoryEntryResponseObject, error) {
	cmd := command.RegisterEntry{
		TargetID: request.Body.TargetId,
		TargetType: request.Body.TargetType,
		DisplayName: request.Body.DisplayName,
		Handle: request.Body.Handle,
	}

	_, err := h.commandBus.Dispatch[any](ctx, cmd)
	if err != nil {
		return nil, err
	}

	return RegisterDirectoryEntry201Response{}, nil
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

func (h *ApiHandler) RegisterUser(ctx context.Context, request RegisterUserRequestObject) (RegisterUserResponseObject, error) {
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
	RegisterUser200JSONResponse
}

func (response Register200JSONResponseCustom) VisitRegisterUserResponse(w http.ResponseWriter) error {
	setAuthCookies(w, response.AccessToken, response.RefreshToken)
	return response.RegisterUser200JSONResponse.VisitRegisterUserResponse(w)
}

type Login200JSONResponseCustom struct {
	LoginUser200JSONResponse
}

func (response Login200JSONResponseCustom) VisitLoginUserResponse(w http.ResponseWriter) error {
	setAuthCookies(w, response.AccessToken, response.RefreshToken)
	return response.LoginUser200JSONResponse.VisitLoginUserResponse(w)
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

func (h *ApiHandler) Test(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	_, err := auth.AuthFromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to accept websocket connection: unauthorized", http.StatusUnauthorized)
		return
	}

	// var rqst struct {
	// 	Body        string    `json:"body"`
	// 	RecipientID uuid.UUID `json:"recipietId"`
	// }

	// err = json.UnmarshalRead(r.Body, &rqst)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// cmd := cnvCommand.SendFirstDirectMessage{
	// 	Sender:    a.Principal.Subject.ID,
	// 	Recipient: rqst.RecipientID,
	// 	Body:      rqst.Body,
	// }

	// rsp, err := h.commandBus.Dispatch[any](r.Context(), cmd)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	var rqst struct {
		Search string `json:"search"`
	}

	err = json.UnmarshalRead(r.Body, &rqst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	qry := query.Search{
		Str: rqst.Search,
	}

	rsp, err := h.queryBus.Dispatch[any](r.Context(), qry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.MarshalWrite(w, rsp)
	// w.Write([]byte("foo"))
}
