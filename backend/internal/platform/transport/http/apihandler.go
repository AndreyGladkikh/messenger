package http

import (
	"context"
	"encoding/json/v2"

	webBffQry "messenger/messenger/internal/bff/web/app/query"
	directoryCmd "messenger/messenger/internal/directory/app/command"
	directoryQry "messenger/messenger/internal/directory/app/query"
	idApp "messenger/messenger/internal/identity/app"
	idCommand "messenger/messenger/internal/identity/app/command"
	idQuery "messenger/messenger/internal/identity/app/query"
	"messenger/messenger/internal/identity/app/token"
	"messenger/messenger/internal/identity/domain/model/session"
	cnvCmd "messenger/messenger/internal/messaging/app/command"
	"messenger/messenger/internal/platform/auth"
	"messenger/messenger/internal/platform/command_handler"
	"messenger/messenger/internal/platform/query_handler"
	"messenger/messenger/internal/platform/transport/http/cookies"
	"net/http"
	"net/netip"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

var _ StrictServerInterface = (*ApiHandler)(nil)

type ApiHandler struct {
	commands *command_handler.Registry
	queries  *query_handler.Registry
}

func NewApiHandler(
	commands *command_handler.Registry,
	queries *query_handler.Registry,
) *ApiHandler {
	return &ApiHandler{
		commands: commands,
		queries:  queries,
	}
}

func (h *ApiHandler) GetMe(ctx context.Context, request GetMeRequestObject) (GetMeResponseObject, error) {
	a, err := getAuth(ctx)
	if err != nil {
		return nil, err
	}

	qry := idQuery.GetAccount{
		ID: a.Principal.Subject.ID,
	}
	response, err := h.queries.GetAccount.Handle(ctx, qry)
	if err != nil {
		return nil, err
	}

	return GetMe200JSONResponse{
		Id:   response.ID.String(),
		Type: response.Type,
	}, nil
}

func (h *ApiHandler) LoginUser(ctx context.Context, request LoginUserRequestObject) (LoginUserResponseObject, error) {
	clientIP := middleware.GetClientIP(ctx)
	ip, err := netip.ParseAddr(clientIP)
	if err != nil {
		return nil, err
	}

	cmd := idCommand.LoginUser{
		Login:     request.Body.Login,
		Password:  request.Body.Password,
		IP:        ip,
		UserAgent: request.Params.UserAgent,
	}
	response, err := h.commands.LoginUser.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return Login200JSONResponseCustom{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}, nil
}

func (h *ApiHandler) RegisterDirectoryEntry(ctx context.Context, request RegisterDirectoryEntryRequestObject) (RegisterDirectoryEntryResponseObject, error) {
	cmd := directoryCmd.RegisterEntry{
		TargetID:    request.Body.TargetId,
		TargetType:  request.Body.TargetType,
		DisplayName: request.Body.DisplayName,
		Handle:      request.Body.Handle,
	}

	_, err := h.commands.RegisterDirectoryEntry.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return RegisterDirectoryEntry201Response{}, nil
}

func (h *ApiHandler) RefreshTokens(ctx context.Context, request RefreshTokensRequestObject) (RefreshTokensResponseObject, error) {
	cmd := idCommand.RefreshUserTokens{
		RefreshToken: request.Params.RefreshToken,
	}
	response, err := h.commands.RefreshUserTokens.Handle(ctx, cmd)
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

	cmd := idCommand.RegisterUser{
		Login:     request.Body.Login,
		Password:  request.Body.Password,
		IP:        ip,
		UserAgent: *request.Params.UserAgent,
	}
	response, err := h.commands.RegisterUser.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return Register200JSONResponseCustom{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}, nil
}

func (h *ApiHandler) StartDirectConversation(ctx context.Context, request StartDirectConversationRequestObject) (StartDirectConversationResponseObject, error) {
	a, err := getAuth(ctx)
	if err != nil {
		return nil, err
	}

	cmd := cnvCmd.StartDirectConversation{
		Sender:    a.Principal.Subject.ID,
		Recipient: request.Body.ParticipantId,
		Body:      request.Body.FirstMessage.Body,
	}
	resp, err := h.commands.StartDirectConversation.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return StartDirectConversation201JSONResponse{
		ConversationId: resp.ConversationID,
		MessageId:      resp.MessageID,
		MessageSentAt:  resp.MessageSentAt,
	}, nil
}

func (h *ApiHandler) SendMessage(ctx context.Context, request SendMessageRequestObject) (SendMessageResponseObject, error) {
	a, err := getAuth(ctx)
	if err != nil {
		return nil, err
	}

	cmd := cnvCmd.SendMessage{
		SenderID:       a.Principal.Subject.ID,
		ConversationID: request.ConversationId,
		Body:           request.Body.Body,
	}
	_, err = h.commands.SendMessage.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return SendMessage201Response{}, nil
}

func (h *ApiHandler) GetWebHomePage(ctx context.Context, request GetWebHomePageRequestObject) (GetWebHomePageResponseObject, error) {
	a, err := getAuth(ctx)
	if err != nil {
		return nil, err
	}

	qry := webBffQry.GetHomePage{
		UserID: a.Principal.Subject.ID,
	}
	resp, err := h.queries.GetHomePage.Handle(ctx, qry)
	if err != nil {
		return nil, err
	}

	peerItem := make([]PeerItem, 0, len(resp.Peers))
	for _, p := range resp.Peers {
		peerItem = append(peerItem, PeerItem{
			Id:          p.ID,
			Type:        p.Type,
			AvatarUrl:   p.AvatarUrl,
			Title:       p.Title,
			Preview:     &p.Preview,
			PreviewTime: &p.PreviewTime,
			UnreadCount: &p.UnreadCount,
		})
	}

	return GetWebHomePage200JSONResponse{
		Peers: peerItem,
	}, nil
}

func (h *ApiHandler) CreateMyProfile(ctx context.Context, request CreateMyProfileRequestObject) (CreateMyProfileResponseObject, error) {
	a, err := getAuth(ctx)
	if err != nil {
		return nil, err
	}

	cmd := directoryCmd.RegisterEntry{
		TargetID:    a.Principal.Subject.ID,
		TargetType:  string(a.Principal.Subject.Type),
		DisplayName: request.Body.DisplayName,
		Handle:      request.Body.Handle,
		AvatarURL:   request.Body.AvatarUrl,
	}
	_, err = h.commands.RegisterDirectoryEntry.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return CreateMyProfile201Response{}, nil
}

func (h *ApiHandler) GetMyProfile(ctx context.Context, request GetMyProfileRequestObject) (GetMyProfileResponseObject, error) {
	a, err := getAuth(ctx)
	if err != nil {
		return nil, err
	}

	qry := directoryQry.GetEntry{
		ID:   a.Principal.Subject.ID,
		Type: string(a.Principal.Subject.Type),
	}
	resp, err := h.queries.GetDirectoryEntry.Handle(ctx, qry)
	if err != nil {
		return nil, err
	}

	return GetMyProfile200JSONResponse{
		Handle:      resp.Handle,
		DisplayName: resp.DisplayName,
		AvatarUrl:   resp.AvatarURL,
	}, nil
}

func getAuth(ctx context.Context) (auth.Auth, error) {
	a, err := auth.AuthFromContext(ctx)
	if err != nil {
		return auth.Auth{}, idApp.ErrNeedLogin
	}
	return a, nil
}

type Register200JSONResponseCustom struct {
	RegisterUser201JSONResponse
}

func (response Register200JSONResponseCustom) VisitRegisterUserResponse(w http.ResponseWriter) error {
	setAuthCookies(w, response.AccessToken, response.RefreshToken)
	return response.RegisterUser201JSONResponse.VisitRegisterUserResponse(w)
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

// todo rm
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

	qry := directoryQry.Search{
		Str: rqst.Search,
	}

	rsp, err := h.queries.SearchDirectoryEntry.Handle(r.Context(), qry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.MarshalWrite(w, rsp)
	// w.Write([]byte("foo"))
}
