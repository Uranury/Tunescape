package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/Uranury/tunescape/pkg/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

type Client struct {
	httpClient *http.Client
	oauth2Cfg  *oauth2.Config
}

type SpotifyProfile struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Country     string `json:"country"`
	Product     string `json:"product"`
	Images      []struct {
		URL string `json:"url"`
	} `json:"images"`
}

type TopTrackItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
}

type topTracksResponse struct {
	Items []TopTrackItem `json:"items"`
}

func NewClient(cfg config.Spotify, httpClient *http.Client) *Client {
	oauth2Cfg := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes: []string{
			"user-top-read",
			"user-read-recently-played",
			"user-read-private",
			"user-read-email",
		},
		Endpoint: spotify.Endpoint,
	}
	return &Client{
		httpClient: httpClient,
		oauth2Cfg:  oauth2Cfg,
	}
}

func (c *Client) getMe(ctx context.Context, accessToken string) (*SpotifyProfile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify /me returned %d", resp.StatusCode)
	}

	var profile SpotifyProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (c *Client) GetTopTracks(ctx context.Context, accessToken string, limit int) ([]TopTrackItem, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/top/tracks?limit=%d&time_range=medium_term", limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify /me/top/tracks returned %d", resp.StatusCode)
	}

	var result topTracksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Items, nil
}