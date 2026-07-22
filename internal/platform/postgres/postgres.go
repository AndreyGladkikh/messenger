package postgres

import (
	"context"
	"fmt"
	"messenger/messenger/internal/platform/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// database/sql
// func NewPool(cfg *config.Config) (*sql.DB, func(), error) {
// 	dsn := fmt.Sprintf(
// 		"%s://%s:%s@%s:%s/%s",
// 		cfg.Database.Schema,
// 		cfg.Database.User,
// 		cfg.Database.Password,
// 		cfg.Database.Host,
// 		cfg.Database.Port,
// 		cfg.Database.Name,
// 	)

// 	pool, err := sql.Open(cfg.Database.Driver, dsn)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("postgres open: %w", err)
// 	}

// 	cleanup := func() {
// 		pool.Close()
// 	}
	
// 	pool.SetMaxOpenConns(50)
// 	pool.SetMaxIdleConns(10)
// 	pool.SetConnMaxLifetime(1 * time.Hour)
// 	pool.SetConnMaxIdleTime(5 * time.Minute)

// 	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
// 	defer cancel()

// 	if err := pool.PingContext(ctx); err != nil {
// 		return nil, cleanup, fmt.Errorf("postgres ping: %w", err)
// 	}

// 	return pool, cleanup, nil
// }

func NewPool(cfg *config.Config) (*pgxpool.Pool, func(), error) {
	dsn := fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s",
		cfg.Database.Schema,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	pgxPoolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse pgx pool config: %w", err)
	}

	pgxPoolConfig.MaxConnLifetime = 1 * time.Hour
	pgxPoolConfig.MaxConnIdleTime = 5 * time.Minute
	pgxPoolConfig.MaxConns = 50

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create pgx pool: %w", err)
	}
	cleanup := func() {
		pool.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, cleanup, fmt.Errorf("database ping: %w", err)
	}
	
	return pool, cleanup, nil
}