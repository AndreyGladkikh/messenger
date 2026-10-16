package e2e

import (
	"context"
	"messenger/messenger/internal/app/api"
	"path/filepath"
	"testing"
	"time"

	// "github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestRegisterUser(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:18.4-alpine",
		postgres.WithInitScripts(filepath.Join("..", "..", "db", "migrations", "000001_init.up.sql")),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	// connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	// assert.NoError(t, err)

	// customerRepo, err := NewRepository(ctx, connStr)
	// assert.NoError(t, err)

	// c, err := customerRepo.CreateCustomer(ctx, Customer{
	// 	Name:  "Henry",
	// 	Email: "henry@gmail.com",
	// })
	// assert.NoError(t, err)
	// assert.NotNil(t, c)

	// customer, err := customerRepo.GetCustomerByEmail(ctx, "henry@gmail.com")
	// assert.NoError(t, err)
	// assert.NotNil(t, customer)
	// assert.Equal(t, "Henry", customer.Name)
	// assert.Equal(t, "henry@gmail.com", customer.Email)
}

func initApp(t *testing.T) (*api.App, error) {
	t.Helper()

	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:18.4-alpine",
		postgres.WithInitScripts(filepath.Join("..", "..", "db", "migrations", "000001_init.up.sql")),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	app, cleanup, err := api.Initialize()
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		cleanup()
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	return app, nil
}
