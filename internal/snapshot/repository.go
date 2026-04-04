package snapshot

import (
	"context"

	"gitlab.com/Uranury/tunescape/pkg/database"
)

type Repository interface {
	CreateSnapshot(ctx context.Context, s *Snapshot) error
	CreateSnapshotTrack(ctx context.Context, st *SnapshotTrack) error
}

type repository struct {
	exec database.Executor
}

func NewRepository(exec database.Executor) Repository {
	return &repository{exec: exec}
}

func (r *repository) CreateSnapshot(ctx context.Context, s *Snapshot) error {
	query := `
		INSERT INTO snapshots (user_id)
		VALUES ($1)
		RETURNING id, created_at
	`
	return r.exec.QueryRowxContext(ctx, query, s.UserID).Scan(&s.ID, &s.CreatedAt)
}

func (r *repository) CreateSnapshotTrack(ctx context.Context, st *SnapshotTrack) error {
	query := `
		INSERT INTO snapshot_tracks (snapshot_id, track_id, position)
		VALUES ($1, $2, $3)
	`
	_, err := r.exec.ExecContext(ctx, query, st.SnapshotID, st.TrackID, st.Position)
	return err
}
