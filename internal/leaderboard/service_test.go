package leaderboard

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gitlab.com/Uranury/tunescape/internal/user"
)

type mockLeaderboardStore struct {
	zAddFn                func(ctx context.Context, key string, score float64, member string) error
	zRevRangeWithScoresFn func(ctx context.Context, key string, start, stop int64) ([]redis.Z, error)
	zRevRankFn            func(ctx context.Context, key string, member string) (int64, error)
}

func (m *mockLeaderboardStore) ZAdd(ctx context.Context, key string, score float64, member string) error {
	return m.zAddFn(ctx, key, score, member)
}

func (m *mockLeaderboardStore) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return m.zRevRangeWithScoresFn(ctx, key, start, stop)
}

func (m *mockLeaderboardStore) ZRevRank(ctx context.Context, key string, member string) (int64, error) {
	return m.zRevRankFn(ctx, key, member)
}

type mockUserRepo struct {
	findDisplayNamesByIDsFn func(ctx context.Context, userIDs []string) (map[string]string, error)
}

func (m *mockUserRepo) FindDisplayNamesByIDs(ctx context.Context, userIDs []string) (map[string]string, error) {
	return m.findDisplayNamesByIDsFn(ctx, userIDs)
}

func (m *mockUserRepo) ConnectSpotify(ctx context.Context, userID uuid.UUID, spotifyID *string, avatarURL, country, product *string) error {
	return nil
}

func (m *mockUserRepo) ClearSpotify(ctx context.Context, userID uuid.UUID) error {
	return nil
}

func (m *mockUserRepo) Create(ctx context.Context, u *user.User) error {
	return nil
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return nil, nil
}

func (m *mockUserRepo) FindByID(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	return nil, nil
}

func (m *mockUserRepo) FindDisplayName(ctx context.Context, userID uuid.UUID) (string, error) {
	return "", nil
}

func (m *mockUserRepo) FindAll(ctx context.Context) ([]user.User, error) {
	return nil, nil
}

func (m *mockUserRepo) FindByDisplayName(_ context.Context, _ string) (*user.User, error) {
	return nil, nil
}

func (m *mockUserRepo) FindAvatarURLsByIDs(_ context.Context, _ []string) (map[string]*string, error) {
	return map[string]*string{}, nil
}

type mockCache struct {
	getFn    func(ctx context.Context, key string) ([]byte, error)
	setFn    func(ctx context.Context, key string, data []byte, ttl time.Duration) error
	deleteFn func(ctx context.Context, key string) error
	existsFn func(ctx context.Context, key string) (bool, error)
}

func (m *mockCache) Get(ctx context.Context, key string) ([]byte, error) {
	if m.getFn != nil {
		return m.getFn(ctx, key)
	}
	return nil, errors.New("cache miss")
}

func (m *mockCache) Set(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	if m.setFn != nil {
		return m.setFn(ctx, key, data, ttl)
	}
	return nil
}

func (m *mockCache) Delete(ctx context.Context, key string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, key)
	}
	return nil
}

func (m *mockCache) Exists(ctx context.Context, key string) (bool, error) {
	if m.existsFn != nil {
		return m.existsFn(ctx, key)
	}
	return false, nil
}

// TestLeaderboardService_PushScore_Success tests successful score push
func TestLeaderboardService_PushScore_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	feature := "energy"
	userID := "user-123"
	score := 95.5

	var capturedKey string
	var capturedScore float64
	var capturedMember string

	store := &mockLeaderboardStore{
		zAddFn: func(ctx context.Context, key string, s float64, member string) error {
			capturedKey = key
			capturedScore = s
			capturedMember = member
			return nil
		},
	}

	svc := NewService(store, &mockUserRepo{}, &mockCache{})
	err := svc.PushScore(ctx, feature, userID, score)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if capturedKey != "leaderboard:energy" {
		t.Fatalf("expected key 'leaderboard:energy', got %q", capturedKey)
	}
	if capturedScore != score {
		t.Fatalf("expected score %f, got %f", score, capturedScore)
	}
	if capturedMember != userID {
		t.Fatalf("expected member %q, got %q", userID, capturedMember)
	}
}

// TestLeaderboardService_PushScore_InvalidFeature tests invalid feature rejection
func TestLeaderboardService_PushScore_InvalidFeature(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := &mockLeaderboardStore{}
	svc := NewService(store, &mockUserRepo{}, &mockCache{})

	err := svc.PushScore(ctx, "invalid_feature", "user-123", 50.0)
	if err == nil {
		t.Fatalf("expected error for invalid feature")
	}
}

