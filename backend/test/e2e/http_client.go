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
	"time"
)

type httpClient struct {
	*http.Client
	JWT *string
}

type httpRequestOpts struct {
	Body        any
	Timeout     time.Duration
	ContentType string
	Cookies     []*http.Cookie
}

func (c *httpClient) doRequest(method, url string, opts *httpRequestOpts) (*http.Response, error) {
	if opts == nil {
		opts = new(httpRequestOpts)
	}
	
	if opts.Timeout == 0 {
		// todo remove 10000
		opts.Timeout = time.Second * 10000
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

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("http://www.example.com/%s", url), body)
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
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

// func (c *httpClient) doRequest(method, url string, opts *httpRequestOpts) (*http.Response, error) {
// 	if opts == nil {
// 		opts = new(httpRequestOpts)
// 	}

// 	// todo remove
// 	timeout := time.Second * 10000
// 	// timeout := time.Second
// 	if opts.Timeout != nil {
// 		timeout = *opts.Timeout
// 	}

// 	contType := contenttype.ApplicationJson
// 	if opts.ContentType != nil {
// 		contType = *opts.ContentType
// 	}

// 	body := new(bytes.Buffer)
// 	if opts.Body != nil {
// 		if contType == contenttype.ApplicationJson {
// 			err := json.MarshalWrite(body, opts.Body)
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 	defer cancel()

// 	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("http://www.example.com/%s", url), body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set(headers.ContentType, contType)

// 	if c.JWT != nil {
// 		req.Header.Set(headers.Authorization, fmt.Sprintf("Bearer %s", *c.JWT))
// 	}

// 	for _, c := range opts.Cookies {
// 		req.AddCookie(c)
// 	}

// 	rsp, err := c.Client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return rsp, nil
// }

func (c *httpClient) registerUser(login, pass string) (*http.Response, error) {
	return c.doRequest(http.MethodPost, "auth/users", &httpRequestOpts{Body: appHttp.RegisterUserJSONRequestBody{
		Login:    login,
		Password: pass,
	}})
}

func (c *httpClient) loginUser(login, pass string) (*http.Response, error) {
	return c.doRequest(http.MethodPost, "auth/users/login", &httpRequestOpts{Body: appHttp.LoginUserJSONRequestBody{
		Login:    login,
		Password: pass,
	}})
}

func (c *httpClient) refreshUserTokens() (*http.Response, error) {
	return c.doRequest(http.MethodPost, "auth/users/refreshTokens", nil)
}
