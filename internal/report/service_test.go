package report

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/leaderboard"
	"gitlab.com/Uranury/tunescape/internal/track"
	"gitlab.com/Uranury/tunescape/internal/user"
)

type mockReportRepo struct {
	getLatestSnapshotTopTracksFn func(ctx context.Context, userID uuid.UUID) ([]track.Track, error)
}

func (m *mockReportRepo) GetLatestSnapshotTopTracks(ctx context.Context, userID uuid.UUID) ([]track.Track, error) {
	return m.getLatestSnapshotTopTracksFn(ctx, userID)
}

type mockLeaderboardService struct {
	getUserRankingsFn func(ctx context.Context, userID string) (*leaderboard.UserRankings, error)
}

func (m *mockLeaderboardService) GetUserRankings(ctx context.Context, userID string) (*leaderboard.UserRankings, error) {
	return m.getUserRankingsFn(ctx, userID)
}

func (m *mockLeaderboardService) PushScore(ctx context.Context, feature, userID string, score float64) error {
	return nil
}

func (m *mockLeaderboardService) GetLeaderboard(ctx context.Context, feature string, limit, offset int64) (*leaderboard.LeaderboardResponse, error) {
	return nil, nil
}

type mockUserRepository struct {
	findDisplayNameFn func(ctx context.Context, userID uuid.UUID) (string, error)
}

func (m *mockUserRepository) FindDisplayName(ctx context.Context, userID uuid.UUID) (string, error) {
	return m.findDisplayNameFn(ctx, userID)
}

func (m *mockUserRepository) ConnectSpotify(ctx context.Context, userID uuid.UUID, spotifyID *string, avatarURL, country, product *string) error {
	return nil
}

func (m *mockUserRepository) ClearSpotify(ctx context.Context, userID uuid.UUID) error {
	return nil
}

func (m *mockUserRepository) Create(ctx context.Context, u *user.User) error {
	return nil
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return nil, nil
}

func (m *mockUserRepository) FindByID(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	return nil, nil
}

func (m *mockUserRepository) FindDisplayNamesByIDs(ctx context.Context, userIDs []string) (map[string]string, error) {
	return nil, nil
}

func (m *mockUserRepository) FindAll(ctx context.Context) ([]user.User, error) {
	return nil, nil
}

func (m *mockUserRepository) FindByDisplayName(_ context.Context, _ string) (*user.User, error) {
	return nil, nil
}

func (m *mockUserRepository) FindAvatarURLsByIDs(_ context.Context, _ []string) (map[string]*string, error) {
	return map[string]*string{}, nil
}

// TestReportService_GenerateReport_Success tests successful PDF report generation
func TestReportService_GenerateReport_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	displayName := "Test User"

	tracks := []track.Track{
		{ID: uuid.New(), SpotifyID: "track1", Name: "Song 1", Popularity: 85},
		{ID: uuid.New(), SpotifyID: "track2", Name: "Song 2", Popularity: 80},
		{ID: uuid.New(), SpotifyID: "track3", Name: "Song 3", Popularity: 75},
	}

	repo := &mockReportRepo{
		getLatestSnapshotTopTracksFn: func(ctx context.Context, id uuid.UUID) ([]track.Track, error) {
			if id != userID {
				t.Fatalf("expected userID %s, got %s", userID, id)
			}
			return tracks, nil
		},
	}

	leaderboardSvc := &mockLeaderboardService{
		getUserRankingsFn: func(ctx context.Context, id string) (*leaderboard.UserRankings, error) {
			if id != userID.String() {
				t.Fatalf("expected userID %s, got %s", userID.String(), id)
			}
			rank1 := int64(10)
			rank2 := int64(5)
			return &leaderboard.UserRankings{
				Energy:       &rank1,
				Danceability: &rank2,
			}, nil
		},
	}

	userRepo := &mockUserRepository{
		findDisplayNameFn: func(ctx context.Context, id uuid.UUID) (string, error) {
			if id != userID {
				t.Fatalf("expected userID %s, got %s", userID, id)
			}
			return displayName, nil
		},
	}

	svc := NewService(repo, leaderboardSvc, userRepo)
	pdfData, err := svc.GenerateReport(ctx, userID)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(pdfData) == 0 {
		t.Fatalf("expected non-empty PDF data")
	}
}

