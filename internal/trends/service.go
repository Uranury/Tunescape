package trends

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	GetTrends(ctx context.Context, userID uuid.UUID) (*TrendsResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetTrends(ctx context.Context, userID uuid.UUID) (*TrendsResponse, error) {
	points, err := s.repo.GetTrendsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get trends: %w", err)
	}
	return &TrendsResponse{UserID: userID, Points: points}, nil
}