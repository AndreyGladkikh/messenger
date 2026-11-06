package e2e

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"uuid"

	"messenger/messenger/internal/platform/bootstrap/api"
	"messenger/messenger/internal/platform/bootstrap/consumer"
	"messenger/messenger/internal/platform/bootstrap/outbox_relay"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/migrate"
	"messenger/messenger/internal/platform/postgres"
	appHttp "messenger/messenger/internal/platform/transport/http"

	"messenger/messenger/test/testhelpers"

	"net/http/httptest"

	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
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
	t.Helper()

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

	relayCleanup := startOutboxRelay(t, cfg)
	consumerCleanup := startConsumers(t, cfg)

	t.Cleanup(func() {
		consumerCleanup()
		relayCleanup()
		apiCleanup()

		_, err := pgAdmin.Exec(
			context.Background(),
			fmt.Sprintf("DROP DATABASE %s WITH (FORCE)", testDBName),
		)
		require.NoError(t, err)
	})

	server := httptest.NewTestServer(t, appApi.HttpHandler)
	client := newHttpClient(t, server)

	return server, client
}

func startOutboxRelay(t *testing.T, cfg config.Config) func() {
	t.Helper()

	outboxRelay, cleanup, err := outbox_relay.Initialize(cfg)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		outboxRelay.OutboxRelay.Run(ctx)
	}()

	return func() {
		cancel()
		cleanup()
	}
}

func startConsumers(t *testing.T, cfg config.Config) func() {
	t.Helper()

	consumer, cleanup, err := consumer.Initialize(cfg)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	for _, runner := range consumer.Consumers {
		go func() {
			runner.Run(ctx)
		}()
	}

	return func() {
		cancel()
		cleanup()
	}
}

func decode[T any](t *testing.T, data io.Reader) T {
	t.Helper()
	
	var decoded T
	err := json.UnmarshalRead(data, &decoded)
	require.NoError(t, err)
	return decoded
}
