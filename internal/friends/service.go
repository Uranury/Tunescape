package friends

import (
	"context"
	"log/slog"
	"math"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/analytics"
	"gitlab.com/Uranury/tunescape/internal/playlist"
	"gitlab.com/Uranury/tunescape/internal/snapshot"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

type Service interface {
	SendRequest(ctx context.Context, senderID, receiverID uuid.UUID) error
	AcceptRequest(ctx context.Context, requestID int64, receiverID uuid.UUID) error
	RejectRequest(ctx context.Context, requestID int64, receiverID uuid.UUID) error
	ListIncoming(ctx context.Context, userID uuid.UUID) ([]IncomingRequest, error)
	ListFriends(ctx context.Context, userID uuid.UUID) ([]FriendProfile, error)
	CompareTastes(ctx context.Context, userID, friendID uuid.UUID) (*TasteComparison, error)
	GetFriendPlaylists(ctx context.Context, userID, friendID uuid.UUID) ([]playlist.Playlist, error)
}

type service struct {
	repo          Repository
	snapshotSvc   snapshot.Service
	analyticsRepo analytics.Repository
	playlistSvc   playlist.Service
	txProvider    database.TxProvider
	logger        *slog.Logger
}

func NewService(
	repo Repository,
	snapshotSvc snapshot.Service,
	analyticsRepo analytics.Repository,
	playlistSvc playlist.Service,
	txProvider database.TxProvider,
	logger *slog.Logger,
) Service {
	return &service{
		repo:          repo,
		snapshotSvc:   snapshotSvc,
		analyticsRepo: analyticsRepo,
		playlistSvc:   playlistSvc,
		txProvider:    txProvider,
		logger:        logger,
	}
}

func (s *service) SendRequest(ctx context.Context, senderID, receiverID uuid.UUID) error {
	if senderID == receiverID {
		return apperrors.ErrCannotAddSelf
	}

	ok, err := s.repo.AreFriends(ctx, senderID, receiverID)
	if err != nil {
		return err
	}
	if ok {
		return apperrors.ErrAlreadyFriends
	}

	return s.repo.SendRequest(ctx, senderID, receiverID)
}

func (s *service) AcceptRequest(ctx context.Context, requestID int64, receiverID uuid.UUID) error {
	req, err := s.repo.GetRequest(ctx, requestID)
	if err != nil {
		return err
	}
	if req.ReceiverID != receiverID {
		return apperrors.ErrForbidden
	}
	if req.Status != "pending" {
		return apperrors.ErrRequestNotFound
	}

	return s.txProvider.RunInTx(ctx, func(exec database.Executor) error {
		txRepo := NewRepository(exec)
		return txRepo.AcceptRequest(ctx, requestID, req.SenderID, req.ReceiverID)
	})
}

func (s *service) RejectRequest(ctx context.Context, requestID int64, receiverID uuid.UUID) error {
	return s.repo.RejectRequest(ctx, requestID, receiverID)
}

func (s *service) ListIncoming(ctx context.Context, userID uuid.UUID) ([]IncomingRequest, error) {
	return s.repo.ListIncoming(ctx, userID)
}

func (s *service) ListFriends(ctx context.Context, userID uuid.UUID) ([]FriendProfile, error) {
	return s.repo.ListFriends(ctx, userID)
}

func (s *service) CompareTastes(ctx context.Context, userID, friendID uuid.UUID) (*TasteComparison, error) {
	ok, err := s.repo.AreFriends(ctx, userID, friendID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperrors.ErrNotFriends
	}

	mine, err := s.tasteScoresForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	theirs, err := s.tasteScoresForUser(ctx, friendID)
	if err != nil {
		return nil, err
	}

	return &TasteComparison{
		Mine:               mine,
		Theirs:             theirs,
		CompatibilityScore: compatibilityScore(mine, theirs),
	}, nil
}

func (s *service) GetFriendPlaylists(ctx context.Context, userID, friendID uuid.UUID) ([]playlist.Playlist, error) {
	ok, err := s.repo.AreFriends(ctx, userID, friendID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperrors.ErrNotFriends
	}

	connected, err := s.repo.HasSpotifyConnected(ctx, friendID)
	if err != nil {
		return nil, err
	}
	if !connected {
		return nil, apperrors.ErrSpotifyNotConnected
	}

	return s.playlistSvc.ListByUserID(ctx, friendID)
}

func (s *service) tasteScoresForUser(ctx context.Context, userID uuid.UUID) (TasteScores, error) {
	snap, err := s.snapshotSvc.GetLatestSnapshot(ctx, userID)
	if err != nil {
		return TasteScores{}, err
	}

	avgs, count, err := s.analyticsRepo.GetAveragesBySnapshotID(ctx, snap.ID)
	if err != nil {
		return TasteScores{}, err
	}
	if count == 0 {
		return TasteScores{}, apperrors.ErrNoSnapshot
	}

	return TasteScores{
		Valence:      avgs.Valence,
		Energy:       avgs.Energy,
		Danceability: avgs.Danceability,
		Acousticness: avgs.Acousticness,
	}, nil
}

func compatibilityScore(a, b TasteScores) float64 {
	dv := a.Valence - b.Valence
	de := a.Energy - b.Energy
	dd := a.Danceability - b.Danceability
	da := a.Acousticness - b.Acousticness

	dist := math.Sqrt(dv*dv + de*de + dd*dd + da*da)
	score := math.Max(0, 100*(1-dist/2))
	return math.Round(score*10) / 10
}
