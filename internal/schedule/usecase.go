package schedule

import (
	"test-backend-1-X1ag/internal/logger"
)

type ScheduleUsecase struct {
	repo Repository
	logger *logger.ZerologLogger
}

func NewSheduleUsecase(repo Repository, logger *logger.ZerologLogger) *ScheduleUsecase {
	return &ScheduleUsecase{
		repo: repo,
		logger: logger,
	}
}