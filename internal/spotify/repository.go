package spotify

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

type Repository interface {
	UpsertTokens(ctx context.Context, userID uuid.UUID, accessToken, refreshToken string, expiresAt time.Time) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Token, error)
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type repository struct {
	exec database.Executor
}

func NewRepository(exec database.Executor) Repository {
	return &repository{exec: exec}
}

func (r *repository) UpsertTokens(ctx context.Context, userID uuid.UUID, accessToken, refreshToken string, expiresAt time.Time) error {
	query := `
		INSERT INTO spotify_tokens (user_id, access_token, refresh_token, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			access_token  = EXCLUDED.access_token,
			refresh_token = EXCLUDED.refresh_token,
			expires_at    = EXCLUDED.expires_at,
			updated_at    = NOW()
	`
	_, err := r.exec.ExecContext(ctx, query, userID, accessToken, refreshToken, expiresAt)
	return err
}

func (r *repository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.exec.ExecContext(ctx, `DELETE FROM spotify_tokens WHERE user_id = $1`, userID)
	return err
}

func (r *repository) GetByUserID(ctx context.Context, userID uuid.UUID) (*Token, error) {
	query := `SELECT * FROM spotify_tokens WHERE user_id = $1`
	token := &Token{}
	err := r.exec.QueryRowxContext(ctx, query, userID).StructScan(token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrSpotifyNotConnected
		}
		return nil, err
	}
	return token, nil
}
