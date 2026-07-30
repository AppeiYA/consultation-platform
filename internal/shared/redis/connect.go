package redis

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/shared/config"
	"github.com/redis/go-redis/v9"
)

func Connect(cfg config.RedisConfig) (*Redis, error) {
	opts := &redis.Options{
		Addr:     cfg.Address,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	client := redis.NewClient(opts)

	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return New(client), nil
}