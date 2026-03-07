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

	return &RedisStore{conn: rdb}, nil
}

func (rs *RedisStore) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return rs.conn.Set(ctx, key, value, ttl).Err()
}

func (rs *RedisStore) Delete(ctx context.Context, key string) error {
	return rs.conn.Del(ctx, key).Err()
}

func (rs *RedisStore) Get(ctx context.Context, key string) (string, error) {
	val, err := rs.conn.Get(ctx, key).Result()

	if err != nil {
		return "", err
	}

	if err == redis.Nil {
		return "", nil
	}

	return val, nil

}
