package leaderboard

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gitlab.com/Uranury/tunescape/internal/cache"
	"gitlab.com/Uranury/tunescape/internal/user"
)

const displayNameTTL = time.Hour

var validFeatures = map[string]bool{
	"valence":      true,
	"energy":       true,
	"danceability": true,
	"acousticness": true,
}

type Service interface {
	PushScore(ctx context.Context, feature, userID string, score float64) error
	GetLeaderboard(ctx context.Context, feature string, limit, offset int64) (*LeaderboardResponse, error)
	GetUserRankings(ctx context.Context, userID string) (*UserRankings, error)
}

type service struct {
	store    LeaderboardStore
	userRepo user.Repository
	cache    cache.Cache
}

func NewService(store LeaderboardStore, userRepo user.Repository, cache cache.Cache) Service {
	return &service{store: store, userRepo: userRepo, cache: cache}
}

func (s *service) PushScore(ctx context.Context, feature, userID string, score float64) error {
	if !validFeatures[feature] {
		return fmt.Errorf("invalid feature: %s", feature)
	}
	return s.store.ZAdd(ctx, fmt.Sprintf("leaderboard:%s", feature), score, userID)
}

func (s *service) resolveDisplayNames(ctx context.Context, userIDs []string) (map[string]string, error) {
	names := make(map[string]string, len(userIDs))
	var missing []string

	for _, uid := range userIDs {
		val, err := s.cache.Get(ctx, "displayname:"+uid)
		if err == nil && val != nil {
			names[uid] = string(val)
		} else {
			missing = append(missing, uid)
		}
	}

	if len(missing) == 0 {
		return names, nil
	}

	fetched, err := s.userRepo.FindDisplayNamesByIDs(ctx, missing)
	if err != nil {
		return nil, fmt.Errorf("resolve display names from db: %w", err)
	}

	for uid, name := range fetched {
		names[uid] = name
		if err := s.cache.Set(ctx, "displayname:"+uid, []byte(name), displayNameTTL); err != nil {
			slog.Warn("failed to cache display name", "user_id", uid, "error", err)
		}
	}

	return names, nil
}

func (s *service) GetLeaderboard(ctx context.Context, feature string, limit, offset int64) (*LeaderboardResponse, error) {
	if !validFeatures[feature] {
		return nil, fmt.Errorf("invalid feature: %s", feature)
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	start := offset
	stop := offset + limit - 1
	entries, err := s.store.ZRevRangeWithScores(ctx, fmt.Sprintf("leaderboard:%s", feature), start, stop)
	if err != nil {
		return nil, fmt.Errorf("get top n: %w", err)
	}
	userIDs := make([]string, len(entries))
	for i, e := range entries {
		userIDs[i] = fmt.Sprintf("%v", e.Member)
	}
	names, err := s.resolveDisplayNames(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("resolve display names: %w", err)
	}
	avatars, err := s.userRepo.FindAvatarURLsByIDs(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("resolve avatars: %w", err)
	}
	resp := &LeaderboardResponse{
		Feature: feature,
		Entries: make([]Entry, len(entries)),
	}
	for i, e := range entries {
		uid := fmt.Sprintf("%v", e.Member)
		resp.Entries[i] = Entry{
			Rank:        int(offset) + i + 1,
			UserID:      uid,
			DisplayName: names[uid],
			AvatarURL:   avatars[uid],
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
