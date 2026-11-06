package e2e

import (
	"time"

	appHttp "messenger/messenger/internal/platform/transport/http"
	"messenger/messenger/internal/platform/transport/http/cookies"
	"messenger/messenger/internal/shared/app/apperr"

	"net/http"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser_Success(t *testing.T) {
	_, client := initApi(t)

	resp, err := client.registerUser(t, "alice", "TestPass123!")
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	body := decode[appHttp.Tokens](t, resp.Body)
	assert.NotNil(t, body.AccessToken)
	assert.NotNil(t, body.RefreshToken)

	respCookies := resp.Cookies()
	assert.NotEmpty(t, respCookies)
	var accessTokenCookie, refreshTokenCookie *http.Cookie
	for _, c := range respCookies {
		switch c.Name {
		case cookies.AccessToken:
			accessTokenCookie = c
		case cookies.RefreshToken:
			refreshTokenCookie = c
		}
	}
	require.NotNil(t, accessTokenCookie)
	assert.NotEmpty(t, accessTokenCookie.Value)
	require.NotNil(t, refreshTokenCookie)
	assert.NotEmpty(t, refreshTokenCookie.Value)
}

func TestRegisterUser_LoginAlreadyExists(t *testing.T) {
	_, client := initApi(t)

	login := "alice"
	pass := "TestPass123!"

	resp, err := client.registerUser(t, login, pass)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, err = client.registerUser(t, login, pass)
	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	body := decode[appHttp.ProblemDetails](t, resp.Body)
	assert.Equal(t, problemTypesRegistry.AlreadyExists, *body.Type)
	assert.Equal(t, apperr.KindAlreadyExists.String(), body.Title)
}

func TestLoginUser_Success(t *testing.T) {
	_, client := initApi(t)

	login := "alice"
	pass := "TestPass123!"

	resp, err := client.registerUser(t, login, pass)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, err = client.loginUser(t, login, pass)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := decode[appHttp.Tokens](t, resp.Body)
	assert.NotNil(t, body.AccessToken)
	assert.NotNil(t, body.RefreshToken)

	respCookies := resp.Cookies()
	assert.NotEmpty(t, respCookies)
	var accessTokenCookie, refreshTokenCookie *http.Cookie
	for _, c := range respCookies {
		switch c.Name {
		case cookies.AccessToken:
			accessTokenCookie = c
		case cookies.RefreshToken:
			refreshTokenCookie = c
		}
	}
	require.NotNil(t, accessTokenCookie)
	assert.NotEmpty(t, accessTokenCookie.Value)
	require.NotNil(t, refreshTokenCookie)
	assert.NotEmpty(t, refreshTokenCookie.Value)
}

func TestLoginUser_NonExistentUser(t *testing.T) {
	_, client := initApi(t)

	resp, err := client.loginUser(t, "alice", "TestPass123!")
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	body := decode[appHttp.ProblemDetails](t, resp.Body)
	assert.Equal(t, problemTypesRegistry.Blank, *body.Type)
	assert.NotNil(t, apperr.KindUnauthenticated, body.Title)
}

func TestLoginUser_WrongPassword(t *testing.T) {
	_, client := initApi(t)

	resp, err := client.registerUser(t, "alice", "TestPass123!")
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, err = client.loginUser(t, "alice", "WrongPass123!")
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	body := decode[appHttp.ProblemDetails](t, resp.Body)
	assert.Equal(t, problemTypesRegistry.Blank, *body.Type)
	assert.NotNil(t, apperr.KindUnauthenticated, body.Title)
}

func TestRefreshUserTokens_Success(t *testing.T) {
	_, client := initApi(t)

	resp, err := client.registerUser(t, "alice", "TestPass123!")
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	body := decode[appHttp.Tokens](t, resp.Body)

	firstAccessToken := body.AccessToken
	firstRefreshToken := body.RefreshToken

	time.Sleep(time.Second)

	resp, err = client.refreshUserTokens(t)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body = decode[appHttp.Tokens](t, resp.Body)
	assert.NotEmpty(t, body.AccessToken)
	assert.NotEmpty(t, body.RefreshToken)

	assert.NotEqual(t, firstAccessToken, body.AccessToken)
	assert.NotEqual(t, firstRefreshToken, body.RefreshToken)

	respCookies := resp.Cookies()
	assert.NotEmpty(t, respCookies)
	var accessTokenCookie, refreshTokenCookie *http.Cookie
	for _, c := range respCookies {
		switch c.Name {
		case cookies.AccessToken:
			accessTokenCookie = c
		case cookies.RefreshToken:
			refreshTokenCookie = c
		}
	}
	require.NotNil(t, accessTokenCookie)
	assert.NotEmpty(t, accessTokenCookie.Value)
	require.NotNil(t, refreshTokenCookie)
	assert.NotEmpty(t, refreshTokenCookie.Value)
}

func TestMalformedRequest(t *testing.T) {
	_, client := initApi(t)

	resp, err := client.doRequest(t, http.MethodPost, "auth/users", &httpRequestOpts{Body: map[string]string{"wrongField": "alice"}})
	require.NoError(t, err)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	body := decode[appHttp.ProblemDetails](t, resp.Body)
	assert.Equal(t, problemTypesRegistry.ValidationError, *body.Type)
	assert.Equal(t, "Validation error", body.Title)
}
