package e2e

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	appHttp "messenger/messenger/internal/platform/transport/http"
	"messenger/messenger/internal/platform/transport/http/contenttype"
	"messenger/messenger/internal/platform/transport/http/headers"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/publicsuffix"
)

type httpClient struct {
	*http.Client
	server *httptest.Server
	JWT    *string
}

func newHttpClient(t *testing.T, server *httptest.Server) *httpClient {
	t.Helper()

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	require.NoError(t, err)

	client := *server.Client()

	client.Jar = jar

	return &httpClient{
		Client: &client,
		server: server,
	}
}

type httpRequestOpts struct {
	Body        any
	Timeout     time.Duration
	ContentType string
	Cookies     []*http.Cookie
}

func (c *httpClient) doRequest(t *testing.T, method, url string, opts *httpRequestOpts) (*http.Response, error) {
	if opts == nil {
		opts = new(httpRequestOpts)
	}

	if opts.Timeout == 0 {
		// todo remove multiplier
		opts.Timeout = time.Second * 100000
	}

	if opts.ContentType == "" {
		opts.ContentType = contenttype.ApplicationJson
	}

	body := new(bytes.Buffer)
	if opts.Body != nil {
		switch opts.ContentType {
		case contenttype.ApplicationJson:
			err := json.MarshalWrite(body, opts.Body)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("no encoder for content type")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("https://www.example.com/%s", url), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set(headers.ContentType, opts.ContentType)

	if c.JWT != nil {
		req.Header.Set(headers.Authorization, fmt.Sprintf("Bearer %s", *c.JWT))
	}

	for _, c := range opts.Cookies {
		req.AddCookie(c)
	}

	rsp, err := c.Client.Do(req)
	require.NoError(t, err)

	t.Cleanup(func() {
		rsp.Body.Close()
	})

	return rsp, nil
}

func (c *httpClient) registerUser(t *testing.T, login, pass string) (*http.Response, error) {
	return c.doRequest(t, http.MethodPost, "auth/users", &httpRequestOpts{Body: appHttp.RegisterUserJSONRequestBody{
		Login:    login,
		Password: pass,
	}})
}

func (c *httpClient) loginUser(t *testing.T, login, pass string) (*http.Response, error) {
	return c.doRequest(t, http.MethodPost, "auth/users/login", &httpRequestOpts{Body: appHttp.LoginUserJSONRequestBody{
		Login:    login,
		Password: pass,
	}})
}

func (c *httpClient) logoutUser(t *testing.T) (*http.Response, error) {
	return c.doRequest(t, http.MethodPost, "auth/users/logout", nil)
}

func (c *httpClient) refreshUserTokens(t *testing.T) (*http.Response, error) {
	return c.doRequest(t, http.MethodPost, "auth/users/refresh-tokens", nil)
}

func (c *httpClient) me(t *testing.T) (*http.Response, error) {
	return c.doRequest(t, http.MethodGet, "auth/me", nil)
}

func (c *httpClient) createMyProfile(t *testing.T, r appHttp.CreateMyProfileJSONRequestBody) (*http.Response, error) {
	return c.doRequest(t, http.MethodPost, "me/profile", &httpRequestOpts{Body: r})
}

func (c *httpClient) getMyProfile(t *testing.T) (*http.Response, error) {
	return c.doRequest(t, http.MethodGet, "me/profile", nil)
}

func (c *httpClient) startDirectConversation(t *testing.T, r appHttp.StartDirectConversationJSONRequestBody) (*http.Response, error) {
	return c.doRequest(t, http.MethodPost, "conversations/direct", &httpRequestOpts{Body: r})
}

func (c *httpClient) registerDirectoryEntry(t *testing.T, r appHttp.RegisterDirectoryEntryJSONRequestBody) (*http.Response, error) {
	return c.doRequest(t, http.MethodPost, "directory/entries", &httpRequestOpts{Body: r})
}

func (c *httpClient) getHomePage(t *testing.T) (*http.Response, error) {
	return c.doRequest(t, http.MethodGet, "bff/web/home-page", nil)
}
