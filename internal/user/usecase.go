package user

import (
	"test-backend-1-X1ag/internal/logger"
)

type UserUsecase struct {
	repo Repository
	logger *logger.ZerologLogger
}

func NewUserUsecase(repo Repository, logger *logger.ZerologLogger) *UserUsecase {
	return &UserUsecase{
		repo: repo,
		logger: logger,
	}
}