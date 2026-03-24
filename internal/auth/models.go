package auth

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        int64      `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Role      string     `json:"role" db:"role"`
	TokenHash string     `json:"token_hash" db:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	RevokedAt *time.Time `json:"revoked_at" db:"revoked_at"`
	UserAgent string     `json:"user_agent" db:"user_agent"`
	IP        string     `json:"ip" db:"ip"`
}
