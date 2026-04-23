package trends

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockTrendsRepo struct {
	getTrendsByUserIDFn func(ctx context.Context, userID uuid.UUID) ([]SnapshotPoint, error)
}

func (m *mockTrendsRepo) GetTrendsByUserID(ctx context.Context, userID uuid.UUID) ([]SnapshotPoint, error) {
	return m.getTrendsByUserIDFn(ctx, userID)
}

// TestTrendsService_GetTrends_Success tests successful trends retrieval
func TestTrendsService_GetTrends_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()

	points := []SnapshotPoint{
		{
			SnapshotID:       uuid.New(),
			CreatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Danceability:     0.7,
			Valence:          0.6,
			Energy:           0.8,
			Acousticness:     0.3,
			Instrumentalness: 0.1,
			Liveness:         0.2,
			Speechiness:      0.05,
			Tempo:            130.0,
			Loudness:         -5.0,
			TracksCount:      50,
		},
		{
			SnapshotID:       uuid.New(),
			CreatedAt:        time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			Danceability:     0.75,
			Valence:          0.65,
			Energy:           0.78,
			Acousticness:     0.25,
			Instrumentalness: 0.12,
			Liveness:         0.22,
			Speechiness:      0.06,
			Tempo:            128.0,
			Loudness:         -4.5,
			TracksCount:      48,
		},
	}

	repo := &mockTrendsRepo{
		getTrendsByUserIDFn: func(ctx context.Context, id uuid.UUID) ([]SnapshotPoint, error) {
			if id != userID {
				t.Fatalf("expected userID %s, got %s", userID, id)
			}
			return points, nil
		},
	}

	svc := NewService(repo)
	result, err := svc.GetTrends(ctx, userID)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.UserID != userID {
		t.Fatalf("expected userID %s, got %s", userID, result.UserID)
	}
	if len(result.Points) != len(points) {
		t.Fatalf("expected %d points, got %d", len(points), len(result.Points))
	}

	for i, point := range result.Points {
		if point.Danceability != points[i].Danceability {
			t.Fatalf("point %d: expected danceability %f, got %f", i, points[i].Danceability, point.Danceability)
		}
		if point.Valence != points[i].Valence {
			t.Fatalf("point %d: expected valence %f, got %f", i, points[i].Valence, point.Valence)
		}
		if point.Energy != points[i].Energy {
			t.Fatalf("point %d: expected energy %f, got %f", i, points[i].Energy, point.Energy)
		}
		if point.TracksCount != points[i].TracksCount {
			t.Fatalf("point %d: expected tracks count %d, got %d", i, points[i].TracksCount, point.TracksCount)
		}
	}
}

// TestTrendsService_GetTrends_EmptyTrends tests retrieval with no trends
func TestTrendsService_GetTrends_EmptyTrends(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()

	repo := &mockTrendsRepo{
		getTrendsByUserIDFn: func(ctx context.Context, id uuid.UUID) ([]SnapshotPoint, error) {
			return []SnapshotPoint{}, nil
		},
	}

	svc := NewService(repo)
	result, err := svc.GetTrends(ctx, userID)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.UserID != userID {
		t.Fatalf("expected userID %s, got %s", userID, result.UserID)
	}
	if len(result.Points) != 0 {
		t.Fatalf("expected 0 points, got %d", len(result.Points))
	}
}

// TestTrendsService_GetTrends_RepositoryError tests error handling
func TestTrendsService_GetTrends_RepositoryError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	expectedErr := errors.New("database connection error")

	repo := &mockTrendsRepo{
		getTrendsByUserIDFn: func(ctx context.Context, id uuid.UUID) ([]SnapshotPoint, error) {
			return nil, expectedErr
		},
	}

	svc := NewService(repo)
	_, err := svc.GetTrends(ctx, userID)

	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error, got %v", err)
	}
}

// TestTrendsService_GetTrends_LargeDataset tests trends with many data points
func TestTrendsService_GetTrends_LargeDataset(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()

	// Generate 52 weeks of data (1 year)
	var points []SnapshotPoint
	for i := 0; i < 52; i++ {
		points = append(points, SnapshotPoint{
			SnapshotID:   uuid.New(),
			CreatedAt:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, i*7),
			Danceability: 0.5 + float64(i%10)*0.05,
			Valence:      0.6 + float64(i%8)*0.03,
			Energy:       0.7 + float64(i%6)*0.02,
			TracksCount:  50 + i%20,
		})
	}

	repo := &mockTrendsRepo{
		getTrendsByUserIDFn: func(ctx context.Context, id uuid.UUID) ([]SnapshotPoint, error) {
			return points, nil
		},
	}

	svc := NewService(repo)
	result, err := svc.GetTrends(ctx, userID)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(result.Points) != 52 {
		t.Fatalf("expected 52 points, got %d", len(result.Points))
	}
}

// TestTrendsService_GetTrends_SinglePoint tests trends with only one data point
func TestTrendsService_GetTrends_SinglePoint(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	snapshotID := uuid.New()
	now := time.Now()

	points := []SnapshotPoint{
		{
			SnapshotID:       snapshotID,
			CreatedAt:        now,
			Danceability:     0.7,
			Valence:          0.6,
			Energy:           0.8,
			Acousticness:     0.3,
			Instrumentalness: 0.1,
			Liveness:         0.2,
			Speechiness:      0.05,
			Tempo:            130.0,
			Loudness:         -5.0,
			TracksCount:      50,
		},
	}

	repo := &mockTrendsRepo{
		getTrendsByUserIDFn: func(ctx context.Context, id uuid.UUID) ([]SnapshotPoint, error) {
			return points, nil
		},
	}

	svc := NewService(repo)
	result, err := svc.GetTrends(ctx, userID)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.Points) != 1 {
		t.Fatalf("expected 1 point, got %d", len(result.Points))
	}
	if result.Points[0].SnapshotID != snapshotID {
		t.Fatalf("expected snapshot ID %s, got %s", snapshotID, result.Points[0].SnapshotID)
	}
}
