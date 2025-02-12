package service

import (
	"context"
	"github.com/spanwalla/merch-store/internal/entity"
	"github.com/spanwalla/merch-store/internal/repository"
)

type UserReportService struct {
	userReportRepo repository.UserReport
}

func NewUserReportService(userReportRepo repository.UserReport) *UserReportService {
	return &UserReportService{userReportRepo: userReportRepo}
}

func (s *UserReportService) Get(ctx context.Context, userId int) (entity.UserReport, error) {
	return s.userReportRepo.Get(ctx, userId)
}
