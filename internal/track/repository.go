package track

import (
	"context"

	"gitlab.com/Uranury/tunescape/pkg/database"
)

type Repository interface {
	Upsert(ctx context.Context, track *Track) error
}

type repository struct {
	exec database.Executor
}

func NewRepository(exec database.Executor) Repository {
	return &repository{exec: exec}
}

func (r *repository) Upsert(ctx context.Context, track *Track) error {
	query := `
		INSERT INTO tracks (spotify_id, name, popularity)
		VALUES ($1, $2, $3)
		ON CONFLICT (spotify_id) DO UPDATE SET
			name       = EXCLUDED.name,
			popularity = EXCLUDED.popularity
		RETURNING id
	`
	return r.exec.QueryRowxContext(ctx, query,
		track.SpotifyID, track.Name, track.Popularity,
	).Scan(&track.ID)
}
