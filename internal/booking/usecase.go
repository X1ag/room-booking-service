package booking

import (
	"test-backend-1-X1ag/internal/logger"
)

type BookingUsecase struct {
	repo Repository
	logger *logger.ZerologLogger
}

func NewBookingUsecase(repo Repository, logger *logger.ZerologLogger) *BookingUsecase {
	return &BookingUsecase{
		repo: repo,
		logger: logger,
	}
}