package redis

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

func New(config *Config) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr: config.Uri(),
	})

	return &Redis{client}
}

func (r Redis) Get(key string) (string, error) {
	cmd := r.client.Get(context.Background(), key)

	err := cmd.Err()
	if err != nil {
		return "", err
	}

	return cmd.Val(), nil
}

func (r Redis) GetInt64(key string) (int64, error) {
	res, err := r.Get(key)
	if err != nil {
		return -1, err
	}

	return strconv.ParseInt(res, 10, 64)
}

func (r Redis) Set(key string, value interface{}) error {
	return r.client.Set(context.Background(), key, value, 0).Err()
}
