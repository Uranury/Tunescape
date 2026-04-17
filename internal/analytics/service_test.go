package analytics

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/cache"
	"gitlab.com/Uranury/tunescape/internal/reccobeats"
	"gitlab.com/Uranury/tunescape/internal/snapshot"
	"gitlab.com/Uranury/tunescape/internal/track"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"gitlab.com/Uranury/tunescape/pkg/config"
	"gitlab.com/Uranury/tunescape/pkg/database"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

type mockAnalyticsRepo struct {
	getLatestSnapshotFn       func(ctx context.Context, userID uuid.UUID) (*snapshot.Snapshot, error)
	getTracksBySnapshotFn     func(ctx context.Context, snapshotID uuid.UUID) ([]track.Track, error)
	bulkUpsertAudioFeaturesFn func(ctx context.Context, features []TrackAudioFeatures) error
	getAveragesFn             func(ctx context.Context, snapshotID uuid.UUID) (*AudioFeatureAverages, int, error)
}

func (m *mockAnalyticsRepo) GetLatestSnapshotByUserID(ctx context.Context, userID uuid.UUID) (*snapshot.Snapshot, error) {
	return m.getLatestSnapshotFn(ctx, userID)
}
func (m *mockAnalyticsRepo) GetTracksBySnapshotID(ctx context.Context, snapshotID uuid.UUID) ([]track.Track, error) {
	return m.getTracksBySnapshotFn(ctx, snapshotID)
}
func (m *mockAnalyticsRepo) BulkUpsertAudioFeatures(ctx context.Context, features []TrackAudioFeatures) error {
	return m.bulkUpsertAudioFeaturesFn(ctx, features)
}
func (m *mockAnalyticsRepo) GetAveragesBySnapshotID(ctx context.Context, snapshotID uuid.UUID) (*AudioFeatureAverages, int, error) {
	return m.getAveragesFn(ctx, snapshotID)
}

type mockCache struct {
	getFn    func(ctx context.Context, key string) ([]byte, error)
	setFn    func(ctx context.Context, key string, data []byte, ttl time.Duration) error
	deleteFn func(ctx context.Context, key string) error
	existsFn func(ctx context.Context, key string) (bool, error)
}

func (m *mockCache) Get(ctx context.Context, key string) ([]byte, error) {
	return m.getFn(ctx, key)
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

func noopCache() cache.Cache {
	return &mockCache{
		getFn: func(_ context.Context, _ string) ([]byte, error) {
			return nil, errors.New("cache miss")
		},
	}
}

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func newDB(t *testing.T) (sqlmock.Sqlmock, database.TxProvider) {
	t.Helper()
	rawDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	db := sqlx.NewDb(rawDB, "sqlmock")
	t.Cleanup(func() { _ = rawDB.Close() })
	return mock, database.NewTxProvider(db)
}

func makeTracks(n int) []track.Track {
	tracks := make([]track.Track, n)
	for i := range tracks {
		tracks[i] = track.Track{
			ID:        uuid.New(),
			SpotifyID: fmt.Sprintf("spotify-%d", i+1),
			Name:      fmt.Sprintf("Track %d", i+1),
		}
	}
	return tracks
}

func makeAudioFeaturesResponse(tracks []track.Track) []reccobeats.AudioFeatures {
	features := make([]reccobeats.AudioFeatures, len(tracks))
	for i, t := range tracks {
		features[i] = reccobeats.AudioFeatures{
			Href:         fmt.Sprintf("https://open.spotify.com/track/%s", t.SpotifyID),
			Danceability: 0.5 + float64(i)*0.01,
			Energy:       0.6 + float64(i)*0.01,
			Valence:      0.7 + float64(i)*0.01,
		}
	}
	return features
}

func newReccobeatsService(transport http.RoundTripper) reccobeats.Service {
	client := reccobeats.NewClient(
		config.Reccobeats{BaseURL: "http://reccobeats.test"},
		&http.Client{Transport: transport},
	)
	return reccobeats.NewService(client)
}

func reccobeatsOKTransport(t *testing.T, features []reccobeats.AudioFeatures) http.RoundTripper {
	t.Helper()
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		body, err := json.Marshal(map[string]any{"content": features})
		if err != nil {
			t.Fatalf("marshal reccobeats response: %v", err)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Request:    req,
		}, nil
	})
}

