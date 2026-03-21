package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id" db:"id"`
	SpotifyID   *string   `json:"spotify_id" db:"spotify_id"`
	Email       string    `json:"email" db:"email"`
	Password    string    `json:"-" db:"password"`
	DisplayName string    `json:"display_name" db:"display_name"`
	AvatarURL   *string   `json:"avatar_url" db:"avatar_url"`
	Country     *string   `json:"country" db:"country"`
	Product     *string   `json:"product" db:"product"`
	Role        string    `json:"role" db:"role"`
	IsDeleted   bool      `json:"is_deleted" db:"is_deleted"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
