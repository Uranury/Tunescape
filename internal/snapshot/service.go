package snapshot

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/spotify"
	"gitlab.com/Uranury/tunescape/internal/track"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

const topTracksLimit = 50

type Service interface {
	CreateSnapshot(ctx context.Context, userID uuid.UUID) (*Snapshot, error)
}

type service struct {
	repo          Repository
	spotifyRepo   spotify.Repository
	spotifyClient *spotify.Client
	txProvider    database.TxProvider
}

func NewService(
	repo Repository,
	spotifyRepo spotify.Repository,
	spotifyClient *spotify.Client,
	txProvider database.TxProvider,
) Service {
	return &service{
		repo:          repo,
		spotifyRepo:   spotifyRepo,
		spotifyClient: spotifyClient,
		txProvider:    txProvider,
	}
}

func (s *service) CreateSnapshot(ctx context.Context, userID uuid.UUID) (*Snapshot, error) {
	token, err := s.spotifyRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrSpotifyNotConnected) {
			return nil, apperrors.ErrSpotifyNotConnected
		}
		return nil, fmt.Errorf("get spotify token: %w", err)
	}

	topTracks, err := s.spotifyClient.GetTopTracks(ctx, token.AccessToken, topTracksLimit)
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

		tracks := make([]track.Track, 0, len(topTracks))
		for i, item := range topTracks {
			t := &track.Track{
				SpotifyID:  item.ID,
				Name:       item.Name,
				Popularity: item.Popularity,
			}

			if err := trackRepo.Upsert(ctx, t); err != nil {
				return fmt.Errorf("upsert track %q: %w", item.ID, err)
			}

			if err := snapRepo.CreateSnapshotTrack(ctx, &SnapshotTrack{
				SnapshotID: snap.ID,
				TrackID:    t.ID,
				Position:   i + 1,
			}); err != nil {
				return fmt.Errorf("link track %q to snapshot: %w", item.ID, err)
			}

			tracks = append(tracks, *t)
		}

		snap.Tracks = tracks
		result = snap
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
