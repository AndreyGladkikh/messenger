package e2e

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"uuid"

	"messenger/messenger/internal/app/api"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/migrate"
	"messenger/messenger/internal/platform/postgres"
	appHttp "messenger/messenger/internal/platform/transport/http"
	"messenger/messenger/internal/platform/transport/http/cookies"
	"messenger/messenger/internal/shared/app/apperr"
	"messenger/messenger/test/testhelpers"

	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"

	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	"golang.org/x/net/publicsuffix"
)

var (
	pgContainer          *testhelpers.PostgresContainer
	redisContainer       *testhelpers.RedisContainer
	rmqContainer         *testhelpers.RabbitMQContainer
	pgAdmin              *pgxpool.Pool
	problemTypesRegistry appHttp.ProblemTypesRegistry
	problemsRegistryUrl  = "https://problems-registry.messenger.com"
)

func TestMain(m *testing.M) {
	cleanup, err := setup()
	if err != nil {
		log.Fatalf("failed to setup tests: %v", err)
	}
	defer cleanup()

	code := m.Run()
	os.Exit(code)
}

func setup() (func(), error) {
	ctx := context.Background()

	var err error

	pgContainer, err = testhelpers.CreatePostgresContainer(ctx)
	if err != nil {
		return nil, err
	}
	redisContainer, err = testhelpers.CreateRedisContainer(ctx)
	if err != nil {
		return nil, err
	}
	rmqContainer, err = testhelpers.CreateRabbitMQContainer(ctx)
	if err != nil {
		return nil, err
	}

	migrator, migratorCleanup, err := migrate.New(config.Config{Migrations: config.Migrations{DBDSN: pgContainer.MigrationsConnectionString, Path: "../../db/migrations"}})
	if err != nil {
		return nil, err
	}
	if err := migrator.Up(); err != nil {
		return nil, err
	}
	migratorCleanup()

	var pgAdminCleanup func()
	pgAdmin, pgAdminCleanup, err = postgres.NewPool(config.Config{Database: config.Database{Dsn: pgContainer.AdminConnectionString}})
	if err != nil {
		return nil, err
	}

	problemTypesRegistry = appHttp.NewProblemTypesRegistry(config.Config{ProblemsRegistryUrl: problemsRegistryUrl})

	return func() {
		pgAdminCleanup()
		if err := testcontainers.TerminateContainer(pgContainer); err != nil {
			log.Fatalf("failed to terminate pg container: %s", err)
		}
		if err := testcontainers.TerminateContainer(redisContainer); err != nil {
			log.Fatalf("failed to terminate redis container: %s", err)
		}
		if err := testcontainers.TerminateContainer(rmqContainer); err != nil {
			log.Fatalf("failed to terminate rmq container: %s", err)
		}
	}, nil
}

func initApi(t *testing.T) (*httptest.Server, *httpClient) {
	testDBName := fmt.Sprintf("test_%s", strings.ReplaceAll(uuid.New().String(), "-", ""))

	_, err := pgAdmin.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s TEMPLATE %s", testDBName, pgContainer.TemplateDBName))
	require.NoError(t, err)

	pgConnString, err := pgContainer.ConnectionString(testDBName)
	require.NoError(t, err)

	cfg := config.Config{
		Env:                 config.EnvLocal,
		ProblemsRegistryUrl: problemsRegistryUrl,
		Auth: config.Auth{
			Jwt: struct {
				PrivateKeyPath string
				PublicKeyPath  string
			}{
				PrivateKeyPath: "../../keys/private.pem",
				PublicKeyPath:  "../../keys/public.pem",
			},
		},
		Database: config.Database{Dsn: pgConnString},
		Cache:    config.Cache{Dsn: redisContainer.ConnectionString},
		Mq:       config.Mq{Dsn: rmqContainer.ConnectionString},
	}

	appApi, apiCleanup, err := api.Initialize(cfg)
	require.NoError(t, err)

	t.Cleanup(func() {
		apiCleanup()

		_, err := pgAdmin.Exec(
			context.Background(),
			fmt.Sprintf("DROP DATABASE %s WITH (FORCE)", testDBName),
		)
		require.NoError(t, err)
	})

	server := httptest.NewTestServer(t, appApi.HttpHandler)

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	require.NoError(t, err)

	server.Client().Jar = jar

	return server, &httpClient{
		Client: server.Client(),
	}
}

func TestRegisterUser_Success(t *testing.T) {
	_, client := initApi(t)

	resp, err := client.registerUser("alice", "TestPass123!")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

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

	body := decode[appHttp.Tokens](t, resp.Body)
	assert.NotNil(t, body.AccessToken)
	assert.NotNil(t, body.RefreshToken)
}

func TestRegisterUser_LoginAlreadyExists(t *testing.T) {
	_, client := initApi(t)

	login := "alice"
	pass := "TestPass123!"

	resp, err := client.registerUser(login, pass)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = client.registerUser(login, pass)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	body := decode[appHttp.ProblemDetails](t, resp.Body)
	assert.Equal(t, problemTypesRegistry.AlreadyExists, *body.Type)
	assert.Equal(t, apperr.KindAlreadyExists.String(), body.Title)
}

func TestLoginUser_Success(t *testing.T) {
	_, client := initApi(t)

	login := "alice"
	pass := "TestPass123!"

	resp, err := client.registerUser(login, pass)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = client.loginUser(login, pass)
	require.NoError(t, err)
	defer resp.Body.Close()

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

	resp, err := client.loginUser("alice", "TestPass123!")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	body := decode[appHttp.ProblemDetails](t, resp.Body)
	assert.Equal(t, problemTypesRegistry.Blank, *body.Type)
	assert.NotNil(t, apperr.KindUnauthenticated, body.Title)
}

func TestLoginUser_WrongPassword(t *testing.T) {
	_, client := initApi(t)

	resp, err := client.registerUser("alice", "TestPass123!")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = client.loginUser("alice", "WrongPass123!")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	body := decode[appHttp.ProblemDetails](t, resp.Body)
	assert.Equal(t, problemTypesRegistry.Blank, *body.Type)
	assert.NotNil(t, apperr.KindUnauthenticated, body.Title)
}

func TestRefreshUserTokens_Success(t *testing.T) {
	_, client := initApi(t)

	resp, err := client.registerUser("alice", "TestPass123!")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	body := decode[appHttp.Tokens](t, resp.Body)

	firstAccessToken := body.AccessToken
	firstRefreshToken := body.RefreshToken

	time.Sleep(time.Second)

	resp, err = client.refreshUserTokens()
	require.NoError(t, err)
	defer resp.Body.Close()

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

	resp, err := client.doRequest(http.MethodPost, "auth/users", &httpRequestOpts{Body: map[string]string{"wrongField": "alice"}})
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	body := decode[appHttp.ProblemDetails](t, resp.Body)
	assert.Equal(t, problemTypesRegistry.ValidationError, *body.Type)
	assert.Equal(t, "Validation error", body.Title)
}

// func validateErrorResponse(t *testing.T, resp *http.Response) {

// }

func decode[T any](t *testing.T, data io.Reader) T {
	var decoded T
	err := json.UnmarshalRead(data, &decoded)
	require.NoError(t, err)
	return decoded
}
