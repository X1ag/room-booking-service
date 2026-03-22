package auth

import (
	"context"
	"errors"
	"test-backend-1-X1ag/internal/config"
	"test-backend-1-X1ag/internal/logger"

	"github.com/google/uuid"
)

type TokenManager interface {
	Generate(userID uuid.UUID, role string) (string, error)
	Parse(token string) (*Claims, error)
}

var (
	ErrInvalidRole = errors.New("invalid role")
	ErrInvalidGenerate = errors.New("invalid generate JWT token")
	ErrInvalidToken = errors.New("invalid token")
)

type AuthUsecase struct {
	jwtManager TokenManager
	cfg        config.AuthConfig
	logger 		*logger.ZerologLogger
}

func NewAuthUsecase(jwtManager TokenManager, cfg config.AuthConfig, l *logger.ZerologLogger) *AuthUsecase {
	return &AuthUsecase{
		jwtManager: jwtManager,
		cfg:        cfg,
		logger:     l,
	}
}

func (u *AuthUsecase) GenerateToken(userID uuid.UUID, role string) (string, error) {
	return u.jwtManager.Generate(userID, role)
}

func (u *AuthUsecase) DummyLogin(ctx context.Context, role string) (string, error) {
	var userID uuid.UUID
	switch role {
	case "admin":
		 userID = u.cfg.DummyAdminID
	case "user":
		userID = u.cfg.DummyUserID
	default:
		return "", ErrInvalidRole
	}
	
	return u.jwtManager.Generate(userID, role)
}