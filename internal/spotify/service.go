package spotify

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/track"
	"gitlab.com/Uranury/tunescape/internal/user"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

type Service interface {
	AuthURL(state string) string
	ConnectAccount(ctx context.Context, userID uuid.UUID, code string) error
	Disconnect(ctx context.Context, userID uuid.UUID) error
	GetValidToken(ctx context.Context, userID uuid.UUID) (string, error)
	GetTopTracks(ctx context.Context, userID uuid.UUID, limit int, timeRange TimeRange) ([]track.Track, error)
	CreatePlaylist(ctx context.Context, userID uuid.UUID, name string, trackURIs []string) (*PlaylistResult, error)
	UpsertTokens(ctx context.Context, userID uuid.UUID, accessToken, refreshToken string, expiresAt time.Time) error
}

type service struct {
	repo       Repository
	userRepo   user.Repository
	client     *Client
	txProvider database.TxProvider
	logger     *slog.Logger
}

func NewService(repo Repository, userRepo user.Repository, client *Client, txProvider database.TxProvider, logger *slog.Logger) Service {
	return &service{repo: repo, userRepo: userRepo, client: client, txProvider: txProvider, logger: logger}
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

func (s *service) GetTopTracks(ctx context.Context, userID uuid.UUID, limit int, timeRange TimeRange) ([]track.Track, error) {
  accessToken, err := s.GetValidToken(ctx, userID)
  if err != nil {
    return nil, err
  }
  items, err := s.client.GetTopTracks(ctx, accessToken, limit, timeRange)
  if err != nil {
    return nil, err
  }
  tracks := make([]track.Track, len(items))
  for i, item := range items {
    var imageURL *string
    if len(item.Album.Images) > 0 {
      url := item.Album.Images[0].URL
      imageURL = &url
    }
    tracks[i] = track.Track{
      SpotifyID:  item.ID,
      Name:       item.Name,
      Popularity: item.Popularity,
      ImageURL:   imageURL,
    }
  }
  return tracks, nil
}

func (s *service) Disconnect(ctx context.Context, userID uuid.UUID) error {
	return s.txProvider.RunInTx(ctx, func(exec database.Executor) error {
		if err := NewRepository(exec).DeleteByUserID(ctx, userID); err != nil {
			return fmt.Errorf("delete spotify tokens: %w", err)
		}
		if err := s.userRepo.ClearSpotify(ctx, userID); err != nil {
			return fmt.Errorf("clear spotify fields: %w", err)
		}
		return nil
	})
}

func (s *service) UpsertTokens(ctx context.Context, userID uuid.UUID, accessToken, refreshToken string, expiresAt time.Time) error {
	return s.repo.UpsertTokens(ctx, userID, accessToken, refreshToken, expiresAt)
}

func (s *service) CreatePlaylist(ctx context.Context, userID uuid.UUID, name string, trackURIs []string) (*PlaylistResult, error) {
	accessToken, err := s.GetValidToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	result, err := s.client.CreatePlaylist(ctx, accessToken, name)
	if err != nil {
		return nil, fmt.Errorf("create spotify playlist: %w", err)
	}

	if err := s.client.AddTracksToPlaylist(ctx, accessToken, result.ID, trackURIs); err != nil {
		return nil, fmt.Errorf("add tracks to playlist: %w", err)
	}

	return result, nil
}

func (s *service) GetValidToken(ctx context.Context, userID uuid.UUID) (string, error) {
	token, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return "", err
	}

	s.logger.Info("spotify token expiry check", "user_id", userID, "expires_at", token.ExpiresAt, "now", time.Now().UTC())

	if time.Now().Before(token.ExpiresAt.Add(-30 * time.Second)) {
		s.logger.Info("using cached spotify token", "user_id", userID)
		return token.AccessToken, nil
	}

	s.logger.Info("refreshing spotify token", "user_id", userID)
	refreshed, err := s.client.RefreshToken(ctx, token.RefreshToken)
	if err != nil {
		return "", fmt.Errorf("refresh spotify token: %w", err)
	}

	newRefreshToken := refreshed.RefreshToken
	if newRefreshToken == "" {
		newRefreshToken = token.RefreshToken
	}

	if err := s.repo.UpsertTokens(ctx, userID, refreshed.AccessToken, newRefreshToken, refreshed.Expiry); err != nil {
		return "", fmt.Errorf("persist refreshed token: %w", err)
	}

	return refreshed.AccessToken, nil
}