package leaderboard

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

type Repository interface {
	PushScore(ctx context.Context, feature, userID string, score float64) error
	GetTopN(ctx context.Context, feature string, limit int64) ([]redis.Z, error)
	GetUserRank(ctx context.Context, feature, userID string) (int64, error)
	ResolveDisplayNames(ctx context.Context, userIDs []string) (map[string]string, error)
}

type repository struct {
	store LeaderboardStore
	exec  database.Executor
}

func NewRepository(store LeaderboardStore, exec database.Executor) Repository {
	return &repository{store: store, exec: exec}
}

func (r *repository) PushScore(ctx context.Context, feature, userID string, score float64) error {
	return r.store.ZAdd(ctx, fmt.Sprintf("leaderboard:%s", feature), score, userID)
}

func (r *repository) GetTopN(ctx context.Context, feature string, limit int64) ([]redis.Z, error) {
	return r.store.ZRevRangeWithScores(ctx, fmt.Sprintf("leaderboard:%s", feature), 0, limit-1)
}

func (r *repository) GetUserRank(ctx context.Context, feature, userID string) (int64, error) {
	rank, err := r.store.ZRevRank(ctx, fmt.Sprintf("leaderboard:%s", feature), userID)
	if err != nil {
		return 0, err
	}
	return rank + 1, nil
}

func (r *repository) ResolveDisplayNames(ctx context.Context, userIDs []string) (map[string]string, error) {
	if len(userIDs) == 0 {
		return map[string]string{}, nil
	}
	query := `SELECT id::text AS id, display_name FROM users WHERE id::text = ANY($1)`
	var rows []struct {
		ID          string `db:"id"`
		DisplayName string `db:"display_name"`
	}
	if err := r.exec.SelectContext(ctx, &rows, query, pq.Array(userIDs)); err != nil {
		return nil, err
	}
	result := make(map[string]string, len(rows))
	for _, row := range rows {
		result[row.ID] = row.DisplayName
	}
	return result, nil
}