package analytics

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/snapshot"
	"gitlab.com/Uranury/tunescape/internal/track"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

type Repository interface {
	GetLatestSnapshotByUserID(ctx context.Context, userID uuid.UUID) (*snapshot.Snapshot, error)
	GetTracksBySnapshotID(ctx context.Context, snapshotID uuid.UUID) ([]track.Track, error)
	BulkUpsertAudioFeatures(ctx context.Context, features []TrackAudioFeatures) error
	GetAveragesBySnapshotID(ctx context.Context, snapshotID uuid.UUID) (*AudioFeatureAverages, int, error)
}

type repository struct {
	exec database.Executor
}

func NewRepository(exec database.Executor) Repository {
	return &repository{exec: exec}
}

func (r *repository) GetLatestSnapshotByUserID(ctx context.Context, userID uuid.UUID) (*snapshot.Snapshot, error) {
	query := `
		SELECT id, user_id, created_at
		FROM snapshots
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var s snapshot.Snapshot
	err := r.exec.QueryRowxContext(ctx, query, userID).StructScan(&s)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNoSnapshot
		}
		return nil, err
	}
	return &s, nil
}

func (r *repository) GetTracksBySnapshotID(ctx context.Context, snapshotID uuid.UUID) ([]track.Track, error) {
	query := `
		SELECT t.id, t.spotify_id, t.name, t.popularity
		FROM snapshot_tracks st
		JOIN tracks t ON t.id = st.track_id
		WHERE st.snapshot_id = $1
		ORDER BY st.position
	`
	var tracks []track.Track
	if err := r.exec.SelectContext(ctx, &tracks, query, snapshotID); err != nil {
		return nil, err
	}
	return tracks, nil
}

func (r *repository) BulkUpsertAudioFeatures(ctx context.Context, features []TrackAudioFeatures) error {
	if len(features) == 0 {
		return nil
	}

	const cols = 10
	placeholders := make([]string, len(features))
	args := make([]interface{}, 0, len(features)*cols)

	for i, f := range features {
		base := i * cols
		placeholders[i] = fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			base+1, base+2, base+3, base+4, base+5,
			base+6, base+7, base+8, base+9, base+10,
		)
		args = append(args,
			f.TrackID,
			f.Danceability,
			f.Valence,
			f.Energy,
			f.Acousticness,
			f.Instrumentalness,
			f.Liveness,
			f.Speechiness,
			f.Tempo,
			f.Loudness,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO track_audio_features
			(track_id, danceability, valence, energy, acousticness, instrumentalness, liveness, speechiness, tempo, loudness)
		VALUES %s
		ON CONFLICT (track_id) DO UPDATE SET
			danceability     = EXCLUDED.danceability,
			valence          = EXCLUDED.valence,
			energy           = EXCLUDED.energy,
			acousticness     = EXCLUDED.acousticness,
			instrumentalness = EXCLUDED.instrumentalness,
			liveness         = EXCLUDED.liveness,
			speechiness      = EXCLUDED.speechiness,
			tempo            = EXCLUDED.tempo,
			loudness         = EXCLUDED.loudness
	`, strings.Join(placeholders, ", "))

	_, err := r.exec.ExecContext(ctx, query, args...)
	return err
}

func (r *repository) GetAveragesBySnapshotID(ctx context.Context, snapshotID uuid.UUID) (*AudioFeatureAverages, int, error) {
	query := `
		SELECT
			AVG(taf.danceability)     AS danceability,
			AVG(taf.valence)          AS valence,
			AVG(taf.energy)           AS energy,
			AVG(taf.acousticness)     AS acousticness,
			AVG(taf.instrumentalness) AS instrumentalness,
			AVG(taf.liveness)         AS liveness,
			AVG(taf.speechiness)      AS speechiness,
			AVG(taf.tempo)            AS tempo,
			AVG(taf.loudness)         AS loudness,
			COUNT(*)                  AS tracks_count
		FROM snapshot_tracks st
		JOIN track_audio_features taf ON taf.track_id = st.track_id
		WHERE st.snapshot_id = $1
	`

	row := r.exec.QueryRowxContext(ctx, query, snapshotID)

	var avgs AudioFeatureAverages
	var count int
	if err := row.Scan(
		&avgs.Danceability,
		&avgs.Valence,
		&avgs.Energy,
		&avgs.Acousticness,
		&avgs.Instrumentalness,
		&avgs.Liveness,
		&avgs.Speechiness,
		&avgs.Tempo,
		&avgs.Loudness,
		&count,
	); err != nil {
		return nil, 0, err
	}

	return &avgs, count, nil
}
