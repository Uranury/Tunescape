package leaderboard

import (
	"context"
	"fmt"

	"gitlab.com/Uranury/tunescape/internal/user"
)

var validFeatures = map[string]bool{
	"valence":      true,
	"energy":       true,
	"danceability": true,
	"acousticness": true,
}

type Service interface {
	PushScore(ctx context.Context, feature, userID string, score float64) error
	GetLeaderboard(ctx context.Context, feature string, limit int64) (*LeaderboardResponse, error)
	GetUserRankings(ctx context.Context, userID string) (*UserRankings, error)
}

type service struct {
	store    LeaderboardStore
	userRepo user.Repository
}

func NewService(store LeaderboardStore, userRepo user.Repository) Service {
	return &service{store: store, userRepo: userRepo}
}

func (s *service) PushScore(ctx context.Context, feature, userID string, score float64) error {
	if !validFeatures[feature] {
		return fmt.Errorf("invalid feature: %s", feature)
	}
	return s.store.ZAdd(ctx, fmt.Sprintf("leaderboard:%s", feature), score, userID)
}

func (s *service) GetLeaderboard(ctx context.Context, feature string, limit int64) (*LeaderboardResponse, error) {
	if !validFeatures[feature] {
		return nil, fmt.Errorf("invalid feature: %s", feature)
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	entries, err := s.store.ZRevRangeWithScores(ctx, fmt.Sprintf("leaderboard:%s", feature), 0, limit-1)
	if err != nil {
		return nil, fmt.Errorf("get top n: %w", err)
	}
	userIDs := make([]string, len(entries))
	for i, e := range entries {
		userIDs[i] = fmt.Sprintf("%v", e.Member)
	}
	names, err := s.userRepo.FindDisplayNamesByIDs(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("resolve display names: %w", err)
	}
	resp := &LeaderboardResponse{
		Feature: feature,
		Entries: make([]Entry, len(entries)),
	}
	for i, e := range entries {
		uid := fmt.Sprintf("%v", e.Member)
		resp.Entries[i] = Entry{
			Rank:        i + 1,
			UserID:      uid,
			DisplayName: names[uid],
			Score:       e.Score,
		}
	}
	return resp, nil
}

func (s *service) GetUserRankings(ctx context.Context, userID string) (*UserRankings, error) {
	rankings := &UserRankings{}
	for _, feature := range []string{"valence", "energy", "danceability", "acousticness"} {
		rank, err := s.store.ZRevRank(ctx, fmt.Sprintf("leaderboard:%s", feature), userID)
		if err != nil {
			continue
		}
		r := rank + 1
		switch feature {
		case "valence":
			rankings.Valence = &r
		case "energy":
			rankings.Energy = &r
		case "danceability":
			rankings.Danceability = &r
		case "acousticness":
			rankings.Acousticness = &r
		}
	}
	return rankings, nil
}