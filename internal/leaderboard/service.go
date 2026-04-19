package leaderboard

import (
	"context"
	"fmt"
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
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) PushScore(ctx context.Context, feature, userID string, score float64) error {
	if !validFeatures[feature] {
		return fmt.Errorf("invalid feature: %s", feature)
	}
	return s.repo.PushScore(ctx, feature, userID, score)
}

func (s *service) GetLeaderboard(ctx context.Context, feature string, limit int64) (*LeaderboardResponse, error) {
	if !validFeatures[feature] {
		return nil, fmt.Errorf("invalid feature: %s", feature)
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	entries, err := s.repo.GetTopN(ctx, feature, limit)
	if err != nil {
		return nil, fmt.Errorf("get top n: %w", err)
	}
	userIDs := make([]string, len(entries))
	for i, e := range entries {
		userIDs[i] = fmt.Sprintf("%v", e.Member)
	}
	names, err := s.repo.ResolveDisplayNames(ctx, userIDs)
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
		rank, err := s.repo.GetUserRank(ctx, feature, userID)
		if err != nil {
			continue
		}
		r := rank
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