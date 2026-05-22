package playlist

import (
	"context"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

type Repository interface {
	Insert(ctx context.Context, userID uuid.UUID, p *Playlist) error
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]Playlist, error)
}

type repository struct {
	db database.Executor
}

func NewRepository(db database.Executor) Repository {
	return &repository{db: db}
}

func (r *repository) Insert(ctx context.Context, userID uuid.UUID, p *Playlist) error {
	const q = `
		INSERT INTO playlists (user_id, spotify_playlist_id, name, external_url, embed_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at`
	return r.db.QueryRowxContext(ctx, q,
		userID, p.SpotifyPlaylistID, p.Name, p.ExternalURL, p.EmbedURL,
	).Scan(&p.CreatedAt)
}

func (r *repository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]Playlist, error) {
	const q = `
		SELECT spotify_playlist_id, name, external_url, embed_url, created_at
		FROM playlists
		WHERE user_id = $1
		ORDER BY created_at DESC`
	var out []Playlist
	if err := r.db.SelectContext(ctx, &out, q, userID); err != nil {
		return nil, err
	}
	return out, nil
}
