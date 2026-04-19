package leaderboard

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type LeaderboardStore interface {
	ZAdd(ctx context.Context, key string, score float64, member string) error
	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error)
	ZRevRank(ctx context.Context, key string, member string) (int64, error)
}

type redisStore struct {
	rds *redis.Client
}

func NewStore(rds *redis.Client) LeaderboardStore {
	return &redisStore{rds: rds}
}

func (s *redisStore) ZAdd(ctx context.Context, key string, score float64, member string) error {
	return s.rds.ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Err()
}

func (s *redisStore) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return s.rds.ZRevRangeWithScores(ctx, key, start, stop).Result()
}

func (s *redisStore) ZRevRank(ctx context.Context, key string, member string) (int64, error) {
	return s.rds.ZRevRank(ctx, key, member).Result()
}