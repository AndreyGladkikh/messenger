package redis

import (
	"context"
	"fmt"
	"messenger/messenger/internal/platform/config"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg config.Config) (*redis.Client, func(), error) {
	opt, err := redis.ParseURL(cfg.Cache.Dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("redis: failed to parse url: %w", err)
	}
	rdb := redis.NewClient(opt)

	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     cfg.Cache.Dsn,
	// 	Password: "", // no password
	// 	DB:       0,  // use default DB
	// 	Protocol: 2,
	// })

	cleanup := func() {
		rdb.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return nil, nil, fmt.Errorf("redis: failed to ping: %w", err)
	}

	return rdb, cleanup, nil
}
