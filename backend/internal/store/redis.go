package store

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	conn *redis.Client
}

func NewRedisStore(url string) (*RedisStore, error) {

	opts, err := redis.ParseURL(url)

	if err != nil {
		return nil, errors.New("failed to parse redis url")
	}
	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisStore{}, nil
}

func (rs *RedisStore) Set(ctx context.Context, key string, value string) {

}
