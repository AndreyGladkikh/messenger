package http

import (
	"context"
	"encoding/json/v2"
	webBffQry "messenger/messenger/internal/bff/web/application/query"
	drctrQry "messenger/messenger/internal/directory/app/query"
	idQry "messenger/messenger/internal/identity/app/query"
	"messenger/messenger/internal/platform/auth"
	appQuery "messenger/messenger/internal/shared/app/query"
	"net/http/httptest"
	"strings"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApiHandlerGetMeUsesAuthenticatedAccount(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	getAccount := appQuery.HandlerFunc[idQry.GetAccount, idQry.AccountView](
		func(_ context.Context, qry idQry.GetAccount) (idQry.AccountView, error) {
			assert.Equal(t, userID, qry.ID)
			return idQry.AccountView{ID: userID, Type: "user"}, nil
		},
	)
	handler := &ApiHandler{getAccount: getAccount}
	ctx := auth.NewContextWithAuth(t.Context(), auth.Auth{
		Principal: auth.Principal{Subject: auth.Subject{ID: userID, Type: auth.SubjectTypeUser}},
	})

	response, err := handler.GetMe(ctx, GetMeRequestObject{})

	require.NoError(t, err)
	assert.Equal(t, GetMe200JSONResponse{Id: userID.String(), Type: "user"}, response)
}

func TestApiHandlerGetWebHomePageUsesAuthenticatedUser(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	getHomePage := appQuery.HandlerFunc[webBffQry.GetHomePage, webBffQry.HomePage](
		func(_ context.Context, qry webBffQry.GetHomePage) (webBffQry.HomePage, error) {
			assert.Equal(t, userID, qry.UserID)
			return webBffQry.HomePage{Peers: []webBffQry.PeerItem{}}, nil
		},
	)
	handler := &ApiHandler{getWebHomePage: getHomePage}
	ctx := auth.NewContextWithAuth(t.Context(), auth.Auth{
		Principal: auth.Principal{Subject: auth.Subject{ID: userID, Type: auth.SubjectTypeUser}},
	})

	response, err := handler.GetWebHomePage(ctx, GetWebHomePageRequestObject{})

	require.NoError(t, err)
	assert.Equal(t, GetWebHomePage200JSONResponse{Peers: []PeerItem{}}, response)
}

func TestApiHandlerTestSearchUsesTypedHandler(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	expected := []drctrQry.SearchResultItem{{
		ID:          uuid.New(),
		Type:        "user",
		Handle:      "alice",
		DisplayName: "Alice",
	}}
	searchDirectory := appQuery.HandlerFunc[drctrQry.Search, []drctrQry.SearchResultItem](
		func(_ context.Context, qry drctrQry.Search) ([]drctrQry.SearchResultItem, error) {
			assert.Equal(t, "alice", qry.Str)
			return expected, nil
		},
	)
	handler := &ApiHandler{searchDirectory: searchDirectory}
	request := httptest.NewRequest("POST", "/test", strings.NewReader(`{"search":"alice"}`))
	request = request.WithContext(auth.NewContextWithAuth(request.Context(), auth.Auth{
		Principal: auth.Principal{Subject: auth.Subject{ID: userID, Type: auth.SubjectTypeUser}},
	}))
	recorder := httptest.NewRecorder()

	handler.Test(recorder, request)

	assert.Equal(t, 200, recorder.Code)
	var response []drctrQry.SearchResultItem
	require.NoError(t, json.UnmarshalRead(recorder.Body, &response))
	assert.Equal(t, expected, response)
}
