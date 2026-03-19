package spotify

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SaveAccess(ctx context.Context, accessToken string) error {
	return s.repo.SaveAccess(ctx, accessToken)
}

func (s *Service) SaveRefresh(ctx context.Context, refreshToken string) error {
	return s.repo.SaveRefresh(ctx, refreshToken)
}
