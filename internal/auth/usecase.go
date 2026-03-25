package auth

import (
	"context"
	"errors"

	"test-backend-1-X1ag/internal/config"
	"test-backend-1-X1ag/internal/logger"
	"test-backend-1-X1ag/internal/user"

	"github.com/google/uuid"
)

type TokenManager interface {
	Generate(userID uuid.UUID, role string) (string, error)
	Parse(token string) (*Claims, error)
}

type UserService interface {
	Register(ctx context.Context, email, password, role string) (*user.User, error)
	Authenticate(ctx context.Context, email, password string) (*user.User, error)
}

var (
	ErrInvalidRole     = errors.New("invalid role")
	ErrInvalidGenerate = errors.New("invalid generate JWT token")
	ErrInvalidToken    = errors.New("invalid token")
)

type AuthUsecase struct {
	jwtManager  TokenManager
	userService UserService
	cfg         config.AuthConfig
	logger      *logger.ZerologLogger
}

func NewAuthUsecase(jwtManager TokenManager, userService UserService, cfg config.AuthConfig, l *logger.ZerologLogger) *AuthUsecase {
	return &AuthUsecase{
		jwtManager:  jwtManager,
		userService: userService,
		cfg:         cfg,
		logger:      l,
	}
}

func (u *AuthUsecase) GenerateToken(userID uuid.UUID, role string) (string, error) {
	token, err := u.jwtManager.Generate(userID, role)
	if err != nil {
		return "", ErrInvalidGenerate
	}

	return token, nil
}

func (u *AuthUsecase) Register(ctx context.Context, email, password, role string) (*user.User, error) {
	return u.userService.Register(ctx, email, password, role)
}

func (u *AuthUsecase) Login(ctx context.Context, email, password string) (string, error) {
	foundUser, err := u.userService.Authenticate(ctx, email, password)
	if err != nil {
		return "", err
	}

	return u.GenerateToken(foundUser.ID, foundUser.Role)
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

	return u.GenerateToken(userID, role)
}
