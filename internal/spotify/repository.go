package spotify

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

type Repository interface {
	UpsertTokens(ctx context.Context, userID uuid.UUID, accessToken, refreshToken string, expiresAt time.Time) error
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
