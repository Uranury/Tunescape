package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

type Repository interface {
	ConnectSpotify(ctx context.Context, userID uuid.UUID, spotifyID *string, avatarURL, country, product *string) error
	ClearSpotify(ctx context.Context, userID uuid.UUID) error
	Create(ctx context.Context, u *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, userID uuid.UUID) (*User, error)
	FindDisplayName(ctx context.Context, userID uuid.UUID) (string, error)
	FindDisplayNamesByIDs(ctx context.Context, userIDs []string) (map[string]string, error)
	FindAll(ctx context.Context) ([]User, error)
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

func (r *repository) ClearSpotify(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users
		SET spotify_id = NULL,
		    avatar_url = NULL,
		    country    = NULL,
		    product    = NULL,
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.exec.ExecContext(ctx, query, userID)
	return err
}

func (r *repository) FindByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	query := `SELECT * FROM users WHERE id = $1 AND is_deleted = FALSE`
	u := &User{}
	err := r.exec.QueryRowxContext(ctx, query, userID).StructScan(u)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *repository) FindDisplayName(ctx context.Context, userID uuid.UUID) (string, error) {
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

func (r *repository) FindDisplayNamesByIDs(ctx context.Context, userIDs []string) (map[string]string, error) {
	if len(userIDs) == 0 {
		return map[string]string{}, nil
	}
	query := `SELECT id::text AS id, display_name FROM users WHERE id::text = ANY($1)`
	var rows []struct {
		ID          string `db:"id"`
		DisplayName string `db:"display_name"`
	}
	if err := r.exec.SelectContext(ctx, &rows, query, pq.Array(userIDs)); err != nil {
		return nil, err
	}
	result := make(map[string]string, len(rows))
	for _, row := range rows {
		result[row.ID] = row.DisplayName
	}
	return result, nil
}

func (r *repository) FindAll(ctx context.Context) ([]User, error) {
	query := `SELECT * FROM users WHERE is_deleted = FALSE ORDER BY created_at DESC`
	var users []User
	if err := r.exec.SelectContext(ctx, &users, query); err != nil {
		return nil, err
	}
	return users, nil
}