func defaultAverages() *AudioFeatureAverages {
	return &AudioFeatureAverages{
		Danceability: 0.55,
		Energy:       0.65,
		Valence:      0.75,
		Tempo:        120.0,
		Loudness:     -5.0,
	}
}

func anyArgs(n int) []driver.Value {
	args := make([]driver.Value, n)
	for i := range args {
		args[i] = sqlmock.AnyArg()
	}
	return args
}

// Success: snapshot found, tracks fetched, audio features upserted, averages returned.
func TestAnalyticsService_GetMusicTaste_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	snapID := uuid.New()
	tracks := makeTracks(3)
	features := makeAudioFeaturesResponse(tracks)
	avgs := defaultAverages()

	mock, txProvider := newDB(t)
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO track_audio_features`).
		WithArgs(anyArgs(len(tracks) * 10)...).
		WillReturnResult(sqlmock.NewResult(0, int64(len(tracks))))
	mock.ExpectCommit()

	repo := &mockAnalyticsRepo{
		getLatestSnapshotFn: func(_ context.Context, gotUserID uuid.UUID) (*snapshot.Snapshot, error) {
			if gotUserID != userID {
				t.Fatalf("expected userID %s, got %s", userID, gotUserID)
			}
			return &snapshot.Snapshot{ID: snapID, UserID: userID, CreatedAt: time.Now()}, nil
		},
		getTracksBySnapshotFn: func(_ context.Context, gotSnapshotID uuid.UUID) ([]track.Track, error) {
			if gotSnapshotID != snapID {
				t.Fatalf("expected snapshotID %s, got %s", snapID, gotSnapshotID)
			}
			return tracks, nil
		},
		getAveragesFn: func(_ context.Context, gotSnapshotID uuid.UUID) (*AudioFeatureAverages, int, error) {
			if gotSnapshotID != snapID {
				t.Fatalf("expected snapshotID %s, got %s", snapID, gotSnapshotID)
			}
			return avgs, len(tracks), nil
		},
	}

	svc := NewService(repo, newReccobeatsService(reccobeatsOKTransport(t, features)), txProvider, slog.Default(), noopCache())
	resp, err := svc.GetMusicTaste(ctx, userID)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.SnapshotID != snapID {
		t.Fatalf("expected snapshotID %s, got %s", snapID, resp.SnapshotID)
	}
	if resp.TracksCount != len(tracks) {
		t.Fatalf("expected tracksCount %d, got %d", len(tracks), resp.TracksCount)
	}
	if resp.Averages != *avgs {
		t.Fatalf("expected averages %+v, got %+v", *avgs, resp.Averages)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

// NoSnapshot: GetLatestSnapshotByUserID returns ErrNoSnapshot — propagated as-is.
func TestAnalyticsService_GetMusicTaste_NoSnapshot(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	_, txProvider := newDB(t)

	repo := &mockAnalyticsRepo{
		getLatestSnapshotFn: func(_ context.Context, _ uuid.UUID) (*snapshot.Snapshot, error) {
			return nil, apperrors.ErrNoSnapshot
		},
		getTracksBySnapshotFn: func(_ context.Context, _ uuid.UUID) ([]track.Track, error) {
			t.Fatal("GetTracksBySnapshotID must not be called when snapshot is missing")
			return nil, nil
		},
	}

	svc := NewService(repo, nil, txProvider, slog.Default(), noopCache())
	resp, err := svc.GetMusicTaste(ctx, uuid.New())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, apperrors.ErrNoSnapshot) {
		t.Fatalf("expected ErrNoSnapshot, got %v", err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
}

// GetSnapshotError: GetLatestSnapshotByUserID returns an unexpected error — wrapped with "get latest snapshot".
func TestAnalyticsService_GetMusicTaste_GetSnapshotError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expectedErr := errors.New("db down")
	_, txProvider := newDB(t)

	repo := &mockAnalyticsRepo{
		getLatestSnapshotFn: func(_ context.Context, _ uuid.UUID) (*snapshot.Snapshot, error) {
			return nil, expectedErr
		},
		getTracksBySnapshotFn: func(_ context.Context, _ uuid.UUID) ([]track.Track, error) {
			t.Fatal("GetTracksBySnapshotID must not be called on snapshot error")
			return nil, nil
		},
	}

	svc := NewService(repo, nil, txProvider, slog.Default(), noopCache())
	resp, err := svc.GetMusicTaste(ctx, uuid.New())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "get latest snapshot") {
		t.Fatalf("expected 'get latest snapshot' in error, got: %v", err)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
}

// GetTracksError: GetTracksBySnapshotID fails — transaction must not be started.
func TestAnalyticsService_GetMusicTaste_GetTracksError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expectedErr := errors.New("tracks query failed")
	snapID := uuid.New()
	mock, txProvider := newDB(t)

	repo := &mockAnalyticsRepo{
		getLatestSnapshotFn: func(_ context.Context, _ uuid.UUID) (*snapshot.Snapshot, error) {
			return &snapshot.Snapshot{ID: snapID, CreatedAt: time.Now()}, nil
		},
		getTracksBySnapshotFn: func(_ context.Context, _ uuid.UUID) ([]track.Track, error) {
			return nil, expectedErr
		},
		bulkUpsertAudioFeaturesFn: func(_ context.Context, _ []TrackAudioFeatures) error {
			t.Fatal("BulkUpsertAudioFeatures must not be called when GetTracks fails")
			return nil
		},
	}

	svc := NewService(repo, nil, txProvider, slog.Default(), noopCache())
	resp, err := svc.GetMusicTaste(ctx, uuid.New())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "get snapshot tracks") {
		t.Fatalf("expected 'get snapshot tracks' in error, got: %v", err)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

// ReccobeatsError: GetAudioFeaturesBatch fails — transaction is rolled back, GetAverages not called.
func TestAnalyticsService_GetMusicTaste_ReccobeatsError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	snapID := uuid.New()
	tracks := makeTracks(2)

	mock, txProvider := newDB(t)
	mock.ExpectBegin()
	mock.ExpectRollback()

	repo := &mockAnalyticsRepo{
		getLatestSnapshotFn: func(_ context.Context, _ uuid.UUID) (*snapshot.Snapshot, error) {
			return &snapshot.Snapshot{ID: snapID, CreatedAt: time.Now()}, nil
		},
		getTracksBySnapshotFn: func(_ context.Context, _ uuid.UUID) ([]track.Track, error) {
			return tracks, nil
		},
		bulkUpsertAudioFeaturesFn: func(_ context.Context, _ []TrackAudioFeatures) error {
			t.Fatal("BulkUpsertAudioFeatures must not be called when Reccobeats fails")
			return nil
		},
		getAveragesFn: func(_ context.Context, _ uuid.UUID) (*AudioFeatureAverages, int, error) {
			t.Fatal("GetAverages must not be called when Reccobeats fails")
			return nil, 0, nil
		},
	}

	reccobeatsService := newReccobeatsService(roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusServiceUnavailable,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	}))

	svc := NewService(repo, reccobeatsService, txProvider, slog.Default(), noopCache())
	resp, err := svc.GetMusicTaste(ctx, uuid.New())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "fetch audio features batch") {
		t.Fatalf("expected 'fetch audio features batch' in error, got: %v", err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

// UpsertAudioFeaturesError: BulkUpsertAudioFeatures fails — transaction is rolled back, GetAverages not called.
func TestAnalyticsService_GetMusicTaste_UpsertAudioFeaturesError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expectedErr := errors.New("upsert failed")
	snapID := uuid.New()
	tracks := makeTracks(2)
	features := makeAudioFeaturesResponse(tracks)

	mock, txProvider := newDB(t)
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO track_audio_features`).
		WithArgs(anyArgs(len(tracks) * 10)...).
		WillReturnError(expectedErr)
	mock.ExpectRollback()

	repo := &mockAnalyticsRepo{
		getLatestSnapshotFn: func(_ context.Context, _ uuid.UUID) (*snapshot.Snapshot, error) {
			return &snapshot.Snapshot{ID: snapID, CreatedAt: time.Now()}, nil
		},
		getTracksBySnapshotFn: func(_ context.Context, _ uuid.UUID) ([]track.Track, error) {
			return tracks, nil
		},
		getAveragesFn: func(_ context.Context, _ uuid.UUID) (*AudioFeatureAverages, int, error) {
			t.Fatal("GetAverages must not be called when upsert fails")
			return nil, 0, nil
		},
	}

	svc := NewService(repo, newReccobeatsService(reccobeatsOKTransport(t, features)), txProvider, slog.Default(), noopCache())
	resp, err := svc.GetMusicTaste(ctx, uuid.New())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "bulk upsert audio features batch") {
		t.Fatalf("expected 'bulk upsert audio features batch' in error, got: %v", err)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

// GetAveragesError: GetAveragesBySnapshotID fails after a successful transaction.
func TestAnalyticsService_GetMusicTaste_GetAveragesError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expectedErr := errors.New("aggregation failed")
	snapID := uuid.New()
	tracks := makeTracks(2)
	features := makeAudioFeaturesResponse(tracks)

	mock, txProvider := newDB(t)
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO track_audio_features`).
		WithArgs(anyArgs(len(tracks) * 10)...).
		WillReturnResult(sqlmock.NewResult(0, int64(len(tracks))))
	mock.ExpectCommit()

	repo := &mockAnalyticsRepo{
		getLatestSnapshotFn: func(_ context.Context, _ uuid.UUID) (*snapshot.Snapshot, error) {
			return &snapshot.Snapshot{ID: snapID, CreatedAt: time.Now()}, nil
		},
		getTracksBySnapshotFn: func(_ context.Context, _ uuid.UUID) ([]track.Track, error) {
			return tracks, nil
		},
		getAveragesFn: func(_ context.Context, _ uuid.UUID) (*AudioFeatureAverages, int, error) {
			return nil, 0, expectedErr
		},
	}

	svc := NewService(repo, newReccobeatsService(reccobeatsOKTransport(t, features)), txProvider, slog.Default(), noopCache())
	resp, err := svc.GetMusicTaste(ctx, uuid.New())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "aggregate audio features") {
		t.Fatalf("expected 'aggregate audio features' in error, got: %v", err)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

// BatchPagination: tracks count exceeds audioFeaturesBatchSize — two HTTP requests are made (40 + 5).
func TestAnalyticsService_GetMusicTaste_BatchPagination(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	snapID := uuid.New()
	tracks := makeTracks(audioFeaturesBatchSize + 5)
	avgs := defaultAverages()

	mock, txProvider := newDB(t)
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO track_audio_features`).
		WithArgs(anyArgs(audioFeaturesBatchSize * 10)...).
		WillReturnResult(sqlmock.NewResult(0, audioFeaturesBatchSize))
	mock.ExpectExec(`INSERT INTO track_audio_features`).
		WithArgs(anyArgs(5 * 10)...).
		WillReturnResult(sqlmock.NewResult(0, 5))
	mock.ExpectCommit()

	var httpCallCount int
	transport := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		start := httpCallCount * audioFeaturesBatchSize
		httpCallCount++

		batch := tracks[start:]
		if len(batch) > audioFeaturesBatchSize {
			batch = batch[:audioFeaturesBatchSize]
		}

		features := makeAudioFeaturesResponse(batch)
		body, _ := json.Marshal(map[string]any{"content": features})
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Request:    req,
		}, nil
	})

	repo := &mockAnalyticsRepo{
		getLatestSnapshotFn: func(_ context.Context, _ uuid.UUID) (*snapshot.Snapshot, error) {
			return &snapshot.Snapshot{ID: snapID, CreatedAt: time.Now()}, nil
		},
		getTracksBySnapshotFn: func(_ context.Context, _ uuid.UUID) ([]track.Track, error) {
			return tracks, nil
		},
		getAveragesFn: func(_ context.Context, _ uuid.UUID) (*AudioFeatureAverages, int, error) {
			return avgs, len(tracks), nil
		},
	}

	svc := NewService(repo, newReccobeatsService(transport), txProvider, slog.Default(), noopCache())
	resp, err := svc.GetMusicTaste(ctx, uuid.New())

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if httpCallCount != 2 {
		t.Fatalf("expected 2 HTTP calls for %d tracks (batchSize=%d), got %d",
			len(tracks), audioFeaturesBatchSize, httpCallCount)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

// UnknownSpotifyID: Reccobeats returns a feature with an href that does not match any track — silently skipped.
func TestAnalyticsService_GetMusicTaste_UnknownSpotifyID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	snapID := uuid.New()
	tracks := makeTracks(2)
	avgs := defaultAverages()

	mock, txProvider := newDB(t)
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO track_audio_features`).
		WithArgs(anyArgs(len(tracks) * 10)...).
		WillReturnResult(sqlmock.NewResult(0, int64(len(tracks))))
	mock.ExpectCommit()

	features := makeAudioFeaturesResponse(tracks)
	features = append(features, reccobeats.AudioFeatures{
		Href: "https://open.spotify.com/track/unknown-id",
	})

	repo := &mockAnalyticsRepo{
		getLatestSnapshotFn: func(_ context.Context, _ uuid.UUID) (*snapshot.Snapshot, error) {
			return &snapshot.Snapshot{ID: snapID, CreatedAt: time.Now()}, nil
		},
		getTracksBySnapshotFn: func(_ context.Context, _ uuid.UUID) ([]track.Track, error) {
			return tracks, nil
		},
		getAveragesFn: func(_ context.Context, _ uuid.UUID) (*AudioFeatureAverages, int, error) {
			return avgs, len(tracks), nil
		},
	}

	svc := NewService(repo, newReccobeatsService(reccobeatsOKTransport(t, features)), txProvider, slog.Default(), noopCache())
	resp, err := svc.GetMusicTaste(ctx, uuid.New())

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

// CacheHit: valid cached response is returned immediately — no DB or ReccoBeats calls made.
func TestAnalyticsService_GetMusicTaste_CacheHit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	snapID := uuid.New()
	avgs := defaultAverages()

	expected := MusicTasteResponse{
		SnapshotID:  snapID,
		TracksCount: 3,
		Averages:    *avgs,
	}
	cachedBytes, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("marshal cached response: %v", err)
	}

	c := &mockCache{
		getFn: func(_ context.Context, key string) ([]byte, error) {
			if key != "music_taste:"+userID.String() {
				t.Fatalf("unexpected cache key: %s", key)
			}
			return cachedBytes, nil
		},
	}

	repo := &mockAnalyticsRepo{
		getLatestSnapshotFn: func(_ context.Context, _ uuid.UUID) (*snapshot.Snapshot, error) {
			t.Fatal("GetLatestSnapshotByUserID must not be called on cache hit")
			return nil, nil
		},
		getTracksBySnapshotFn: func(_ context.Context, _ uuid.UUID) ([]track.Track, error) {
			t.Fatal("GetTracksBySnapshotID must not be called on cache hit")
			return nil, nil
		},
	}

	_, txProvider := newDB(t)
	svc := NewService(repo, nil, txProvider, slog.Default(), c)
	resp, err := svc.GetMusicTaste(ctx, userID)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.SnapshotID != snapID {
		t.Fatalf("expected snapshotID %s, got %s", snapID, resp.SnapshotID)
	}
	if resp.TracksCount != expected.TracksCount {
		t.Fatalf("expected tracksCount %d, got %d", expected.TracksCount, resp.TracksCount)
	}
	if resp.Averages != *avgs {
		t.Fatalf("expected averages %+v, got %+v", *avgs, resp.Averages)
	}
}
