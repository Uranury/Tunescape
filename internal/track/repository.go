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
  const query = `
    INSERT INTO tracks (spotify_id, name, popularity, image_url)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (spotify_id) DO UPDATE SET
      name       = EXCLUDED.name,
      popularity = EXCLUDED.popularity,
      image_url  = EXCLUDED.image_url
    RETURNING id
  `
  return r.exec.QueryRowxContext(ctx, query,
    track.SpotifyID, track.Name, track.Popularity, track.ImageURL,
  ).Scan(&track.ID)
}