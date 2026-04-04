package analytics

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/reccobeats"
	"gitlab.com/Uranury/tunescape/internal/track"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

const audioFeaturesBatchSize = 40

type Service interface {
	GetMusicTaste(ctx context.Context, userID uuid.UUID) (*MusicTasteResponse, error)
}

type service struct {
	repo             Repository
	reccobeatsClient *reccobeats.Client
	txProvider       database.TxProvider
}

func NewService(
	repo Repository,
	reccobeatsClient *reccobeats.Client,
	txProvider database.TxProvider,
) Service {
	return &service{
		repo:             repo,
		reccobeatsClient: reccobeatsClient,
		txProvider:       txProvider,
	}
}

func (s *service) GetMusicTaste(ctx context.Context, userID uuid.UUID) (*MusicTasteResponse, error) {
	snap, err := s.repo.GetLatestSnapshotByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNoSnapshot) {
			return nil, apperrors.ErrNoSnapshot
		}
		return nil, fmt.Errorf("get latest snapshot: %w", err)
	}

	tracks, err := s.repo.GetTracksBySnapshotID(ctx, snap.ID)
	if err != nil {
		return nil, fmt.Errorf("get snapshot tracks: %w", err)
	}

	// Build a lookup map so features can be matched by Spotify ID, not by
	// positional index. The Reccobeats response includes an href field of the
	// form "https://open.spotify.com/track/<spotifyID>" which is the only
	// reliable per-item identifier the API returns.
	trackBySpotifyID := make(map[string]track.Track, len(tracks))
	for _, t := range tracks {
		trackBySpotifyID[t.SpotifyID] = t
	}

	err = s.txProvider.RunInTx(ctx, func(exec database.Executor) error {
		repo := NewRepository(exec)

		for start := 0; start < len(tracks); start += audioFeaturesBatchSize {
			end := start + audioFeaturesBatchSize
			if end > len(tracks) {
				end = len(tracks)
			}
			batch := tracks[start:end]

			spotifyIDs := make([]string, len(batch))
			for i, t := range batch {
				spotifyIDs[i] = t.SpotifyID
			}

			features, err := s.reccobeatsClient.GetAudioFeaturesBatch(ctx, spotifyIDs)
			if err != nil {
				return fmt.Errorf("fetch audio features batch [%d:%d]: %w", start, end, err)
			}

			for _, f := range features {
				// Extract the Spotify ID from the href URL path segment.
				spotifyID := path.Base(f.Href)
				t, ok := trackBySpotifyID[spotifyID]
				if !ok {
					continue
				}
				if err := repo.UpsertAudioFeatures(ctx, &TrackAudioFeatures{
					TrackID:          t.ID,
					Danceability:     f.Danceability,
					Valence:          f.Valence,
					Energy:           f.Energy,
					Acousticness:     f.Acousticness,
					Instrumentalness: f.Instrumentalness,
					Liveness:         f.Liveness,
					Speechiness:      f.Speechiness,
					Tempo:            f.Tempo,
					Loudness:         f.Loudness,
				}); err != nil {
					return fmt.Errorf("upsert audio features for track %s: %w", t.ID, err)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	avgs, count, err := s.repo.GetAveragesBySnapshotID(ctx, snap.ID)
	if err != nil {
		return nil, fmt.Errorf("aggregate audio features: %w", err)
	}

	return &MusicTasteResponse{
		SnapshotID:  snap.ID,
		TracksCount: count,
		Averages:    *avgs,
	}, nil
}
