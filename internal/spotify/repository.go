package spotify

import (
	"context"
	"gitlab.com/Uranury/tunescape/pkg/database"
	"time"
)

type Repository interface {
	SaveAccess(ctx context.Context, accessToken string) error
	SaveRefresh(ctx context.Context, refreshToken string) error
}

type repository struct {
	exec database.Executor
}

func NewRepository(exec database.Executor) Repository {
	return &repository{exec: exec}
}

func (r *repository) SaveRefresh(ctx context.Context, refreshToken string) error {
	query := `INSERT INTO spotify_refresh (refresh_token) VALUES (?)`
	_, err := r.exec.ExecContext(ctx, query, refreshToken)
	return err
}

func (r *repository) SaveAccess(ctx context.Context, accessToken string) error {
	query := `INSERT INTO spotify_access (access_token) VALUES (?)`
	_, err := r.exec.ExecContext(ctx, query, accessToken)
	return err
}

func (r *repository) Delete(ctx context.Context, tokenExpiry time.Time) error {
	query := `DELETE FROM spotify_access WHERE access_expiry`
}
