package e2e

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"log"
	"os"
	"uuid"

	"messenger/messenger/internal/app/api"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/migrate"
	"messenger/messenger/internal/platform/postgres"
	pltfrmHttp "messenger/messenger/internal/platform/transport/http"
	"messenger/messenger/internal/platform/transport/http/cookies"
	"messenger/messenger/test/testhelpers"

	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"

	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"

	"golang.org/x/net/publicsuffix"
)

var (
	pgContainer    *testhelpers.PostgresContainer
	redisContainer *testhelpers.RedisContainer
	rmqContainer   *testhelpers.RabbitMQContainer
	pgAdmin        *pgxpool.Pool
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

	migrator, migratorCleanup, err := migrate.New(config.Migrations{DBDSN: pgContainer.MigrationsConnectionString, Path: "../../db/migrations"})
	if err != nil {
		return nil, err
	}
	if err := migrator.Up(); err != nil {
		return nil, err
	}
	migratorCleanup()

	var pgAdminCleanup func()
	pgAdmin, pgAdminCleanup, err = postgres.NewPool(config.Database{Dsn: pgContainer.AdminConnectionString})
	if err != nil {
		return nil, err
	}

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
	testDBName := pgx.Identifier{fmt.Sprintf("test_%s", uuid.New().String())}.Sanitize()
	_, err := pgAdmin.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s TEMPLATE %s", testDBName, pgContainer.TemplateDBName))
	assert.NoError(t, err)

	pgConnString, err := pgContainer.ConnectionString(testDBName)
	assert.NoError(t, err)

	cfg := config.Config{
		Env: config.EnvLocal,
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
		Cache: config.Cache{Dsn: redisContainer.ConnectionString},
		Mq: config.Mq{Dsn: rmqContainer.ConnectionString},
	}

	appApi, cleanup, err := api.Initialize(cfg)
	assert.NoError(t, err)
	t.Cleanup(cleanup)

	server := httptest.NewTestServer(t, appApi.HttpHandler)

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	assert.NoError(t, err)

	server.Client().Jar = jar

	return server, &httpClient{
		Client: server.Client(),
	}
}

func TestRegisterUserSuccess(t *testing.T) {
	_, client := initApi(t)

	resp, err := client.registerUser("alice", "foo")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, 200)

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
	assert.NotNil(t, accessTokenCookie)
	assert.NotEmpty(t, accessTokenCookie.Value)
	assert.NotNil(t, refreshTokenCookie)
	assert.NotEmpty(t, refreshTokenCookie.Value)

	var data pltfrmHttp.Tokens
	err = json.UnmarshalRead(resp.Body, &data)
	assert.NoError(t, err)
	assert.NotNil(t, data.AccessToken)
	assert.NotNil(t, data.RefreshToken)
}
