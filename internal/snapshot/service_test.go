package snapshot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gitlab.com/Uranury/tunescape/internal/spotify"
	"gitlab.com/Uranury/tunescape/internal/track"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

type mockCache struct{}

func (m *mockCache) Get(_ context.Context, _ string) ([]byte, error)                  { return nil, nil }
func (m *mockCache) Set(_ context.Context, _ string, _ []byte, _ time.Duration) error { return nil }
func (m *mockCache) Delete(_ context.Context, _ string) error                         { return nil }
func (m *mockCache) Exists(_ context.Context, _ string) (bool, error)                 { return false, nil }

type mockSpotifyService struct {
	getTopTracksFn func(ctx context.Context, userID uuid.UUID, limit int, timeRange spotify.TimeRange) ([]track.Track, error)
}

func (m *mockSpotifyService) GetTopTracks(ctx context.Context, userID uuid.UUID, limit int, timeRange spotify.TimeRange) ([]track.Track, error) {
	return m.getTopTracksFn(ctx, userID, limit, timeRange)
}
func (m *mockSpotifyService) AuthURL(_ string) string { return "" }
func (m *mockSpotifyService) ConnectAccount(_ context.Context, _ uuid.UUID, _ string) error {
	return nil
}
func (m *mockSpotifyService) Disconnect(_ context.Context, _ uuid.UUID) error { return nil }
func (m *mockSpotifyService) GetValidToken(_ context.Context, _ uuid.UUID) (string, error) {
	return "", nil
}
func (m *mockSpotifyService) UpsertTokens(_ context.Context, _ uuid.UUID, _, _ string, _ time.Time) error {
	return nil
}

func newDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, database.TxProvider) {
	t.Helper()
	rawDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	db := sqlx.NewDb(rawDB, "sqlmock")
	return rawDB, mock, database.NewTxProvider(db)
}

func makeTopTracks(n int) []track.Track {
	tracks := make([]track.Track, n)
	for i := range tracks {
		tracks[i] = track.Track{
			SpotifyID:  fmt.Sprintf("spotify-id-%d", i+1),
			Name:       fmt.Sprintf("Track %d", i+1),
			Popularity: i * 2,
		}
	}
	return tracks
}

