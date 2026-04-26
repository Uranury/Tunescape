package spotify

import (
	"bytes"
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
			"playlist-modify-public",
			"playlist-modify-private",
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

type createPlaylistRequest struct {
	Name        string `json:"name"`
	Public      bool   `json:"public"`
	Description string `json:"description"`
}

type createPlaylistResponse struct {
	ID           string `json:"id"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}

type addTracksRequest struct {
	URIs []string `json:"uris"`
}

func (c *Client) CreatePlaylist(ctx context.Context, accessToken, name string) (*PlaylistResult, error) {
	body, _ := json.Marshal(createPlaylistRequest{
		Name:        name,
		Public:      false,
		Description: "Created by Tunescape",
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.spotify.com/v1/me/playlists", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("spotify create playlist returned %d: %w", resp.StatusCode, apperrors.ErrUpstreamUnavailable)
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("spotify create playlist returned %d: %w", resp.StatusCode, apperrors.ErrSpotifyNotConnected)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("spotify create playlist returned %d", resp.StatusCode)
	}

	var result createPlaylistResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &PlaylistResult{ID: result.ID, ExternalURL: result.ExternalURLs.Spotify}, nil
}

func (c *Client) AddTracksToPlaylist(ctx context.Context, accessToken, playlistID string, trackURIs []string) error {
	const batchSize = 10

	for start := 0; start < len(trackURIs); start += batchSize {
		end := start + batchSize
		if end > len(trackURIs) {
			end = len(trackURIs)
		}

		body, _ := json.Marshal(addTracksRequest{URIs: trackURIs[start:end]})
		url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/items", playlistID)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}
		_ = resp.Body.Close()

		if resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests {
			return fmt.Errorf("spotify add tracks returned %d: %w", resp.StatusCode, apperrors.ErrUpstreamUnavailable)
		}
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return fmt.Errorf("spotify add tracks returned %d: %w", resp.StatusCode, apperrors.ErrSpotifyNotConnected)
		}
		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			return fmt.Errorf("spotify add tracks returned %d", resp.StatusCode)
		}
	}

	return nil
}
