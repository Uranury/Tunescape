package report

import (
	"context"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/track"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

type Repository interface {
	GetLatestSnapshotTopTracks(ctx context.Context, userID uuid.UUID) ([]track.Track, error)
}

type repository struct {
	exec database.Executor
}

func NewRepository(exec database.Executor) Repository {
	return &repository{exec: exec}
}

func (r *repository) GetLatestSnapshotTopTracks(ctx context.Context, userID uuid.UUID) ([]track.Track, error) {
	query := `
		SELECT t.id,
			t.spotify_id,
			t.name,
			COALESCE(t.popularity, 0) AS popularity,
			t.image_url,
			t.artist_name
		FROM snapshot_tracks st
		JOIN tracks t ON t.id = st.track_id
		WHERE st.snapshot_id = (
			SELECT id FROM snapshots
			WHERE user_id = $1
			ORDER BY created_at DESC
			LIMIT 1
		)
		ORDER BY st.position
		LIMIT 10
	`
	var tracks []track.Track
	if err := r.exec.SelectContext(ctx, &tracks, query, userID); err != nil {
		return nil, err
	}
	return tracks, nil
}