// Success: fetch tracks -> INSERT snapshots -> INSERT tracks (upsert) -> INSERT snapshot_tracks.
func TestSnapshotService_CreateSnapshot_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	snapID := uuid.New()
	createdAt := time.Now()
	topTracks := makeTopTracks(3)
	trackIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}

	rawDB, mock, txProvider := newDB(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO snapshots`).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "created_at"}).
				AddRow(snapID, createdAt),
		)
	for i, tr := range topTracks {
		mock.ExpectQuery(`INSERT INTO tracks`).
			WithArgs(tr.SpotifyID, tr.Name, tr.Popularity).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).AddRow(trackIDs[i]),
			)
		mock.ExpectExec(`INSERT INTO snapshot_tracks`).
			WithArgs(snapID, trackIDs[i], i+1).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
	mock.ExpectCommit()
	mock.ExpectClose()

	spotifySvc := &mockSpotifyService{
		getTopTracksFn: func(_ context.Context, gotUserID uuid.UUID, limit int, timeRange spotify.TimeRange) ([]track.Track, error) {
			if gotUserID != userID {
				t.Fatalf("expected userID %s, got %s", userID, gotUserID)
			}
			if limit != topTracksLimit {
				t.Fatalf("expected limit %d, got %d", topTracksLimit, limit)
			}
			if timeRange != spotify.MediumTerm {
				t.Fatalf("expected timeRange %s, got %s", spotify.MediumTerm, timeRange)
			}
			return topTracks, nil
		},
	}

	svc := NewService(nil, spotifySvc, txProvider, &mockCache{}, slog.Default())
	snap, err := svc.CreateSnapshot(ctx, userID, spotify.MediumTerm)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if snap == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if snap.ID != snapID {
		t.Fatalf("expected snapID %s, got %s", snapID, snap.ID)
	}
	if snap.UserID != userID {
		t.Fatalf("expected userID %s, got %s", userID, snap.UserID)
	}
	if !snap.CreatedAt.Equal(createdAt) {
		t.Fatalf("expected createdAt %v, got %v", createdAt, snap.CreatedAt)
	}
	if len(snap.Tracks) != len(topTracks) {
		t.Fatalf("expected %d tracks, got %d", len(topTracks), len(snap.Tracks))
	}
	for i, tr := range snap.Tracks {
		if tr.ID != trackIDs[i] {
			t.Fatalf("track %d: expected ID %s, got %s", i, trackIDs[i], tr.ID)
		}
	}

	if err := rawDB.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

// TimeRangePropagation: verifies ShortTerm and LongTerm are forwarded to spotify correctly.
func TestSnapshotService_CreateSnapshot_TimeRangePropagation(t *testing.T) {
	t.Parallel()

	for _, tr := range []spotify.TimeRange{spotify.ShortTerm, spotify.MediumTerm, spotify.LongTerm} {
		timeRange := tr
		t.Run(string(timeRange), func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			userID := uuid.New()
			snapID := uuid.New()

			rawDB, mock, txProvider := newDB(t)
			mock.ExpectBegin()
			mock.ExpectQuery(`INSERT INTO snapshots`).
				WithArgs(userID).
				WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(snapID, time.Now()))
			mock.ExpectCommit()
			mock.ExpectClose()

			var gotTimeRange spotify.TimeRange
			spotifySvc := &mockSpotifyService{
				getTopTracksFn: func(_ context.Context, _ uuid.UUID, _ int, tr spotify.TimeRange) ([]track.Track, error) {
					gotTimeRange = tr
					return []track.Track{}, nil
				},
			}

			svc := NewService(nil, spotifySvc, txProvider, &mockCache{}, slog.Default())
			_, err := svc.CreateSnapshot(ctx, userID, timeRange)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotTimeRange != timeRange {
				t.Fatalf("expected timeRange %s forwarded to spotify, got %s", timeRange, gotTimeRange)
			}
			if err := rawDB.Close(); err != nil {
				t.Fatalf("close db: %v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet sqlmock expectations: %v", err)
			}
		})
	}
}

// EmptyTracks: Spotify returns empty list — snapshot is created, tracks and snapshot_tracks are not touched.
func TestSnapshotService_CreateSnapshot_EmptyTracks(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	snapID := uuid.New()
	createdAt := time.Now()

	rawDB, mock, txProvider := newDB(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO snapshots`).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "created_at"}).AddRow(snapID, createdAt),
		)
	mock.ExpectCommit()
	mock.ExpectClose()

	spotifySvc := &mockSpotifyService{
		getTopTracksFn: func(_ context.Context, _ uuid.UUID, _ int, _ spotify.TimeRange) ([]track.Track, error) {
			return []track.Track{}, nil
		},
	}

	svc := NewService(nil, spotifySvc, txProvider, &mockCache{}, slog.Default())
	snap, err := svc.CreateSnapshot(ctx, userID, spotify.MediumTerm)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if snap == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if snap.ID != snapID {
		t.Fatalf("expected snapID %s, got %s", snapID, snap.ID)
	}
	if len(snap.Tracks) != 0 {
		t.Fatalf("expected 0 tracks, got %d", len(snap.Tracks))
	}

	if err := rawDB.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

// SpotifyError: Spotify returns an error — transaction must not be started at all.
func TestSnapshotService_CreateSnapshot_SpotifyError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expectedErr := errors.New("spotify unavailable")

	rawDB, mock, txProvider := newDB(t)
	mock.ExpectClose()

	spotifySvc := &mockSpotifyService{
		getTopTracksFn: func(_ context.Context, _ uuid.UUID, _ int, _ spotify.TimeRange) ([]track.Track, error) {
			return nil, expectedErr
		},
	}

	svc := NewService(nil, spotifySvc, txProvider, &mockCache{}, slog.Default())
	snap, err := svc.CreateSnapshot(ctx, uuid.New(), spotify.MediumTerm)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "fetch top tracks from spotify") {
		t.Fatalf("expected 'fetch top tracks from spotify' in error, got: %v", err)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}
	if snap != nil {
		t.Fatal("expected nil snapshot on error")
	}

	if err := rawDB.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

