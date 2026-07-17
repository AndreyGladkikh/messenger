package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"messenger/messenger/internal/platform/config"
	"time"
)

func NewPool(cfg *config.Config) (*sql.DB, func(), error) {
	dsn := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		cfg.Database.Schema,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	pool, err := sql.Open(cfg.Database.Driver, dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("postgres open: %w", err)
	}

	cleanup := func() {
		pool.Close()
	}
	
	pool.SetMaxOpenConns(50)
	pool.SetMaxIdleConns(10)
	pool.SetConnMaxLifetime(1 * time.Hour)
	pool.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		return nil, cleanup, fmt.Errorf("postgres ping: %w", err)
	}

	return pool, cleanup, nil
}
