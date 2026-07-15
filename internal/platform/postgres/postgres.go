package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"messenger/messenger/internal/platform/config"
	"time"
)

func NewPool(ctx context.Context, cfg *config.Config) (*sql.DB, func(), error) {
	pool, err := sql.Open(cfg.DB.Driver, cfg.DB.DSN)
	// pool, err := sql.Open("pgx", "postgres://app:secret@postgres:5432/app")
	if err != nil {
		return nil, nil, err
	}

	cleanup := func () {
		pool.Close()
	}
	pool.SetMaxOpenConns(50)
	pool.SetMaxIdleConns(10)
	pool.SetConnMaxLifetime(1*time.Hour)
	pool.SetConnMaxIdleTime(5*time.Minute)

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		return nil, cleanup, fmt.Errorf("unable to connect to database: %w", err)
	}

	return pool, cleanup, nil
}