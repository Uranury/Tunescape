package friends

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gitlab.com/Uranury/tunescape/internal/analytics"
	"gitlab.com/Uranury/tunescape/internal/playlist"
	"gitlab.com/Uranury/tunescape/internal/snapshot"
	"gitlab.com/Uranury/tunescape/internal/spotify"
	"gitlab.com/Uranury/tunescape/internal/track"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

type mockFriendsRepo struct {
	sendRequestFn         func(ctx context.Context, senderID, receiverID uuid.UUID) error
	getRequestFn          func(ctx context.Context, requestID int64) (*FriendRequest, error)
	acceptRequestFn       func(ctx context.Context, requestID int64, senderID, receiverID uuid.UUID) error
	rejectRequestFn       func(ctx context.Context, requestID int64, receiverID uuid.UUID) error
	listIncomingFn        func(ctx context.Context, userID uuid.UUID) ([]IncomingRequest, error)
	listFriendsFn         func(ctx context.Context, userID uuid.UUID) ([]FriendProfile, error)
	areFriendsFn          func(ctx context.Context, userID, friendID uuid.UUID) (bool, error)
	removeFriendFn        func(ctx context.Context, userID, friendID uuid.UUID) error
	hasSpotifyConnectedFn func(ctx context.Context, userID uuid.UUID) (bool, error)
}

func (m *mockFriendsRepo) SendRequest(ctx context.Context, senderID, receiverID uuid.UUID) error {
	return m.sendRequestFn(ctx, senderID, receiverID)
}
func (m *mockFriendsRepo) GetRequest(ctx context.Context, requestID int64) (*FriendRequest, error) {
	return m.getRequestFn(ctx, requestID)
}
func (m *mockFriendsRepo) AcceptRequest(ctx context.Context, requestID int64, senderID, receiverID uuid.UUID) error {
	return m.acceptRequestFn(ctx, requestID, senderID, receiverID)
}
func (m *mockFriendsRepo) RejectRequest(ctx context.Context, requestID int64, receiverID uuid.UUID) error {
	return m.rejectRequestFn(ctx, requestID, receiverID)
}
func (m *mockFriendsRepo) ListIncoming(ctx context.Context, userID uuid.UUID) ([]IncomingRequest, error) {
	return m.listIncomingFn(ctx, userID)
}
func (m *mockFriendsRepo) ListFriends(ctx context.Context, userID uuid.UUID) ([]FriendProfile, error) {
	return m.listFriendsFn(ctx, userID)
}
func (m *mockFriendsRepo) AreFriends(ctx context.Context, userID, friendID uuid.UUID) (bool, error) {
	return m.areFriendsFn(ctx, userID, friendID)
}
func (m *mockFriendsRepo) RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	return m.removeFriendFn(ctx, userID, friendID)
}
func (m *mockFriendsRepo) HasSpotifyConnected(ctx context.Context, userID uuid.UUID) (bool, error) {
	return m.hasSpotifyConnectedFn(ctx, userID)
}

type mockTxProvider struct {
	runInTxFn func(ctx context.Context, fn func(database.Executor) error) error
}

func (m *mockTxProvider) RunInTx(ctx context.Context, fn func(database.Executor) error) error {
	return m.runInTxFn(ctx, fn)
}

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

type mockAnalyticsRepo struct {
	getAveragesFn func(ctx context.Context, snapshotID uuid.UUID) (*analytics.AudioFeatureAverages, int, error)
}

func (m *mockAnalyticsRepo) GetLatestSnapshotByUserID(_ context.Context, _ uuid.UUID) (*snapshot.Snapshot, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) GetTracksBySnapshotID(_ context.Context, _ uuid.UUID) ([]track.Track, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) BulkUpsertAudioFeatures(_ context.Context, _ []analytics.TrackAudioFeatures) error {
	return nil
}
func (m *mockAnalyticsRepo) GetAveragesBySnapshotID(ctx context.Context, snapshotID uuid.UUID) (*analytics.AudioFeatureAverages, int, error) {
	return m.getAveragesFn(ctx, snapshotID)
}

type mockPlaylistService struct {
	listByUserIDFn func(ctx context.Context, userID uuid.UUID) ([]playlist.Playlist, error)
}

