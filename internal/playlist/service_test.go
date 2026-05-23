package playlist

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/snapshot"
	"gitlab.com/Uranury/tunescape/internal/spotify"
	"gitlab.com/Uranury/tunescape/internal/track"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
)

type mockSnapshotService struct {
	getLatestSnapshotFn func(ctx context.Context, userID uuid.UUID) (*snapshot.Snapshot, error)
}

func (m *mockSnapshotService) GetLatestSnapshot(ctx context.Context, userID uuid.UUID) (*snapshot.Snapshot, error) {
	return m.getLatestSnapshotFn(ctx, userID)
}
func (m *mockSnapshotService) CreateSnapshot(_ context.Context, _ uuid.UUID, _ spotify.TimeRange) (*snapshot.Snapshot, error) {
	return nil, nil
}
func (m *mockSnapshotService) ListSnapshots(_ context.Context, _ uuid.UUID) ([]snapshot.SnapshotSummary, error) {
	return nil, nil
}
func (m *mockSnapshotService) GetSnapshot(_ context.Context, _, _ uuid.UUID) (*snapshot.Snapshot, error) {
	return nil, nil
}

type mockSpotifyService struct {
	createPlaylistFn func(ctx context.Context, userID uuid.UUID, name string, trackURIs []string) (*spotify.PlaylistResult, error)
}

func (m *mockSpotifyService) CreatePlaylist(ctx context.Context, userID uuid.UUID, name string, trackURIs []string) (*spotify.PlaylistResult, error) {
	return m.createPlaylistFn(ctx, userID, name, trackURIs)
}
func (m *mockSpotifyService) AuthURL(_ string) string { return "" }
func (m *mockSpotifyService) ConnectAccount(_ context.Context, _ uuid.UUID, _ string) error {
	return nil
}
func (m *mockSpotifyService) Disconnect(_ context.Context, _ uuid.UUID) error { return nil }
func (m *mockSpotifyService) GetValidToken(_ context.Context, _ uuid.UUID) (string, error) {
	return "", nil
}
func (m *mockSpotifyService) GetTopTracks(_ context.Context, _ uuid.UUID, _ int, _ spotify.TimeRange) ([]track.Track, error) {
	return nil, nil
}
func (m *mockSpotifyService) UpsertTokens(_ context.Context, _ uuid.UUID, _, _ string, _ time.Time) error {
	return nil
}

type mockPlaylistRepo struct {
	insertFn func(ctx context.Context, userID uuid.UUID, p *Playlist) error
}

func (m *mockPlaylistRepo) Insert(ctx context.Context, userID uuid.UUID, p *Playlist) error {
	if m.insertFn != nil {
		return m.insertFn(ctx, userID, p)
	}
	p.CreatedAt = time.Now().UTC()
	return nil
}

func (m *mockPlaylistRepo) ListByUserID(_ context.Context, _ uuid.UUID) ([]Playlist, error) {
	return nil, nil
}

func makeSnapshot(trackCount int) *snapshot.Snapshot {
	tracks := make([]track.Track, trackCount)
	for i := range tracks {
		tracks[i] = track.Track{
			ID:        uuid.New(),
			SpotifyID: "spotify-id-" + uuid.New().String()[:8],
			Name:      "Track",
		}
	}
	return &snapshot.Snapshot{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		CreatedAt: time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC),
		Tracks:    tracks,
	}
}

