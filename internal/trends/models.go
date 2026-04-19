package trends

import (
	"time"

	"github.com/google/uuid"
)

type SnapshotPoint struct {
	SnapshotID       uuid.UUID `json:"snapshot_id"      db:"snapshot_id"`
	CreatedAt        time.Time `json:"created_at"       db:"created_at"`
	Danceability     float64   `json:"danceability"     db:"danceability"`
	Valence          float64   `json:"valence"          db:"valence"`
	Energy           float64   `json:"energy"           db:"energy"`
	Acousticness     float64   `json:"acousticness"     db:"acousticness"`
	Instrumentalness float64   `json:"instrumentalness" db:"instrumentalness"`
	Liveness         float64   `json:"liveness"         db:"liveness"`
	Speechiness      float64   `json:"speechiness"      db:"speechiness"`
	Tempo            float64   `json:"tempo"            db:"tempo"`
	Loudness         float64   `json:"loudness"         db:"loudness"`
	TracksCount      int       `json:"tracks_count"     db:"tracks_count"`
}

type TrendsResponse struct {
	UserID uuid.UUID       `json:"user_id"`
	Points []SnapshotPoint `json:"points"`
}