package user

import "time"

type User struct {
	ID          int64     `json:"id" db:"id"`
	SpotifyID   string    `json:"spotify_id" db:"spotify_id"`
	Email       string    `json:"email" db:"email"`
	DisplayName string    `json:"display_name" db:"display_name"`
	AvatarURL   *string   `json:"avatar_url" db:"avatar_url"`
	Country     *string   `json:"country" db:"country"`
	Product     *string   `json:"product" db:"product"`
	IsDeleted   bool      `json:"is_deleted" db:"is_deleted"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}