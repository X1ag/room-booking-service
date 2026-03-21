package slot

import (
	"test-backend-1-X1ag/internal/logger"
)

type SlotUsecase struct {
	repo Repository
	logger *logger.ZerologLogger
}

func NewSlotUsecase(repo Repository, logger *logger.ZerologLogger) *SlotUsecase {
	return &SlotUsecase{
		repo: repo,
		logger: logger,
	}
}