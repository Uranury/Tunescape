package playlist

import "time"

type Playlist struct {
	SpotifyPlaylistID string    `json:"id"           db:"spotify_playlist_id"`
	Name              string    `json:"name"         db:"name"`
	ExternalURL       string    `json:"external_url" db:"external_url"`
	EmbedURL          string    `json:"embed_url"    db:"embed_url"`
	CreatedAt         time.Time `json:"created_at"   db:"created_at"`
}

type Response struct {
	PlaylistID  string `json:"playlist_id"`
	Name        string `json:"name"`
	ExternalURL string `json:"external_url"`
	EmbedURL    string `json:"embed_url"`
}
