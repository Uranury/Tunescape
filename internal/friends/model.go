package friends

import (
	"time"

	"github.com/google/uuid"
)

type FriendRequest struct {
	ID         int64     `db:"id"`
	SenderID   uuid.UUID `db:"sender_id"`
	ReceiverID uuid.UUID `db:"receiver_id"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type FriendProfile struct {
	UserID           uuid.UUID `json:"user_id"`
	DisplayName      string    `json:"display_name"`
	SpotifyConnected bool      `json:"spotify_connected"`
	SpotifyID        *string   `json:"spotify_id,omitempty"`
}

type IncomingRequest struct {
	RequestID   int64     `json:"request_id" db:"request_id"`
	SenderID    uuid.UUID `json:"sender_id" db:"sender_id"`
	DisplayName string    `json:"display_name" db:"display_name"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type TasteComparison struct {
	Mine               TasteScores `json:"mine"`
	Theirs             TasteScores `json:"theirs"`
	CompatibilityScore float64     `json:"compatibility_score"`
}

type TasteScores struct {
	Valence      float64 `json:"valence"`
	Energy       float64 `json:"energy"`
	Danceability float64 `json:"danceability"`
	Acousticness float64 `json:"acousticness"`
}
