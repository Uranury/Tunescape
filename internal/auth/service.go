package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

var RefreshTokenTTL = time.Hour * 24 * 30

type RefreshTokenService interface {
	Generate(ctx context.Context, userID uuid.UUID, role string, userAgent string, ip string) (string, error)
	Validate(ctx context.Context, tokenString string) (*RefreshToken, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	Revoke(ctx context.Context, refreshToken string) error

	StartCleanup(ctx context.Context)
}

type refreshTokenService struct {
	tokenSvc   TokenService
	repo       Repository
	txProvider database.TxProvider
	logger     *slog.Logger
}

func NewRefreshService(tokenSvc TokenService, repo Repository, txProvider database.TxProvider, logger *slog.Logger) RefreshTokenService {
	return &refreshTokenService{
		tokenSvc:   tokenSvc,
		repo:       repo,
		txProvider: txProvider,
		logger:     logger,
	}
}

func (s *refreshTokenService) Generate(ctx context.Context, userID uuid.UUID, role string, userAgent string, ip string) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(b)
	tokenHash := hashToken(token)

	refreshToken := RefreshToken{
		UserID:    userID,
		Role:      role,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(RefreshTokenTTL),
		UserAgent: userAgent,
		IP:        ip,
	}
	if err := s.repo.Save(ctx, &refreshToken); err != nil {
		return "", fmt.Errorf("failed to save refresh token: %w", err)
	}
	return token, nil
}

func (s *refreshTokenService) Validate(ctx context.Context, tokenString string) (*RefreshToken, error) {
	tokenHash := hashToken(tokenString)

	token, err := s.repo.FindByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, fmt.Errorf("failed to validate refresh token: %w", err)
	}

	return token, nil
}

func (s *refreshTokenService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	var (
		accessToken string
		newToken    string
	)

	err := s.txProvider.RunInTx(ctx, func(exec database.Executor) error {
		repo := NewRepository(exec)

		tokenHash := hashToken(refreshToken)

		token, err := repo.FindByHashForUpdate(ctx, tokenHash)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("refresh token not found")
			}
			return fmt.Errorf("validate refresh token: %w", err)
		}

		accessToken, err = s.tokenSvc.Generate(token.UserID, token.Role)
		if err != nil {
			return fmt.Errorf("generate access token: %w", err)
		}

		if err := repo.RevokeByHash(ctx, tokenHash); err != nil {
			return fmt.Errorf("revoke refresh token: %w", err)
		}

		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			return fmt.Errorf("generate refresh token bytes: %w", err)
		}

		newToken = base64.URLEncoding.EncodeToString(b)
		newHash := hashToken(newToken)

		newRefreshToken := RefreshToken{
			UserID:    token.UserID,
			Role:      token.Role,
			TokenHash: newHash,
			ExpiresAt: time.Now().Add(RefreshTokenTTL),
			UserAgent: token.UserAgent,
			IP:        token.IP,
		}

		if err := repo.Save(ctx, &newRefreshToken); err != nil {
			return fmt.Errorf("save refresh token: %w", err)
		}

		return nil
	})

	if err != nil {
		s.logger.Error("refresh token failed", "err", err)
		return "", "", err
	}

	return accessToken, newToken, nil
}

func (s *refreshTokenService) Revoke(ctx context.Context, refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	if err := s.repo.RevokeByHash(ctx, tokenHash); err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	return nil
}

func (s *refreshTokenService) StartCleanup(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Hour * 24)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := s.repo.DeleteExpired(ctx); err != nil {
					s.logger.Warn("failed to delete expired refresh tokens", "err", err)
				}
			}
		}
	}()
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}
