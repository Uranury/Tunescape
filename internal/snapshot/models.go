package snapshot

import (
	"time"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/track"
)

type Snapshot struct {
	ID        uuid.UUID     `json:"id"         db:"id"`
	UserID    uuid.UUID     `json:"user_id"    db:"user_id"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	Tracks    []track.Track `json:"tracks"`
}

type SnapshotSummary struct {
	ID        uuid.UUID `json:"id"         db:"id"`
	UserID    uuid.UUID `json:"user_id"    db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type SnapshotTrack struct {
	SnapshotID uuid.UUID `db:"snapshot_id"`
	TrackID    uuid.UUID `db:"track_id"`
	Position   int       `db:"position"`
}
