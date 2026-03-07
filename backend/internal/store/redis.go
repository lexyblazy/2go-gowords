package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	LEADERBOARD_DAILY                = "leaderboard:daily"
	LEADERBOARD_WEEKLY               = "leaderboard:weekly"
	LEADERBOARD_ALL_TIME_HIGH_SCORES = "leaderboard:all_time_high_scores"

	USERS_MONIKERS = "users_monikers"
)

type RedisStore struct {
	conn *redis.Client
}

type LeaderboardEntry struct {
	UserId  string `json:"userId"`
	Moniker string `json:"moniker"`
	Score   int    `json:"score"`
	Rank    int    `json:"rank"`
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

func getLeaderBoardKeys() (daily string, weekly string) {
	now := time.Now()

	daily = fmt.Sprintf("%s:%s", LEADERBOARD_DAILY, now.Format("2006-01-02"))

	year, week := now.ISOWeek()

	weekly = fmt.Sprintf("%s:%d-%02d", LEADERBOARD_WEEKLY, year, week)

	return daily, weekly

}

func (rs *RedisStore) UpdateLeaderBoards(ctx context.Context, scoresMap map[string]int) error {
	pipe := rs.conn.Pipeline()

	dailyKey, weeklyKey := getLeaderBoardKeys()

	for userId, score := range scoresMap {

		pipe.ZIncrBy(ctx, dailyKey, float64(score), userId)
		pipe.ZIncrBy(ctx, weeklyKey, float64(score), userId)
		pipe.ZAddArgs(ctx, LEADERBOARD_ALL_TIME_HIGH_SCORES, redis.ZAddArgs{
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

func (rs *RedisStore) aggregateLeaderBoardFromSet(ctx context.Context, set []redis.Z) ([]*LeaderboardEntry, error) {

	entries := make([]*LeaderboardEntry, 0)
	if len(set) == 0 {
		return entries, nil
	}
	ids := make([]string, 0, len(set))

	for i, z := range set {
		userId := z.Member.(string)
		entry := &LeaderboardEntry{
			UserId: userId,
			Score:  int(z.Score),
			Rank:   i + 1,
		}
		ids = append(ids, userId)
		entries = append(entries, entry)
	}

	monikers, err := rs.getUserMonikers(ctx, ids...)
	idsToMonikers := make(map[string]string)

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(monikers); i++ {

		moniker := monikers[i]

		if moniker == nil || moniker == redis.Nil {
			moniker = fmt.Sprintf("Guest:%s", ids[i])
		}
		idsToMonikers[ids[i]] = moniker.(string)
	}

	// modify the entry in place
	for i := 0; i < len(entries); i++ {
		entries[i].Moniker = idsToMonikers[entries[i].UserId]
	}

	return entries, nil

}

func (rs *RedisStore) GetDailyLeaderBoard(ctx context.Context) ([]*LeaderboardEntry, error) {

	dailyKey, _ := getLeaderBoardKeys()
	set, err := rs.conn.ZRevRangeWithScores(ctx, dailyKey, 0, 9).Result()

	if err != nil {
		return nil, err
	}

	return rs.aggregateLeaderBoardFromSet(ctx, set)

}

func (rs *RedisStore) GetWeeklyLeaderboard(ctx context.Context) ([]*LeaderboardEntry, error) {

	_, weeklyKey := getLeaderBoardKeys()
	set, err := rs.conn.ZRevRangeWithScores(ctx, weeklyKey, 0, 9).Result()

	if err != nil {
		return nil, err
	}

	return rs.aggregateLeaderBoardFromSet(ctx, set)

}

func (rs *RedisStore) GetAllTimeHighScores(ctx context.Context) ([]*LeaderboardEntry, error) {

	set, err := rs.conn.ZRevRangeWithScores(ctx, LEADERBOARD_ALL_TIME_HIGH_SCORES, 0, 9).Result()

	if err != nil {
		return nil, err
	}

	return rs.aggregateLeaderBoardFromSet(ctx, set)

}

func (rs *RedisStore) getUserMonikers(ctx context.Context, ids ...string) ([]any, error) {
	return rs.conn.HMGet(ctx, USERS_MONIKERS, ids...).Result()
}
