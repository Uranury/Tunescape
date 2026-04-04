package snapshot

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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
}

func NewService(
	repo Repository,
	spotifySvc spotify.Service,
	txProvider database.TxProvider,
) Service {
	return &service{
		repo:       repo,
		spotifySvc: spotifySvc,
		txProvider: txProvider,
	}
}

func (s *service) CreateSnapshot(ctx context.Context, userID uuid.UUID) (*Snapshot, error) {
	topTracks, err := s.spotifySvc.GetTopTracks(ctx, userID, topTracksLimit)
	if err != nil {
		return nil, fmt.Errorf("fetch top tracks from spotify: %w", err)
	}

	var result *Snapshot

	err = s.txProvider.RunInTx(ctx, func(exec database.Executor) error {
		snapRepo := NewRepository(exec)
		trackRepo := track.NewRepository(exec)

		snap := &Snapshot{UserID: userID}
		if err := snapRepo.CreateSnapshot(ctx, snap); err != nil {
			return fmt.Errorf("create snapshot: %w", err)
		}

		for i := range topTracks {
			if err := trackRepo.Upsert(ctx, &topTracks[i]); err != nil {
				return fmt.Errorf("upsert track %q: %w", topTracks[i].SpotifyID, err)
			}

			if err := snapRepo.CreateSnapshotTrack(ctx, &SnapshotTrack{
				SnapshotID: snap.ID,
				TrackID:    topTracks[i].ID,
				Position:   i + 1,
			}); err != nil {
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

	return result, nil
}
