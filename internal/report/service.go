package report

import (
	"context"

	"github.com/google/uuid"

	"gitlab.com/Uranury/tunescape/internal/leaderboard"
	"gitlab.com/Uranury/tunescape/internal/user"
)

type Service interface {
	GenerateReport(ctx context.Context, userID uuid.UUID) ([]byte, error)
}

type service struct {
	repo           Repository
	userRepo       user.Repository
	leaderboardSvc leaderboard.Service
}

func NewService(repo Repository, leaderboardSvc leaderboard.Service, userRepo user.Repository) Service {
	return &service{repo: repo, userRepo: userRepo, leaderboardSvc: leaderboardSvc}
}

func (s *service) GenerateReport(ctx context.Context, userID uuid.UUID) ([]byte, error) {
	name, err := s.userRepo.FindDisplayName(ctx, userID)
	if err != nil {
		return nil, err
	}

	tracks, err := s.repo.GetLatestSnapshotTopTracks(ctx, userID)
	if err != nil {
		return nil, err
	}

	rankings, err := s.leaderboardSvc.GetUserRankings(ctx, userID.String())
	if err != nil {
		rankings = &leaderboard.UserRankings{}
	}

	return buildPDF(name, tracks, rankings)
}