package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

func NewPool(ctx context.Context) (*sqlx.DB, func(), error) {
	db := sqlx.MustConnect("pgx", "postgres://app:secret@postgres:5432/app")

	cleanup := func () {
		db.Close()
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(1*time.Hour)
	db.SetConnMaxIdleTime(5*time.Minute)

	return db, cleanup, nil
}

// func NewPool(ctx context.Context) (*sql.DB, func(), error) {
// 	pool, err := sql.Open("pgx", "postgres://app:secret@postgres:5432/app")
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	db := sqlx.NewDb(pool, "pgx")
// 	db.Close()

// 	cleanup := func () {
// 		pool.Close()
// 	}

// 	pool.SetMaxOpenConns(50)
// 	pool.SetMaxIdleConns(10)
// 	pool.SetConnMaxLifetime(time.Minute)

// 	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
// 	defer cancel()

// 	if err := pool.PingContext(ctx); err != nil {
// 		return nil, cleanup, fmt.Errorf("unable to connect to database: %w", err)
// 	}

// 	return pool, cleanup, nil
// }