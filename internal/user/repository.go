package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

type Repository interface {
	ConnectSpotify(ctx context.Context, userID uuid.UUID, spotifyID *string, avatarURL, country, product *string) error
	Create(ctx context.Context, u *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	exec database.Executor
}

func NewRepository(exec database.Executor) Repository {
	return &repository{exec: exec}
}

func (r *repository) ConnectSpotify(ctx context.Context, userID uuid.UUID, spotifyID *string, avatarURL, country, product *string) error {
	query := `
		UPDATE users
		SET spotify_id = $2,
		    avatar_url = $3,
		    country    = $4,
		    product    = $5,
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.exec.ExecContext(ctx, query, userID, spotifyID, avatarURL, country, product)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return apperrors.ErrSpotifyIDTaken
		}
	}
	return err
}

func (r *repository) Create(ctx context.Context, u *User) error {
	query := `
		INSERT INTO users (email, password, display_name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := r.exec.QueryRowxContext(ctx, query,
		u.Email, u.Password, u.DisplayName, u.Role,
	).Scan(&u.ID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return apperrors.ErrEmailTaken
		}
		return err
	}
	return nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT * FROM users WHERE email = $1 AND is_deleted = FALSE`
	u := &User{}
	err := r.exec.QueryRowxContext(ctx, query, email).StructScan(u)
	return u, err
}
