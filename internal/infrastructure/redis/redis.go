package redis

import (
	"context"
	"fmt"
	"messenger/messenger/internal/infrastructure/config"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *config.Config) (*redis.Client, func(), error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // use default DB
		Protocol: 2,
	})

	cleanup := func() {
		rdb.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, func() {}, fmt.Errorf("redis: failed to ping: %w", err)
	}

	return rdb, cleanup, nil
}
