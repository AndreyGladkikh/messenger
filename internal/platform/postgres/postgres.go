package postgres

import (
	"context"
	"fmt"
	"messenger/messenger/internal/platform/config"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	pool, err := pgxpool.NewWithConfig(context.Background(), &pgxpool.Config{
		ConnConfig: &pgx.ConnConfig{
			Config: pgconn.Config{
				Host: cfg.Database.Host,
				Port: cfg.Database.Port,
				Database: cfg.Database.Name,
				User: cfg.Database.User,
				Password: cfg.Database.Password,
				ConnectTimeout: 3*time.Second,
			},
		},
		MaxConnLifetime: 1 * time.Hour,
		MaxConnIdleTime: 5 * time.Minute,
		MaxConns: 50,
		MinIdleConns: 10,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to create connection pool: %w", err)
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