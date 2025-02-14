package service

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
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
	report, err := s.userReportRepo.Get(ctx, userId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return entity.UserReport{}, ErrUserNotFound
		}
		log.Errorf("UserReportService.Get: %v", err)
		return entity.UserReport{}, ErrCannotGetReport
	}
	return report, nil
}
