package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
	ctx    context.Context
}

func New(config *Config) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr: config.Uri(),
	})

	return &Redis{
		client: client,
		ctx:    context.Background(),
	}
}

func (redis *Redis) Get(key string) *redis.StringCmd {
	return redis.client.Get(redis.ctx, key)
}

func (redis *Redis) Set(key string, value interface{}) *redis.StatusCmd {
	return redis.client.Set(redis.ctx, key, value, 0)
}
