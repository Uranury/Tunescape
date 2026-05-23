package track

import "github.com/google/uuid"

type Track struct {
    ID         uuid.UUID `json:"id"          db:"id"`
    SpotifyID  string    `json:"spotify_id"  db:"spotify_id"`
    Name       string    `json:"name"        db:"name"`
    Popularity int       `json:"popularity"  db:"popularity"`
    ImageURL   *string   `json:"image_url"   db:"image_url"`
    ArtistName *string   `json:"artist_name" db:"artist_name"`
}
