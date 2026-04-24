package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"gitlab.com/Uranury/tunescape/pkg/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

type Client struct {
	httpClient *http.Client
	oauth2Cfg  *oauth2.Config
}

type spotifyProfile struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Country     string `json:"country"`
	Product     string `json:"product"`
	Images      []struct {
		URL string `json:"url"`
	} `json:"images"`
}

type topTrackItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
}

type topTracksResponse struct {
	Items []topTrackItem `json:"items"`
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

func (c *Client) getMe(ctx context.Context, accessToken string) (*spotifyProfile, error) {
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

	var profile spotifyProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (c *Client) GetTopTracks(ctx context.Context, accessToken string, limit int, timeRange TimeRange) ([]topTrackItem, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/top/tracks?limit=%d&time_range=%s", limit, timeRange)

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
		if resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests {
			return nil, fmt.Errorf("spotify /me/top/tracks returned %d: %w", resp.StatusCode, apperrors.ErrUpstreamUnavailable)
		}
		return nil, fmt.Errorf("spotify /me/top/tracks returned %d", resp.StatusCode)
	}

	var result topTracksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	ts := c.oauth2Cfg.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	return ts.Token()
}
