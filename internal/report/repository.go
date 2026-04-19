package report

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/track"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

type Repository interface {
	GetUserDisplayName(ctx context.Context, userID uuid.UUID) (string, error)
	GetLatestSnapshotTopTracks(ctx context.Context, userID uuid.UUID) ([]track.Track, error)
}

type repository struct {
	exec database.Executor
}

func NewRepository(exec database.Executor) Repository {
	return &repository{exec: exec}
}

func (r *repository) GetUserDisplayName(ctx context.Context, userID uuid.UUID) (string, error) {
	var name string
	err := r.exec.QueryRowxContext(ctx, `SELECT display_name FROM users WHERE id = $1 AND is_deleted = FALSE`, userID).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", apperrors.ErrNotFound
		}
		return "", err
	}
	return name, nil
}

func (r *repository) GetLatestSnapshotTopTracks(ctx context.Context, userID uuid.UUID) ([]track.Track, error) {
	query := `
		SELECT t.id, t.spotify_id, t.name, t.popularity
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