// TestReportService_GenerateReport_UserNotFound tests error when user not found
func TestReportService_GenerateReport_UserNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	expectedErr := errors.New("user not found")

	repo := &mockReportRepo{
		getLatestSnapshotTopTracksFn: func(ctx context.Context, id uuid.UUID) ([]track.Track, error) {
			return nil, nil
		},
	}

	userRepo := &mockUserRepository{
		findDisplayNameFn: func(ctx context.Context, id uuid.UUID) (string, error) {
			return "", expectedErr
		},
	}

	svc := NewService(repo, &mockLeaderboardService{}, userRepo)
	_, err := svc.GenerateReport(ctx, userID)

	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}
}

// TestReportService_GenerateReport_NoTracks tests error when no tracks found
func TestReportService_GenerateReport_NoTracks(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	displayName := "Test User"
	expectedErr := errors.New("no snapshot found")

	repo := &mockReportRepo{
		getLatestSnapshotTopTracksFn: func(ctx context.Context, id uuid.UUID) ([]track.Track, error) {
			return nil, expectedErr
		},
	}

	userRepo := &mockUserRepository{
		findDisplayNameFn: func(ctx context.Context, id uuid.UUID) (string, error) {
			return displayName, nil
		},
	}

	svc := NewService(repo, &mockLeaderboardService{}, userRepo)
	_, err := svc.GenerateReport(ctx, userID)

	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}
}

// TestReportService_GenerateReport_LeaderboardServiceError tests graceful handling of leaderboard service errors
func TestReportService_GenerateReport_LeaderboardServiceError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	displayName := "Test User"

	tracks := []track.Track{
		{ID: uuid.New(), SpotifyID: "track1", Name: "Song 1", Popularity: 85},
	}

	repo := &mockReportRepo{
		getLatestSnapshotTopTracksFn: func(ctx context.Context, id uuid.UUID) ([]track.Track, error) {
			return tracks, nil
		},
	}

	leaderboardSvc := &mockLeaderboardService{
		getUserRankingsFn: func(ctx context.Context, id string) (*leaderboard.UserRankings, error) {
			return nil, errors.New("leaderboard service error")
		},
	}

	userRepo := &mockUserRepository{
		findDisplayNameFn: func(ctx context.Context, id uuid.UUID) (string, error) {
			return displayName, nil
		},
	}

	svc := NewService(repo, leaderboardSvc, userRepo)
	pdfData, err := svc.GenerateReport(ctx, userID)

	// Should still succeed with empty rankings
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(pdfData) == 0 {
		t.Fatalf("expected non-empty PDF data")
	}
}

// TestReportService_GenerateReport_EmptyTracks tests report generation with no tracks
func TestReportService_GenerateReport_EmptyTracks(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	displayName := "Test User"

	repo := &mockReportRepo{
		getLatestSnapshotTopTracksFn: func(ctx context.Context, id uuid.UUID) ([]track.Track, error) {
			return []track.Track{}, nil
		},
	}

	leaderboardSvc := &mockLeaderboardService{
		getUserRankingsFn: func(ctx context.Context, id string) (*leaderboard.UserRankings, error) {
			return &leaderboard.UserRankings{}, nil
		},
	}

	userRepo := &mockUserRepository{
		findDisplayNameFn: func(ctx context.Context, id uuid.UUID) (string, error) {
			return displayName, nil
		},
	}

	svc := NewService(repo, leaderboardSvc, userRepo)
	pdfData, err := svc.GenerateReport(ctx, userID)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(pdfData) == 0 {
		t.Fatalf("expected non-empty PDF data")
	}
}
