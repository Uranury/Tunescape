package snapshot

import (
	"time"

	"github.com/google/uuid"
)

type Track struct {
	ID         uuid.UUID `json:"id"         db:"id"`
	SpotifyID  string    `json:"spotify_id" db:"spotify_id"`
	Name       string    `json:"name"       db:"name"`
	Popularity int       `json:"popularity" db:"popularity"`
}

type Snapshot struct {
	ID        uuid.UUID `json:"id"         db:"id"`
	UserID    uuid.UUID `json:"user_id"    db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Tracks    []Track   `json:"tracks"`
}

type SnapshotTrack struct {
	SnapshotID uuid.UUID `db:"snapshot_id"`
	TrackID    uuid.UUID `db:"track_id"`
	Position   int       `db:"position"`
}