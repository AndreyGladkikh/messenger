package testhelpers

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	TemplateDBName string
	MigrationsConnectionString string
	AdminConnectionString string
}

type RedisContainer struct {
	*redis.RedisContainer
	ConnectionString string
}

type RabbitMQContainer struct {
	*rabbitmq.RabbitMQContainer
	ConnectionString string
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	templateDBName := "test_template"

	pgContainer, err := postgres.Run(ctx,
		"postgres:18.4-alpine",
		postgres.WithDatabase(templateDBName),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	migrationsConnStr, err := pgContainer.ConnectionString(ctx, fmt.Sprintf("dbname=%s", templateDBName), "sslmode=disable")
	if err != nil {
		return nil, err
	}

	adminConnStr, err := pgContainer.ConnectionString(ctx, "dbname=postgres", "sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		TemplateDBName: templateDBName,
		MigrationsConnectionString: migrationsConnStr,
		AdminConnectionString: adminConnStr,
	}, nil
}

func (p *PostgresContainer) ConnectionString(dbName string) (string, error) {
	connStr, err := p.PostgresContainer.ConnectionString(context.Background(), fmt.Sprintf("dbname=%s", dbName), "sslmode=disable")
	if err != nil {
		return "", err
	}

	return connStr, nil
}

func CreateRedisContainer(ctx context.Context) (*RedisContainer, error) {
	redisContainer, err := redis.Run(ctx,
		"redis:8.10.1-alpine",
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelVerbose),
		// tcredis.WithConfigFile(filepath.Join("testdata", "redis7.conf")),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}

	return &RedisContainer{
		RedisContainer: redisContainer,
		ConnectionString:  connStr,
	}, nil
}

func CreateRabbitMQContainer(ctx context.Context) (*RabbitMQContainer, error) {
	rabbitmqContainer, err := rabbitmq.Run(ctx,
		"rabbitmq:4.3.6-alpine",
		rabbitmq.WithAdminUsername("rabbit"),
		rabbitmq.WithAdminPassword("rabbit"),
	)
	if err != nil {
		return nil, err
	}
	connStr, err := rabbitmqContainer.AmqpURL(ctx)
	if err != nil {
		return nil, err
	}

	return &RabbitMQContainer{
		RabbitMQContainer: rabbitmqContainer,
		ConnectionString:  connStr,
	}, nil
}
