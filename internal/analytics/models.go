package analytics

import "github.com/google/uuid"

type TrackAudioFeatures struct {
	TrackID          uuid.UUID `db:"track_id"`
	Danceability     float64   `db:"danceability"`
	Valence          float64   `db:"valence"`
	Energy           float64   `db:"energy"`
	Acousticness     float64   `db:"acousticness"`
	Instrumentalness float64   `db:"instrumentalness"`
	Liveness         float64   `db:"liveness"`
	Speechiness      float64   `db:"speechiness"`
	Tempo            float64   `db:"tempo"`
	Loudness         float64   `db:"loudness"`
}

type AudioFeatureAverages struct {
	Danceability     float64 `json:"danceability"     db:"danceability"`
	Valence          float64 `json:"valence"          db:"valence"`
	Energy           float64 `json:"energy"           db:"energy"`
	Acousticness     float64 `json:"acousticness"     db:"acousticness"`
	Instrumentalness float64 `json:"instrumentalness" db:"instrumentalness"`
	Liveness         float64 `json:"liveness"         db:"liveness"`
	Speechiness      float64 `json:"speechiness"      db:"speechiness"`
	Tempo            float64 `json:"tempo"            db:"tempo"`
	Loudness         float64 `json:"loudness"         db:"loudness"`
}

type MusicTasteResponse struct {
	SnapshotID  uuid.UUID            `json:"snapshot_id"`
	TracksCount int                  `json:"tracks_count"`
	Averages    AudioFeatureAverages `json:"averages"`
}
