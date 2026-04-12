package reccobeats

import "context"

type Service interface {
	GetAudioFeaturesBatch(ctx context.Context, spotifyIDs []string) ([]AudioFeatures, error)
}

type service struct {
	client *Client
}

func NewService(client *Client) Service {
	return &service{client: client}
}

func (s *service) GetAudioFeaturesBatch(ctx context.Context, spotifyIDs []string) ([]AudioFeatures, error) {
	return s.client.GetAudioFeaturesBatch(ctx, spotifyIDs)
}
