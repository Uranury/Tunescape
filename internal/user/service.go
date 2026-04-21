package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"golang.org/x/crypto/bcrypt"
)

type ProfileResponse struct {
	DisplayName      string  `json:"display_name"`
	Email            string  `json:"email"`
	AvatarURL        *string `json:"avatar_url"`
	SpotifyConnected bool    `json:"spotify_connected"`
}

type Service interface {
	ValidateCredentials(ctx context.Context, email, password string) (*User, error)
	Create(ctx context.Context, email, password, displayName string) (*User, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*ProfileResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ValidateCredentials(ctx context.Context, email, password string) (*User, error) {
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("find user: %w", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	return u, nil
}

func (s *service) Create(ctx context.Context, email, password, displayName string) (*User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := &User{
		Email:       email,
		Password:    string(hashed),
		DisplayName: displayName,
		Role:        "user",
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *service) GetProfile(ctx context.Context, userID uuid.UUID) (*ProfileResponse, error) {
	u, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &ProfileResponse{
		DisplayName:      u.DisplayName,
		Email:            u.Email,
		AvatarURL:        u.AvatarURL,
		SpotifyConnected: u.SpotifyID != nil,
	}, nil
}
