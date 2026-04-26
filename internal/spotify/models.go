package spotify

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID           int64     `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	AccessToken  string    `db:"access_token"`
	RefreshToken string    `db:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type TimeRange string

const (
	ShortTerm  TimeRange = "short_term"
	MediumTerm TimeRange = "medium_term"
	LongTerm   TimeRange = "long_term"
)

func ParseTimeRange(s string) (TimeRange, error) {
	switch TimeRange(s) {
	case ShortTerm, MediumTerm, LongTerm:
		return TimeRange(s), nil
	case "":
		return MediumTerm, nil
	default:
		return "", fmt.Errorf("invalid time_range %q: must be short_term, medium_term, or long_term", s)
	}
}

type PlaylistResult struct {
	ID          string
	ExternalURL string
}
