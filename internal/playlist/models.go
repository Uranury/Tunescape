package playlist

type Response struct {
	PlaylistID  string `json:"playlist_id"`
	Name        string `json:"name"`
	ExternalURL string `json:"external_url"`
	EmbedURL    string `json:"embed_url"`
}