func (m *mockPlaylistService) CreateFromLatestSnapshot(_ context.Context, _ uuid.UUID) (*playlist.Response, error) {
	return nil, nil
}
func (m *mockPlaylistService) ListByUserID(ctx context.Context, userID uuid.UUID) ([]playlist.Playlist, error) {
	return m.listByUserIDFn(ctx, userID)
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func newFriendsService(
	repo Repository,
	snapshotSvc snapshot.Service,
	analyticsRepo analytics.Repository,
	playlistSvc playlist.Service,
	txProvider database.TxProvider,
) Service {
	return NewService(repo, snapshotSvc, analyticsRepo, playlistSvc, txProvider, testLogger())
}

func TestCompatibilityScore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		a, b TasteScores
		want float64
	}{
		{
			name: "identical tastes",
			a:    TasteScores{Valence: 0.5, Energy: 0.5, Danceability: 0.5, Acousticness: 0.5},
			b:    TasteScores{Valence: 0.5, Energy: 0.5, Danceability: 0.5, Acousticness: 0.5},
			want: 100,
		},
		{
			name: "maximum distance",
			a:    TasteScores{Valence: 0, Energy: 0, Danceability: 0, Acousticness: 0},
			b:    TasteScores{Valence: 1, Energy: 1, Danceability: 1, Acousticness: 1},
			want: 0,
		},
		{
			name: "half distance",
			a:    TasteScores{Valence: 0, Energy: 0, Danceability: 0, Acousticness: 0},
			b:    TasteScores{Valence: 0.5, Energy: 0.5, Danceability: 0.5, Acousticness: 0.5},
			want: 50,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := compatibilityScore(tc.a, tc.b); got != tc.want {
				t.Fatalf("compatibilityScore() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFriendsService_SendRequest_CannotAddSelf(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	svc := newFriendsService(
		&mockFriendsRepo{},
		&mockSnapshotService{},
		&mockAnalyticsRepo{},
		&mockPlaylistService{},
		&mockTxProvider{},
	)

	err := svc.SendRequest(context.Background(), id, id)
	if !errors.Is(err, apperrors.ErrCannotAddSelf) {
		t.Fatalf("expected ErrCannotAddSelf, got %v", err)
	}
}

func TestFriendsService_SendRequest_AlreadyFriends(t *testing.T) {
	t.Parallel()

	senderID := uuid.New()
	receiverID := uuid.New()

	repo := &mockFriendsRepo{
		areFriendsFn: func(_ context.Context, uid, fid uuid.UUID) (bool, error) {
			if uid != senderID || fid != receiverID {
				t.Fatalf("unexpected AreFriends call: %s, %s", uid, fid)
			}
			return true, nil
		},
		sendRequestFn: func(context.Context, uuid.UUID, uuid.UUID) error {
			t.Fatal("SendRequest must not be called when already friends")
			return nil
		},
	}

	svc := newFriendsService(repo, &mockSnapshotService{}, &mockAnalyticsRepo{}, &mockPlaylistService{}, &mockTxProvider{})
	err := svc.SendRequest(context.Background(), senderID, receiverID)
	if !errors.Is(err, apperrors.ErrAlreadyFriends) {
		t.Fatalf("expected ErrAlreadyFriends, got %v", err)
	}
}

func TestFriendsService_SendRequest_Success(t *testing.T) {
	t.Parallel()

	senderID := uuid.New()
	receiverID := uuid.New()
	var sent bool

	repo := &mockFriendsRepo{
		areFriendsFn: func(context.Context, uuid.UUID, uuid.UUID) (bool, error) { return false, nil },
		sendRequestFn: func(_ context.Context, s, r uuid.UUID) error {
			if s != senderID || r != receiverID {
				t.Fatalf("unexpected SendRequest: %s -> %s", s, r)
			}
			sent = true
			return nil
		},
	}

	svc := newFriendsService(repo, &mockSnapshotService{}, &mockAnalyticsRepo{}, &mockPlaylistService{}, &mockTxProvider{})
	if err := svc.SendRequest(context.Background(), senderID, receiverID); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !sent {
		t.Fatal("expected SendRequest to be called on repository")
	}
}

func TestFriendsService_AcceptRequest_Forbidden(t *testing.T) {
	t.Parallel()

	receiverID := uuid.New()
	otherReceiver := uuid.New()

	repo := &mockFriendsRepo{
		getRequestFn: func(_ context.Context, _ int64) (*FriendRequest, error) {
			return &FriendRequest{
				ID:         1,
				SenderID:   uuid.New(),
				ReceiverID: receiverID,
				Status:     "pending",
			}, nil
		},
	}

	svc := newFriendsService(repo, &mockSnapshotService{}, &mockAnalyticsRepo{}, &mockPlaylistService{}, &mockTxProvider{})
	err := svc.AcceptRequest(context.Background(), 1, otherReceiver)
	if !errors.Is(err, apperrors.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestFriendsService_AcceptRequest_NotPending(t *testing.T) {
	t.Parallel()

	receiverID := uuid.New()

	repo := &mockFriendsRepo{
		getRequestFn: func(_ context.Context, _ int64) (*FriendRequest, error) {
			return &FriendRequest{
				ID:         1,
				SenderID:   uuid.New(),
				ReceiverID: receiverID,
				Status:     "accepted",
			}, nil
		},
	}

	svc := newFriendsService(repo, &mockSnapshotService{}, &mockAnalyticsRepo{}, &mockPlaylistService{}, &mockTxProvider{})
	err := svc.AcceptRequest(context.Background(), 1, receiverID)
	if !errors.Is(err, apperrors.ErrRequestNotFound) {
		t.Fatalf("expected ErrRequestNotFound, got %v", err)
	}
}

func TestFriendsService_AcceptRequest_Success(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	senderID := uuid.New()
	receiverID := uuid.New()
	requestID := int64(42)

	rawDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	db := sqlx.NewDb(rawDB, "sqlmock")
	txProvider := database.NewTxProvider(db)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE friend_requests`).
		WithArgs(requestID, receiverID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`INSERT INTO friends`).
		WithArgs(senderID, receiverID).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit()
	mock.ExpectClose()

	repo := &mockFriendsRepo{
		getRequestFn: func(_ context.Context, id int64) (*FriendRequest, error) {
			if id != requestID {
				t.Fatalf("expected requestID %d, got %d", requestID, id)
			}
			return &FriendRequest{
				ID:         requestID,
				SenderID:   senderID,
				ReceiverID: receiverID,
				Status:     "pending",
			}, nil
		},
	}

	svc := newFriendsService(repo, &mockSnapshotService{}, &mockAnalyticsRepo{}, &mockPlaylistService{}, txProvider)
	if err := svc.AcceptRequest(ctx, requestID, receiverID); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if err := rawDB.Close(); err != nil { // ← явный Close до проверки
		t.Errorf("failed to close db: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestFriendsService_CompareTastes_Success(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	friendID := uuid.New()
	userSnapID := uuid.New()
	friendSnapID := uuid.New()

	repo := &mockFriendsRepo{
		areFriendsFn: func(context.Context, uuid.UUID, uuid.UUID) (bool, error) { return true, nil },
	}

	snapshotSvc := &mockSnapshotService{
		getLatestSnapshotFn: func(_ context.Context, id uuid.UUID) (*snapshot.Snapshot, error) {
			switch id {
			case userID:
				return &snapshot.Snapshot{ID: userSnapID, UserID: userID}, nil
			case friendID:
				return &snapshot.Snapshot{ID: friendSnapID, UserID: friendID}, nil
			default:
				t.Fatalf("unexpected user id: %s", id)
				return nil, nil
			}
		},
	}

	analyticsRepo := &mockAnalyticsRepo{
		getAveragesFn: func(_ context.Context, snapID uuid.UUID) (*analytics.AudioFeatureAverages, int, error) {
			switch snapID {
			case userSnapID:
				return &analytics.AudioFeatureAverages{
					Valence: 0.8, Energy: 0.6, Danceability: 0.7, Acousticness: 0.2,
				}, 10, nil
			case friendSnapID:
				return &analytics.AudioFeatureAverages{
					Valence: 0.8, Energy: 0.6, Danceability: 0.7, Acousticness: 0.2,
				}, 8, nil
			default:
				t.Fatalf("unexpected snapshot id: %s", snapID)
				return nil, 0, nil
			}
		},
	}

	svc := newFriendsService(repo, snapshotSvc, analyticsRepo, &mockPlaylistService{}, &mockTxProvider{})
	result, err := svc.CompareTastes(context.Background(), userID, friendID)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result.CompatibilityScore != 100 {
		t.Fatalf("expected compatibility score 100, got %v", result.CompatibilityScore)
	}
}

func TestFriendsService_CompareTastes_NoListeningData(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	friendID := uuid.New()
	snapID := uuid.New()

	repo := &mockFriendsRepo{
		areFriendsFn: func(context.Context, uuid.UUID, uuid.UUID) (bool, error) { return true, nil },
	}

	snapshotSvc := &mockSnapshotService{
		getLatestSnapshotFn: func(_ context.Context, id uuid.UUID) (*snapshot.Snapshot, error) {
			return &snapshot.Snapshot{ID: snapID, UserID: id}, nil
		},
	}

	analyticsRepo := &mockAnalyticsRepo{
		getAveragesFn: func(context.Context, uuid.UUID) (*analytics.AudioFeatureAverages, int, error) {
			return &analytics.AudioFeatureAverages{}, 0, nil
		},
	}

	svc := newFriendsService(repo, snapshotSvc, analyticsRepo, &mockPlaylistService{}, &mockTxProvider{})
	_, err := svc.CompareTastes(context.Background(), userID, friendID)
	if !errors.Is(err, apperrors.ErrNoSnapshot) {
		t.Fatalf("expected ErrNoSnapshot, got %v", err)
	}
}

func TestFriendsService_GetFriendPlaylists_NotFriends(t *testing.T) {
	t.Parallel()

	repo := &mockFriendsRepo{
		areFriendsFn: func(context.Context, uuid.UUID, uuid.UUID) (bool, error) { return false, nil },
	}

	svc := newFriendsService(repo, &mockSnapshotService{}, &mockAnalyticsRepo{}, &mockPlaylistService{}, &mockTxProvider{})
	_, err := svc.GetFriendPlaylists(context.Background(), uuid.New(), uuid.New())
	if !errors.Is(err, apperrors.ErrNotFriends) {
		t.Fatalf("expected ErrNotFriends, got %v", err)
	}
}

func TestFriendsService_GetFriendPlaylists_SpotifyNotConnected(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	friendID := uuid.New()

	repo := &mockFriendsRepo{
		areFriendsFn: func(context.Context, uuid.UUID, uuid.UUID) (bool, error) { return true, nil },
		hasSpotifyConnectedFn: func(_ context.Context, id uuid.UUID) (bool, error) {
			if id != friendID {
				t.Fatalf("expected friendID %s, got %s", friendID, id)
			}
			return false, nil
		},
	}

	playlistSvc := &mockPlaylistService{
		listByUserIDFn: func(context.Context, uuid.UUID) ([]playlist.Playlist, error) {
			t.Fatal("ListByUserID must not be called when Spotify is not connected")
			return nil, nil
		},
	}

	svc := newFriendsService(repo, &mockSnapshotService{}, &mockAnalyticsRepo{}, playlistSvc, &mockTxProvider{})
	_, err := svc.GetFriendPlaylists(context.Background(), userID, friendID)
	if !errors.Is(err, apperrors.ErrSpotifyNotConnected) {
		t.Fatalf("expected ErrSpotifyNotConnected, got %v", err)
	}
}

func TestFriendsService_GetFriendPlaylists_Success(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	friendID := uuid.New()
	createdAt := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	want := []playlist.Playlist{{
		SpotifyPlaylistID: "pl-1",
		Name:              "Top Tracks",
		ExternalURL:       "https://open.spotify.com/playlist/pl-1",
		EmbedURL:          "https://open.spotify.com/embed/playlist/pl-1",
		CreatedAt:         createdAt,
	}}

	repo := &mockFriendsRepo{
		areFriendsFn:          func(context.Context, uuid.UUID, uuid.UUID) (bool, error) { return true, nil },
		hasSpotifyConnectedFn: func(context.Context, uuid.UUID) (bool, error) { return true, nil },
	}

	playlistSvc := &mockPlaylistService{
		listByUserIDFn: func(_ context.Context, id uuid.UUID) ([]playlist.Playlist, error) {
			if id != friendID {
				t.Fatalf("expected friendID %s, got %s", friendID, id)
			}
			return want, nil
		},
	}

	svc := newFriendsService(repo, &mockSnapshotService{}, &mockAnalyticsRepo{}, playlistSvc, &mockTxProvider{})
	got, err := svc.GetFriendPlaylists(context.Background(), userID, friendID)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(got) != 1 || got[0].SpotifyPlaylistID != want[0].SpotifyPlaylistID {
		t.Fatalf("unexpected playlists: %+v", got)
	}
}