// TestLeaderboardService_GetLeaderboard_Success tests successful leaderboard retrieval
func TestLeaderboardService_GetLeaderboard_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	feature := "danceability"
	limit := int64(10)
	offset := int64(0)

	entries := []redis.Z{
		{Score: 100.0, Member: "user-1"},
		{Score: 95.0, Member: "user-2"},
		{Score: 90.0, Member: "user-3"},
	}

	displayNames := map[string]string{
		"user-1": "Alice",
		"user-2": "Bob",
		"user-3": "Charlie",
	}

	store := &mockLeaderboardStore{
		zRevRangeWithScoresFn: func(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
			if key != "leaderboard:danceability" {
				t.Fatalf("expected key 'leaderboard:danceability', got %q", key)
			}
			if start != 0 || stop != 9 {
				t.Fatalf("expected range [0, 9], got [%d, %d]", start, stop)
			}
			return entries, nil
		},
	}

	userRepo := &mockUserRepo{
		findDisplayNamesByIDsFn: func(ctx context.Context, userIDs []string) (map[string]string, error) {
			return displayNames, nil
		},
	}

	svc := NewService(store, userRepo, &mockCache{})
	result, err := svc.GetLeaderboard(ctx, feature, limit, offset)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Feature != feature {
		t.Fatalf("expected feature %q, got %q", feature, result.Feature)
	}
	if len(result.Entries) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(result.Entries))
	}

	for i, entry := range result.Entries {
		if entry.Rank != i+1 {
			t.Fatalf("entry %d: expected rank %d, got %d", i, i+1, entry.Rank)
		}
		if entry.DisplayName != displayNames[entry.UserID] {
			t.Fatalf("entry %d: expected display name %q, got %q", i, displayNames[entry.UserID], entry.DisplayName)
		}
	}
}

// TestLeaderboardService_GetLeaderboard_InvalidFeature tests invalid feature rejection
func TestLeaderboardService_GetLeaderboard_InvalidFeature(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := &mockLeaderboardStore{}
	svc := NewService(store, &mockUserRepo{}, &mockCache{})

	result, err := svc.GetLeaderboard(ctx, "unknown_feature", 10, 0)
	if err == nil {
		t.Fatalf("expected error for invalid feature")
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

// TestLeaderboardService_GetLeaderboard_DefaultLimit tests limit validation and default
func TestLeaderboardService_GetLeaderboard_DefaultLimit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	tests := []struct {
		name       string
		inputLimit int64
		wantLimit  int64
	}{
		{name: "negative limit", inputLimit: -1, wantLimit: 10},
		{name: "zero limit", inputLimit: 0, wantLimit: 10},
		{name: "too large limit", inputLimit: 101, wantLimit: 10},
		{name: "valid limit", inputLimit: 5, wantLimit: 5},
		{name: "max valid limit", inputLimit: 100, wantLimit: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockLeaderboardStore{
				zRevRangeWithScoresFn: func(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
					expectedStop := start + tt.wantLimit - 1
					if stop != expectedStop {
						t.Errorf("expected stop %d, got %d", expectedStop, stop)
					}
					return []redis.Z{}, nil
				},
			}

			userRepo := &mockUserRepo{
				findDisplayNamesByIDsFn: func(ctx context.Context, userIDs []string) (map[string]string, error) {
					return map[string]string{}, nil
				},
			}

			svc := NewService(store, userRepo, &mockCache{})
			_, _ = svc.GetLeaderboard(ctx, "energy", tt.inputLimit, 0)
		})
	}
}

// TestLeaderboardService_GetUserRankings_Success tests successful user rankings retrieval
func TestLeaderboardService_GetUserRankings_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := "user-123"

	store := &mockLeaderboardStore{
		zRevRankFn: func(ctx context.Context, key string, member string) (int64, error) {
			if member != userID {
				t.Fatalf("expected member %q, got %q", userID, member)
			}
			// Return different ranks for each feature
			switch key {
			case "leaderboard:valence":
				return 5, nil
			case "leaderboard:energy":
				return 10, nil
			case "leaderboard:danceability":
				return 3, nil
			case "leaderboard:acousticness":
				return 15, nil
			}
			return 0, errors.New("not found")
		},
	}

	svc := NewService(store, &mockUserRepo{}, &mockCache{})
	result, err := svc.GetUserRankings(ctx, userID)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Valence == nil || *result.Valence != 6 {
		t.Fatalf("expected valence rank 6, got %v", result.Valence)
	}
	if result.Energy == nil || *result.Energy != 11 {
		t.Fatalf("expected energy rank 11, got %v", result.Energy)
	}
	if result.Danceability == nil || *result.Danceability != 4 {
		t.Fatalf("expected danceability rank 4, got %v", result.Danceability)
	}
	if result.Acousticness == nil || *result.Acousticness != 16 {
		t.Fatalf("expected acousticness rank 16, got %v", result.Acousticness)
	}
}

// TestLeaderboardService_GetUserRankings_NotRanked tests when user has no rankings
func TestLeaderboardService_GetUserRankings_NotRanked(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := "user-unknown"

	store := &mockLeaderboardStore{
		zRevRankFn: func(ctx context.Context, key string, member string) (int64, error) {
			return 0, errors.New("member not found")
		},
	}

	svc := NewService(store, &mockUserRepo{}, &mockCache{})
	result, err := svc.GetUserRankings(ctx, userID)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	// All rankings should be nil
	if result.Valence != nil || result.Energy != nil || result.Danceability != nil || result.Acousticness != nil {
		t.Fatalf("expected all rankings to be nil")
	}
}
