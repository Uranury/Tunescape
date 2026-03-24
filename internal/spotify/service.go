package spotify

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/user"
)

type Service interface {
	AuthURL(state string) string
	ConnectAccount(ctx context.Context, userID uuid.UUID, code string) error
	UpsertTokens(ctx context.Context, userID uuid.UUID, accessToken, refreshToken string, expiresAt time.Time) error
}

type service struct {
	repo     Repository
	userRepo user.Repository
	client   *Client
}

func NewService(repo Repository, userRepo user.Repository, client *Client) Service {
	return &service{repo: repo, userRepo: userRepo, client: client}
}

func (s *service) AuthURL(state string) string {
	return s.client.oauth2Cfg.AuthCodeURL(state)
}

func (s *service) ConnectAccount(ctx context.Context, userID uuid.UUID, code string) error {
	token, err := s.client.oauth2Cfg.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("exchange oauth code: %w", err)
	}

	profile, err := s.client.getMe(ctx, token.AccessToken)
	if err != nil {
		return fmt.Errorf("fetch spotify profile: %w", err)
	}

	spotifyID := profile.ID
	var avatarURL, country, product *string
	if len(profile.Images) > 0 {
		avatarURL = &profile.Images[0].URL
	}
	if profile.Country != "" {
		country = &profile.Country
	}
	if profile.Product != "" {
		product = &profile.Product
	}

	if err := s.userRepo.ConnectSpotify(ctx, userID, &spotifyID, avatarURL, country, product); err != nil {
		return err
	}

	return s.repo.UpsertTokens(ctx, userID, token.AccessToken, token.RefreshToken, token.Expiry)
}

func (s *service) UpsertTokens(ctx context.Context, userID uuid.UUID, accessToken, refreshToken string, expiresAt time.Time) error {
	return s.repo.UpsertTokens(ctx, userID, accessToken, refreshToken, expiresAt)
}
