package snapshot

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
	CreateSnapshot(ctx context.Context, s *Snapshot) error
	CreateSnapshotTrack(ctx context.Context, st *SnapshotTrack) error
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]SnapshotSummary, error)
	GetByID(ctx context.Context, snapshotID, userID uuid.UUID) (*Snapshot, error)
	GetLatestByUserID(ctx context.Context, userID uuid.UUID) (*Snapshot, error)
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

func (r *repository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]SnapshotSummary, error) {
	var summaries []SnapshotSummary
	query := `SELECT id, user_id, created_at FROM snapshots WHERE user_id = $1 ORDER BY created_at DESC`
	if err := r.exec.SelectContext(ctx, &summaries, query, userID); err != nil {
		return nil, err
	}
	return summaries, nil
}

func (r *repository) GetLatestByUserID(ctx context.Context, userID uuid.UUID) (*Snapshot, error) {
	var snap Snapshot
	err := r.exec.QueryRowxContext(ctx,
		`SELECT id, user_id, created_at FROM snapshots WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`,
		userID,
	).Scan(&snap.ID, &snap.UserID, &snap.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.ErrNoSnapshot
	}
	if err != nil {
		return nil, err
	}

	var tracks []track.Track
	query := `
		SELECT t.id, t.spotify_id, t.name, t.popularity, t.image_url
		FROM snapshot_tracks st
		JOIN tracks t ON t.id = st.track_id
		WHERE st.snapshot_id = $1
		ORDER BY st.position
	`
	if err := r.exec.SelectContext(ctx, &tracks, query, snap.ID); err != nil {
		return nil, err
	}
	snap.Tracks = tracks
	return &snap, nil
}

func (r *repository) GetByID(ctx context.Context, snapshotID, userID uuid.UUID) (*Snapshot, error) {
	var snap Snapshot
	err := r.exec.QueryRowxContext(ctx,
		`SELECT id, user_id, created_at FROM snapshots WHERE id = $1 AND user_id = $2`,
		snapshotID, userID,
	).Scan(&snap.ID, &snap.UserID, &snap.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	var tracks []track.Track
	query := `
		SELECT t.id, t.spotify_id, t.name, t.popularity, t.image_url
		FROM snapshot_tracks st
		JOIN tracks t ON t.id = st.track_id
		WHERE st.snapshot_id = $1
		ORDER BY st.position
	`
	if err := r.exec.SelectContext(ctx, &tracks, query, snapshotID); err != nil {
		return nil, err
	}
	snap.Tracks = tracks
	return &snap, nil
}