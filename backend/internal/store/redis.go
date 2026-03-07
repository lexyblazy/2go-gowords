package store

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	LEADERBOARD_DAILY  = "leaderboard:daily"
	LEADERBOARD_WEEKLY = "leaderboard:weekly"

	USERS_MONIKERS = "users_monikers"
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

func (rs *RedisStore) CacheUserMoniker(ctx context.Context, userId string, moniker string) error {
	return rs.conn.HSet(ctx, USERS_MONIKERS, userId, moniker).Err()
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

func (rs *RedisStore) UpdateLeaderBoards(ctx context.Context, scoresMap map[string]int) error {
	pipe := rs.conn.Pipeline()

	for userId, score := range scoresMap {

		pipe.ZAddArgs(ctx, LEADERBOARD_DAILY, redis.ZAddArgs{
			GT: true,
			Members: []redis.Z{
				{

					Score:  float64(score),
					Member: userId,
				},
			},
		})

		pipe.ZAddArgs(ctx, LEADERBOARD_WEEKLY, redis.ZAddArgs{
			GT: true,
			Members: []redis.Z{
				{

					Score:  float64(score),
					Member: userId,
				},
			},
		})

	}

	_, err := pipe.Exec(ctx)
	return err
}

func (rs *RedisStore) GetDailyTop(ctx context.Context, rdb *redis.Client) ([]redis.Z, error) {
	return rdb.ZRevRangeWithScores(ctx, LEADERBOARD_DAILY, 0, 9).Result()
}

func (rs *RedisStore) GetWeeklyTop(ctx context.Context, rdb *redis.Client) ([]redis.Z, error) {
	return rdb.ZRevRangeWithScores(ctx, LEADERBOARD_WEEKLY, 0, 9).Result()
}
