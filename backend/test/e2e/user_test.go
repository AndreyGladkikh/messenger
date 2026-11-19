package e2e

import (
	appHttp "messenger/messenger/internal/platform/transport/http"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"uuid"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

type user struct {
	id         uuid.UUID
	login      string
	password   string
	profile    *profile
	httpClient *apiHTTPClient
	wsClient   *wsClient
}

type profile struct {
	handle        *string
	displayName   string
	avatarMediaID *uuid.UUID
	avatarURL     *string
	bio           *string
	birthday      *time.Time
}

func newUser(t *testing.T, server *httptest.Server, login string, opts ...userOption) *user {
	t.Helper()

	httpClient := newHttpClient(t, server)

	password := gofakeit.Password(true, true, true, true, false, 8)

	resp, _ := httpClient.registerUser(t, login, password)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	registerBody := decode[appHttp.Tokens](t, resp.Body)

	httpClient.JWT = &registerBody.AccessToken

	resp, _ = httpClient.me(t)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	meBody := decode[appHttp.GetMe200JSONResponse](t, resp.Body)
	id, err := uuid.Parse(meBody.Id)
	require.NoError(t, err)

	wsClient := newWSClient(t, httpClient)
	t.Cleanup(func() {
		wsClient.CloseNow()
	})

	u := &user{
		id:         id,
		password:   password,
		login:      login,
		httpClient: httpClient,
		wsClient:   wsClient,
	}

	for _, o := range opts {
		o(u)
	}

	return u
}

func alice(t *testing.T, server *httptest.Server) *user {
	handle := "alice_handle"
	return newUser(t, server, "alice", userWithProfile(t, &profileOptions{
		handle:      &handle,
		displayName: "Alice DN",
	}))
}

func bob(t *testing.T, server *httptest.Server) *user {
	handle := "bob_handle"
	return newUser(t, server, "bob", userWithProfile(t, &profileOptions{
		handle:      &handle,
		displayName: "Bob DN",
	}))
}

type userOption func(*user)

func userWithProfile(t *testing.T, opts *profileOptions) userOption {
	return func(u *user) {
		u.profile = u.setProfile(t, opts)
	}
}

type profileOptions struct {
	withoutHandle   bool
	withoutBio      bool
	withoutBirthday bool
	withoutAvatar   bool

	handle      *string
	displayName string
	bio         *string
	birthday    *time.Time
	avatar      []byte
}

func (u *user) setProfile(t *testing.T, opts *profileOptions) *profile {
	t.Helper()

	require.NotEmpty(t, u)

	if opts == nil {
		opts = new(profileOptions)
	}
	if opts.displayName == "" {
		opts.displayName = gofakeit.Name()
	}
	if opts.bio == nil && !opts.withoutBio {
		bio := gofakeit.Bio()
		opts.bio = &bio
	}
	if opts.birthday == nil && !opts.withoutBirthday {
		bd := gofakeit.Date()
		opts.birthday = &bd
	}
	if opts.handle == nil && !opts.withoutHandle {
		handle := gofakeit.Noun()
		opts.handle = &handle
	}
	if opts.avatar == nil && !opts.withoutAvatar {
		avatar := gofakeit.ImagePng(500, 500)
		opts.avatar = avatar
	}

	var avatarMediaID *uuid.UUID
	var avatarURL *string
	if !opts.withoutAvatar {
		avatar := uploadMedia(t, u, gofakeit.Noun(), opts.avatar, "avatar")
		avatarMediaID, avatarURL = &avatar.id, avatar.publicURL
	}

	resp, _ := u.httpClient.createMyProfile(t, appHttp.CreateMyProfileJSONRequestBody{
		Handle:        opts.handle,
		DisplayName:   opts.displayName,
		Bio:           opts.bio,
		Birthday:      dateFromTime(opts.birthday),
		AvatarMediaId: avatarMediaID,
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	u.profile = &profile{
		displayName:   opts.displayName,
		handle:        opts.handle,
		avatarMediaID: avatarMediaID,
		avatarURL:     avatarURL,
		bio:           opts.bio,
		birthday:      opts.birthday,
	}

	return u.profile
}