func TestPlaylistService_CreateFromLatestSnapshot_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	snap := makeSnapshot(3)
	snap.UserID = userID

	playlistID := "pl-abc123"
	externalURL := "https://open.spotify.com/playlist/pl-abc123"

	snapshotSvc := &mockSnapshotService{
		getLatestSnapshotFn: func(_ context.Context, gotUserID uuid.UUID) (*snapshot.Snapshot, error) {
			if gotUserID != userID {
				t.Fatalf("expected userID %s, got %s", userID, gotUserID)
			}
			return snap, nil
		},
	}

	var gotName string
	var gotURIs []string
	spotifySvc := &mockSpotifyService{
		createPlaylistFn: func(_ context.Context, gotUserID uuid.UUID, name string, uris []string) (*spotify.PlaylistResult, error) {
			if gotUserID != userID {
				t.Fatalf("expected userID %s, got %s", userID, gotUserID)
			}
			gotName = name
			gotURIs = uris
			return &spotify.PlaylistResult{ID: playlistID, ExternalURL: externalURL}, nil
		},
	}

	svc := NewService(&mockPlaylistRepo{}, snapshotSvc, spotifySvc)
	resp, err := svc.CreateFromLatestSnapshot(ctx, userID)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.PlaylistID != playlistID {
		t.Fatalf("expected playlist_id %q, got %q", playlistID, resp.PlaylistID)
	}
	if resp.ExternalURL != externalURL {
		t.Fatalf("expected external_url %q, got %q", externalURL, resp.ExternalURL)
	}
	wantEmbed := "https://open.spotify.com/embed/playlist/" + playlistID + "?utm_source=generator&theme=0"
	if resp.EmbedURL != wantEmbed {
		t.Fatalf("expected embed_url %q, got %q", wantEmbed, resp.EmbedURL)
	}
	if !strings.Contains(gotName, "Mar 15, 2025") {
		t.Fatalf("expected playlist name to include snapshot date, got %q", gotName)
	}
	if len(gotURIs) != len(snap.Tracks) {
		t.Fatalf("expected %d track URIs, got %d", len(snap.Tracks), len(gotURIs))
	}
	for i, uri := range gotURIs {
		want := "spotify:track:" + snap.Tracks[i].SpotifyID
		if uri != want {
			t.Fatalf("track URI %d: expected %q, got %q", i, want, uri)
		}
	}
}

func TestPlaylistService_CreateFromLatestSnapshot_NoSnapshot(t *testing.T) {
	t.Parallel()

	snapshotSvc := &mockSnapshotService{
		getLatestSnapshotFn: func(_ context.Context, _ uuid.UUID) (*snapshot.Snapshot, error) {
			return nil, apperrors.ErrNoSnapshot
		},
	}
	spotifySvc := &mockSpotifyService{
		createPlaylistFn: func(_ context.Context, _ uuid.UUID, _ string, _ []string) (*spotify.PlaylistResult, error) {
			t.Fatal("CreatePlaylist must not be called when no snapshot exists")
			return nil, nil
		},
	}

	svc := NewService(&mockPlaylistRepo{}, snapshotSvc, spotifySvc)
	resp, err := svc.CreateFromLatestSnapshot(context.Background(), uuid.New())

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

func TestPlaylistService_CreateFromLatestSnapshot_SpotifyError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("spotify unavailable")
	snap := makeSnapshot(2)

	snapshotSvc := &mockSnapshotService{
		getLatestSnapshotFn: func(_ context.Context, _ uuid.UUID) (*snapshot.Snapshot, error) {
			return snap, nil
		},
	}
	spotifySvc := &mockSpotifyService{
		createPlaylistFn: func(_ context.Context, _ uuid.UUID, _ string, _ []string) (*spotify.PlaylistResult, error) {
			return nil, expectedErr
		},
	}

	svc := NewService(&mockPlaylistRepo{}, snapshotSvc, spotifySvc)
	resp, err := svc.CreateFromLatestSnapshot(context.Background(), uuid.New())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}
	if !strings.Contains(err.Error(), "create spotify playlist") {
		t.Fatalf("expected 'create spotify playlist' in error, got: %v", err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
}

func TestPlaylistService_CreateFromLatestSnapshot_UpstreamUnavailable(t *testing.T) {
	t.Parallel()

	snap := makeSnapshot(2)

	snapshotSvc := &mockSnapshotService{
		getLatestSnapshotFn: func(_ context.Context, _ uuid.UUID) (*snapshot.Snapshot, error) {
			return snap, nil
		},
	}
	spotifySvc := &mockSpotifyService{
		createPlaylistFn: func(_ context.Context, _ uuid.UUID, _ string, _ []string) (*spotify.PlaylistResult, error) {
			return nil, apperrors.ErrUpstreamUnavailable
		},
	}

	svc := NewService(&mockPlaylistRepo{}, snapshotSvc, spotifySvc)
	_, err := svc.CreateFromLatestSnapshot(context.Background(), uuid.New())

	if !errors.Is(err, apperrors.ErrUpstreamUnavailable) {
		t.Fatalf("expected ErrUpstreamUnavailable to be wrapped in error, got: %v", err)
	}
}
