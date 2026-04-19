package trends

import (
	"context"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

type Repository interface {
	GetTrendsByUserID(ctx context.Context, userID uuid.UUID) ([]SnapshotPoint, error)
}

type repository struct {
	exec database.Executor
}

func NewRepository(exec database.Executor) Repository {
	return &repository{exec: exec}
}

func (r *repository) GetTrendsByUserID(ctx context.Context, userID uuid.UUID) ([]SnapshotPoint, error) {
	query := `
		SELECT
			s.id                      AS snapshot_id,
			s.created_at,
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
		FROM snapshots s
		JOIN snapshot_tracks st       ON st.snapshot_id = s.id
		JOIN track_audio_features taf ON taf.track_id   = st.track_id
		WHERE s.user_id = $1
		GROUP BY s.id, s.created_at
		ORDER BY s.created_at ASC
	`
	var points []SnapshotPoint
	if err := r.exec.SelectContext(ctx, &points, query, userID); err != nil {
		return nil, err
	}
	return points, nil
}