package auth

import (
	"context"
	"errors"
	"testing"

	"test-backend-1-X1ag/internal/config"
	"test-backend-1-X1ag/internal/logger"
	"test-backend-1-X1ag/internal/user"

	"github.com/google/uuid"
)

type fakeTokenManager struct {
	generateFn func(userID uuid.UUID, role string) (string, error)
	parseFn    func(token string) (*Claims, error)
}

func (f *fakeTokenManager) Generate(userID uuid.UUID, role string) (string, error) {
	return f.generateFn(userID, role)
}

func (f *fakeTokenManager) Parse(token string) (*Claims, error) {
	if f.parseFn != nil {
		return f.parseFn(token)
	}
	return nil, nil
}

type fakeAuthUserService struct {
	registerFn     func(ctx context.Context, email, password, role string) (*user.User, error)
	authenticateFn func(ctx context.Context, email, password string) (*user.User, error)
}

func (f *fakeAuthUserService) Register(ctx context.Context, email, password, role string) (*user.User, error) {
	return f.registerFn(ctx, email, password, role)
}

func (f *fakeAuthUserService) Authenticate(ctx context.Context, email, password string) (*user.User, error) {
	return f.authenticateFn(ctx, email, password)
}

func TestAuthUsecaseRegisterAndLogin(t *testing.T) {
	t.Run("register delegates to user service", func(t *testing.T) {
		expectedUser := &user.User{ID: uuid.New(), Email: "user@example.com", Role: "user"}
		usecase := NewAuthUsecase(
			&fakeTokenManager{
				generateFn: func(userID uuid.UUID, role string) (string, error) { return "token", nil },
			},
			&fakeAuthUserService{
				registerFn: func(ctx context.Context, email, password, role string) (*user.User, error) {
					return expectedUser, nil
				},
				authenticateFn: func(ctx context.Context, email, password string) (*user.User, error) {
					return nil, nil
				},
			},
			config.AuthConfig{},
			logger.NewTestLogger(),
		)

		createdUser, err := usecase.Register(context.Background(), "user@example.com", "password123", "user")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if createdUser.Email != expectedUser.Email {
			t.Fatalf("expected email %s, got %s", expectedUser.Email, createdUser.Email)
		}
	})

	t.Run("login returns jwt", func(t *testing.T) {
		expectedUser := &user.User{ID: uuid.New(), Email: "user@example.com", Role: "user"}
		usecase := NewAuthUsecase(
			&fakeTokenManager{
				generateFn: func(userID uuid.UUID, role string) (string, error) {
					if userID != expectedUser.ID {
						t.Fatalf("expected user id %s, got %s", expectedUser.ID, userID)
					}
					if role != expectedUser.Role {
						t.Fatalf("expected role %s, got %s", expectedUser.Role, role)
					}
					return "jwt-token", nil
				},
			},
			&fakeAuthUserService{
				registerFn: func(ctx context.Context, email, password, role string) (*user.User, error) {
					return nil, nil
				},
				authenticateFn: func(ctx context.Context, email, password string) (*user.User, error) {
					return expectedUser, nil
				},
			},
			config.AuthConfig{},
			logger.NewTestLogger(),
		)

		token, err := usecase.Login(context.Background(), "user@example.com", "password123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token != "jwt-token" {
			t.Fatalf("expected jwt-token, got %s", token)
		}
	})

	t.Run("login returns invalid credentials", func(t *testing.T) {
		usecase := NewAuthUsecase(
			&fakeTokenManager{
				generateFn: func(userID uuid.UUID, role string) (string, error) { return "token", nil },
			},
			&fakeAuthUserService{
				registerFn: func(ctx context.Context, email, password, role string) (*user.User, error) {
					return nil, nil
				},
				authenticateFn: func(ctx context.Context, email, password string) (*user.User, error) {
					return nil, user.ErrInvalidCredentials
				},
			},
			config.AuthConfig{},
			logger.NewTestLogger(),
		)

		_, err := usecase.Login(context.Background(), "user@example.com", "wrong-password")
		if !errors.Is(err, user.ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("login returns token generation error", func(t *testing.T) {
		usecase := NewAuthUsecase(
			&fakeTokenManager{
				generateFn: func(userID uuid.UUID, role string) (string, error) { return "", errors.New("boom") },
			},
			&fakeAuthUserService{
				registerFn: func(ctx context.Context, email, password, role string) (*user.User, error) {
					return nil, nil
				},
				authenticateFn: func(ctx context.Context, email, password string) (*user.User, error) {
					return &user.User{ID: uuid.New(), Role: "user"}, nil
				},
			},
			config.AuthConfig{},
			logger.NewTestLogger(),
		)

		_, err := usecase.Login(context.Background(), "user@example.com", "password123")
		if !errors.Is(err, ErrInvalidGenerate) {
			t.Fatalf("expected ErrInvalidGenerate, got %v", err)
		}
	})
}

func TestAuthUsecaseDummyLogin(t *testing.T) {
	cfg := config.AuthConfig{
		DummyAdminID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		DummyUserID:  uuid.MustParse("22222222-2222-2222-2222-222222222222"),
	}

	t.Run("returns token for user role", func(t *testing.T) {
		usecase := NewAuthUsecase(
			&fakeTokenManager{
				generateFn: func(userID uuid.UUID, role string) (string, error) {
					if userID != cfg.DummyUserID {
						t.Fatalf("expected user id %s, got %s", cfg.DummyUserID, userID)
					}
					if role != "user" {
						t.Fatalf("expected role user, got %s", role)
					}
					return "dummy-token", nil
				},
			},
			&fakeAuthUserService{
				registerFn: func(ctx context.Context, email, password, role string) (*user.User, error) {
					return nil, nil
				},
				authenticateFn: func(ctx context.Context, email, password string) (*user.User, error) {
					return nil, nil
				},
			},
			cfg,
			logger.NewTestLogger(),
		)

		token, err := usecase.DummyLogin(context.Background(), "user")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token != "dummy-token" {
			t.Fatalf("expected dummy-token, got %s", token)
		}
	})

	t.Run("returns invalid role", func(t *testing.T) {
		usecase := NewAuthUsecase(
			&fakeTokenManager{
				generateFn: func(userID uuid.UUID, role string) (string, error) { return "dummy-token", nil },
			},
			&fakeAuthUserService{
				registerFn: func(ctx context.Context, email, password, role string) (*user.User, error) {
					return nil, nil
				},
				authenticateFn: func(ctx context.Context, email, password string) (*user.User, error) {
					return nil, nil
				},
			},
			cfg,
			logger.NewTestLogger(),
		)

		_, err := usecase.DummyLogin(context.Background(), "manager")
		if !errors.Is(err, ErrInvalidRole) {
			t.Fatalf("expected ErrInvalidRole, got %v", err)
		}
	})
}
