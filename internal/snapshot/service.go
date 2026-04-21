package snapshot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/cache"
	"gitlab.com/Uranury/tunescape/internal/spotify"
	"gitlab.com/Uranury/tunescape/internal/track"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

const topTracksLimit = 50

type Service interface {
	CreateSnapshot(ctx context.Context, userID uuid.UUID) (*Snapshot, error)
}

type service struct {
	repo       Repository
	spotifySvc spotify.Service
	txProvider database.TxProvider
	cache      cache.Cache
	logger     *slog.Logger
}

func NewService(
	repo Repository,
	spotifySvc spotify.Service,
	txProvider database.TxProvider,
	cache cache.Cache,
	logger *slog.Logger,
) Service {
	return &service{
		repo:       repo,
		spotifySvc: spotifySvc,
		txProvider: txProvider,
		cache:      cache,
		logger:     logger,
	}
}

func (s *service) CreateSnapshot(ctx context.Context, userID uuid.UUID) (*Snapshot, error) {
	topTracks, err := s.spotifySvc.GetTopTracks(ctx, userID, topTracksLimit)
	if err != nil {
		s.logger.Error("failed to fetch top tracks from spotify", "user_id", userID, "error", err)
		return nil, fmt.Errorf("fetch top tracks from spotify: %w", err)
	}

	var result *Snapshot

	err = s.txProvider.RunInTx(ctx, func(exec database.Executor) error {
		snapRepo := NewRepository(exec)
		trackRepo := track.NewRepository(exec)

		snap := &Snapshot{UserID: userID}
		if err := snapRepo.CreateSnapshot(ctx, snap); err != nil {
			s.logger.Error("failed to insert snapshot", "user_id", userID, "error", err)
			return fmt.Errorf("create snapshot: %w", err)
		}

		for i := range topTracks {
			if err := trackRepo.Upsert(ctx, &topTracks[i]); err != nil {
				s.logger.Error("failed to upsert track", "user_id", userID, "spotify_id", topTracks[i].SpotifyID, "error", err)
				return fmt.Errorf("upsert track %q: %w", topTracks[i].SpotifyID, err)
			}

			if err := snapRepo.CreateSnapshotTrack(ctx, &SnapshotTrack{
				SnapshotID: snap.ID,
				TrackID:    topTracks[i].ID,
				Position:   i + 1,
			}); err != nil {
				s.logger.Error("failed to link track to snapshot", "user_id", userID, "spotify_id", topTracks[i].SpotifyID, "error", err)
				return fmt.Errorf("link track %q to snapshot: %w", topTracks[i].SpotifyID, err)
			}
		}

		snap.Tracks = topTracks
		result = snap
		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := s.cache.Delete(ctx, "music_taste:"+userID.String()); err != nil {
		s.logger.Warn("failed to invalidate music_taste cache after snapshot creation", "user_id", userID, "error", err)
	}

	return result, nil
}
