package room

import (
	"test-backend-1-X1ag/internal/logger"
)

type RoomUsecase struct {
	repo Repository
	logger *logger.ZerologLogger
}

func NewRoomUsecase(repo Repository, logger *logger.ZerologLogger) *RoomUsecase {
	return &RoomUsecase{
		repo: repo,
		logger: logger,
	}
}