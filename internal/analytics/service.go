package analytics

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"path"
	"time"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/cache"
	"gitlab.com/Uranury/tunescape/internal/reccobeats"
	"gitlab.com/Uranury/tunescape/internal/track"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

const audioFeaturesBatchSize = 40

type ScorePusher interface {
	PushScore(ctx context.Context, feature, userID string, score float64) error
}

type Service interface {
	GetMusicTaste(ctx context.Context, userID uuid.UUID) (*MusicTasteResponse, error)
}

type service struct {
	repo              Repository
	reccobeatsService reccobeats.Service
	txProvider        database.TxProvider
	logger            *slog.Logger
	cache             cache.Cache
	scorePusher       ScorePusher
}

func NewService(
	repo Repository,
	reccobeatsService reccobeats.Service,
	txProvider database.TxProvider,
	logger *slog.Logger,
	cache cache.Cache,
	scorePusher ScorePusher,
) Service {
	return &service{
		repo:              repo,
		reccobeatsService: reccobeatsService,
		txProvider:        txProvider,
		logger:            logger,
		cache:             cache,
		scorePusher:       scorePusher,
	}
}

func (s *service) GetMusicTaste(ctx context.Context, userID uuid.UUID) (*MusicTasteResponse, error) {
	var result MusicTasteResponse

	cached, err := s.cache.Get(ctx, "music_taste:"+userID.String())
	if err == nil && cached != nil {
		if unmarshalErr := json.Unmarshal(cached, &result); unmarshalErr == nil {
			return &result, nil
		}
		s.logger.Error("failed to unmarshal cache", "err", err)
	}

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

			features, err := s.reccobeatsService.GetAudioFeaturesBatch(ctx, spotifyIDs)
			if err != nil {
				return fmt.Errorf("fetch audio features batch [%d:%d]: %w", start, end, err)
			}

			toUpsert := make([]TrackAudioFeatures, 0, len(features))
			for _, f := range features {
				spotifyID := path.Base(f.Href)
				t, ok := trackBySpotifyID[spotifyID]
				if !ok {
					continue
				}
				toUpsert = append(toUpsert, TrackAudioFeatures{
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
				})
			}

			if err := repo.BulkUpsertAudioFeatures(ctx, toUpsert); err != nil {
				return fmt.Errorf("bulk upsert audio features batch [%d:%d]: %w", start, end, err)
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

	result = MusicTasteResponse{
		SnapshotID:  snap.ID,
		TracksCount: count,
		Averages:    *avgs,
	}

	if s.scorePusher != nil {
		for feature, score := range map[string]float64{
			"valence":      avgs.Valence,
			"energy":       avgs.Energy,
			"danceability": avgs.Danceability,
			"acousticness": avgs.Acousticness,
		} {
			if pushErr := s.scorePusher.PushScore(ctx, feature, userID.String(), score); pushErr != nil {
				s.logger.Warn("failed to push leaderboard score", "feature", feature, "error", pushErr)
			}
		}
	}

	resultBytes, err := json.Marshal(result)
	if err == nil {
		if cacheErr := s.cache.Set(ctx, "music_taste:"+userID.String(), resultBytes, time.Hour*24); cacheErr != nil {
			s.logger.Warn("failed to cache music_taste", "error", cacheErr)
		}
	} else {
		s.logger.Warn("failed to marshal music_taste", "error", err)
	}

	return &result, nil
}