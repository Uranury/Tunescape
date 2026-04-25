package reccobeats

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gitlab.com/Uranury/tunescape/pkg/config"
)

type AudioFeatures struct {
	Href             string  `json:"href"`
	Danceability     float64 `json:"danceability"`
	Valence          float64 `json:"valence"`
	Energy           float64 `json:"energy"`
	Acousticness     float64 `json:"acousticness"`
	Instrumentalness float64 `json:"instrumentalness"`
	Liveness         float64 `json:"liveness"`
	Speechiness      float64 `json:"speechiness"`
	Tempo            float64 `json:"tempo"`
	Loudness         float64 `json:"loudness"`
	Key              int     `json:"key"`
	Mode             int     `json:"mode"`
}

type audioFeaturesResponse struct {
	Content []AudioFeatures `json:"content"`
}

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(cfg config.Reccobeats, httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    cfg.BaseURL,
	}
}

func (c *Client) 	GetAudioFeaturesBatch(ctx context.Context, spotifyIDs []string) ([]AudioFeatures, error) {
	url := fmt.Sprintf("%s/v1/audio-features?ids=%s", c.baseURL, strings.Join(spotifyIDs, ","))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("reccobeats /v1/audio-features returned %d", resp.StatusCode)
	}

	var result audioFeaturesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Content, nil
}
