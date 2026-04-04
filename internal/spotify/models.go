package spotify

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	ID           int64     `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	AccessToken  string    `db:"access_token"`
	RefreshToken string    `db:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
