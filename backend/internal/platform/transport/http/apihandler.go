package http

import (
	"context"
	"encoding/json/v2"
	"uuid"

	webBffQry "messenger/messenger/internal/bff/web/app/query"
	directoryCmd "messenger/messenger/internal/directory/app/command"
	directoryQry "messenger/messenger/internal/directory/app/query"
	idApp "messenger/messenger/internal/identity/app"
	idCommand "messenger/messenger/internal/identity/app/command"
	idQuery "messenger/messenger/internal/identity/app/query"
	"messenger/messenger/internal/identity/app/token"
	"messenger/messenger/internal/identity/domain/model/session"
	mediaCmd "messenger/messenger/internal/media/app/command"
	mediaQry "messenger/messenger/internal/media/app/query"
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

func (h *ApiHandler) LogoutUser(ctx context.Context, request LogoutUserRequestObject) (LogoutUserResponseObject, error) {
	cmd := idCommand.LogoutUser{
		RefreshToken: request.Params.RefreshToken,
	}
	_, err := h.commands.LogoutUser.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return LogoutUser204ResponseCustom{}, nil
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
			AvatarUrl:   p.AvatarURL,
			Title:       p.Title,
			Preview:     p.Preview,
			PreviewTime: p.PreviewTime,
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

	cmd := directoryCmd.SetUserProfile{
		UserID:      a.Principal.Subject.ID,
		Handle:      request.Body.Handle,
		DisplayName: request.Body.DisplayName,
		AvatarID:    request.Body.AvatarMediaId,
		Bio:         request.Body.Bio,
		Birthday:    request.Body.Birthday,
	}
	_, err = h.commands.SetUserProfile.Handle(ctx, cmd)
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

	qry := directoryQry.GetUserProfiles{
		UserIDs: []uuid.UUID{a.Principal.Subject.ID},
	}
	resp, err := h.queries.GetUserProfiles.Handle(ctx, qry)
	if err != nil {
		return nil, err
	}

	profile := resp[0]

	return GetMyProfile200JSONResponse{
		Handle:        profile.Handle,
		DisplayName:   profile.DisplayName,
		AvatarMediaId: profile.AvatarID,
	}, nil
}

func (h *ApiHandler) InitMediaUpload(ctx context.Context, request InitMediaUploadRequestObject) (InitMediaUploadResponseObject, error) {
	a, err := getAuth(ctx)
	if err != nil {
		return nil, err
	}

	cmd := mediaCmd.InitUpload{
		InitiatorID: a.Principal.Subject.ID,
		FileName:    request.Body.FileName,
		Purpose:     request.Body.Purpose,
	}
	resp, err := h.commands.InitMediaUpload.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return InitMediaUpload201JSONResponse{
		Id:           resp.ID,
		PresignedUrl: resp.PresignedURL,
	}, nil
}

func (h *ApiHandler) CompleteMediaUpload(ctx context.Context, request CompleteMediaUploadRequestObject) (CompleteMediaUploadResponseObject, error) {
	_, err := getAuth(ctx)
	if err != nil {
		return nil, err
	}

	cmd := mediaCmd.ComlpeteUpload{
		UploadID: request.UploadId,
	}
	mediaID, err := h.commands.CompleteMediaUpload.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return CompleteMediaUpload200JSONResponse{
		MediaId: mediaID,
	}, nil
}

func (h *ApiHandler) GetDownloadMediaUrl(ctx context.Context, request GetDownloadMediaUrlRequestObject) (GetDownloadMediaUrlResponseObject, error) {
	_, err := getAuth(ctx)
	if err != nil {
		return nil, err
	}

	qry := mediaQry.GetDownloadURL{
		MediaID: request.MediaId,
	}
	url, err := h.queries.GetDownloadMediaURL.Handle(ctx, qry)
	if err != nil {
		return nil, err
	}

	return GetDownloadMediaUrl200JSONResponse{
		Url: url,
	}, nil
}

func (h *ApiHandler) GlobalSearch(ctx context.Context, request GlobalSearchRequestObject) (GlobalSearchResponseObject, error) {
	_, err := getAuth(ctx)
	if err != nil {
		return nil, err
	}

	qry := webBffQry.Search{Str: request.Params.String}
	resp, err := h.queries.Search.Handle(ctx, qry)
	if err != nil {
		return nil, err
	}

	globalPeers := make([]PeerItem, 0, len(resp.Global))
	for _, r := range resp.Global {
		globalPeers = append(globalPeers, PeerItem{
			Id:        r.ID,
			Type:      r.Type,
			Title:     r.Title,
			AvatarUrl: r.AvatarURL,
			Handle:    r.Handle,
		})
	}
	return GlobalSearch200JSONResponse{
		Global: &globalPeers,
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

type LogoutUser204ResponseCustom struct {
	LogoutUser204Response
}

func (response LogoutUser204ResponseCustom) VisitLogoutUserResponse(w http.ResponseWriter) error {
	unsetAuthCookies(w)
	return response.LogoutUser204Response.VisitLogoutUserResponse(w)
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

func unsetAuthCookies(w http.ResponseWriter) {
	accessTokenCookie := &http.Cookie{
		Name:     cookies.AccessToken,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	refreshTokenCookie := &http.Cookie{
		Name:     cookies.RefreshToken,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
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
