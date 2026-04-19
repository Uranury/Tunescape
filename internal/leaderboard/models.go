package leaderboard

type Entry struct {
	Rank        int     `json:"rank"`
	UserID      string  `json:"user_id"`
	DisplayName string  `json:"display_name"`
	Score       float64 `json:"score"`
}

type LeaderboardResponse struct {
	Feature string  `json:"feature"`
	Entries []Entry `json:"entries"`
}

type UserRankings struct {
	Valence      *int64 `json:"valence"`
	Energy       *int64 `json:"energy"`
	Danceability *int64 `json:"danceability"`
	Acousticness *int64 `json:"acousticness"`
}