// InsertSnapshotError: INSERT INTO snapshots fails — transaction is rolled back, tracks are not touched.
func TestSnapshotService_CreateSnapshot_InsertSnapshotError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expectedErr := errors.New("db write failed")
	topTracks := makeTopTracks(2)

	rawDB, mock, txProvider := newDB(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO snapshots`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(expectedErr)
	mock.ExpectRollback()
	mock.ExpectClose()

	spotifySvc := &mockSpotifyService{
		getTopTracksFn: func(_ context.Context, _ uuid.UUID, _ int, _ spotify.TimeRange) ([]track.Track, error) {
			return topTracks, nil
		},
	}

	svc := NewService(nil, spotifySvc, txProvider, &mockCache{}, slog.Default())
	snap, err := svc.CreateSnapshot(ctx, uuid.New(), spotify.MediumTerm)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "create snapshot") {
		t.Fatalf("expected 'create snapshot' in error, got: %v", err)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}
	if snap != nil {
		t.Fatal("expected nil snapshot on error")
	}

	if err := rawDB.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

// UpsertTrackError: INSERT INTO tracks fails on the first track — transaction is rolled back, snapshot_tracks is not touched.
func TestSnapshotService_CreateSnapshot_UpsertTrackError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expectedErr := errors.New("track upsert failed")
	topTracks := makeTopTracks(3)
	snapID := uuid.New()

	rawDB, mock, txProvider := newDB(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO snapshots`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "created_at"}).AddRow(snapID, time.Now()),
		)
	mock.ExpectQuery(`INSERT INTO tracks`).
		WithArgs(topTracks[0].SpotifyID, topTracks[0].Name, topTracks[0].Popularity).
		WillReturnError(expectedErr)
	mock.ExpectRollback()
	mock.ExpectClose()

	spotifySvc := &mockSpotifyService{
		getTopTracksFn: func(_ context.Context, _ uuid.UUID, _ int, _ spotify.TimeRange) ([]track.Track, error) {
			return topTracks, nil
		},
	}

	svc := NewService(nil, spotifySvc, txProvider, &mockCache{}, slog.Default())
	snap, err := svc.CreateSnapshot(ctx, uuid.New(), spotify.MediumTerm)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), fmt.Sprintf("upsert track %q", topTracks[0].SpotifyID)) {
		t.Fatalf("expected upsert track error, got: %v", err)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}
	if snap != nil {
		t.Fatal("expected nil snapshot on error")
	}

	if err := rawDB.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

// LinkTrackError: track upsert succeeds but INSERT INTO snapshot_tracks fails — transaction is rolled back.
func TestSnapshotService_CreateSnapshot_LinkTrackError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expectedErr := errors.New("link failed")
	topTracks := makeTopTracks(2)
	snapID := uuid.New()
	trackID := uuid.New()

	rawDB, mock, txProvider := newDB(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO snapshots`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "created_at"}).AddRow(snapID, time.Now()),
		)
	mock.ExpectQuery(`INSERT INTO tracks`).
		WithArgs(topTracks[0].SpotifyID, topTracks[0].Name, topTracks[0].Popularity).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(trackID))
	mock.ExpectExec(`INSERT INTO snapshot_tracks`).
		WithArgs(snapID, trackID, 1).
		WillReturnError(expectedErr)
	mock.ExpectRollback()
	mock.ExpectClose()

	spotifySvc := &mockSpotifyService{
		getTopTracksFn: func(_ context.Context, _ uuid.UUID, _ int, _ spotify.TimeRange) ([]track.Track, error) {
			return topTracks, nil
		},
	}

	svc := NewService(nil, spotifySvc, txProvider, &mockCache{}, slog.Default())
	snap, err := svc.CreateSnapshot(ctx, uuid.New(), spotify.MediumTerm)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), fmt.Sprintf("link track %q to snapshot", topTracks[0].SpotifyID)) {
		t.Fatalf("expected link track error, got: %v", err)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}
	if snap != nil {
		t.Fatal("expected nil snapshot on error")
	}

	if err := rawDB.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
