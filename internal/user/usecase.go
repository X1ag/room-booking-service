package user

import (
	"context"
	"errors"
	"strings"

	"test-backend-1-X1ag/internal/logger"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo   Repository
	logger *logger.ZerologLogger
}

func NewUserUsecase(repo Repository, logger *logger.ZerologLogger) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (u *UserUsecase) Register(ctx context.Context, email, password, role string) (*User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	role = strings.ToLower(strings.TrimSpace(role))

	if email == "" || !strings.Contains(email, "@") {
		u.logger.Info().Str("email", email).Msg("email is invalid")
		return nil, ErrInvalidEmail
	}

	if strings.TrimSpace(password) == "" || len(password) < 8 || len([]byte(password)) > 72 {
		u.logger.Info().Str("email", email).Msg("password is invalid")
		return nil, ErrInvalidPassword
	}

	if role != "admin" && role != "user" {
		u.logger.Info().Str("role", role).Msg("role is invalid")
		return nil, ErrInvalidRole
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userToCreate := User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(passwordHash),
		Role:         role,
	}

	return u.repo.Create(ctx, &userToCreate)
}

func (u *UserUsecase) Authenticate(ctx context.Context, email, password string) (*User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" || !strings.Contains(email, "@") {
		u.logger.Info().Str("email", email).Msg("invalid credentials format")
		return nil, ErrInvalidCredentials
	}

	if password == "" || len([]byte(password)) > 72 {
		u.logger.Info().Str("email", email).Msg("invalid credentials format")
		return nil, ErrInvalidCredentials
	}

	foundUser, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return foundUser, nil
}
