package user

import (
	"context"
	"errors"
	"testing"

	"test-backend-1-X1ag/internal/logger"

	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepo struct {
	createFn     func(ctx context.Context, user *User) (*User, error)
	getByEmailFn func(ctx context.Context, email string) (*User, error)
}

func (f *fakeUserRepo) Create(ctx context.Context, user *User) (*User, error) {
	return f.createFn(ctx, user)
}

func (f *fakeUserRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	return f.getByEmailFn(ctx, email)
}

func TestUserUsecaseRegisterValidation(t *testing.T) {
	usecase := NewUserUsecase(&fakeUserRepo{}, logger.NewTestLogger())

	testCases := []struct {
		name     string
		email    string
		password string
		role     string
		expected error
	}{
		{"invalid email", "bad-email", "password123", "user", ErrInvalidEmail},
		{"invalid password", "user@example.com", "short", "user", ErrInvalidPassword},
		{"invalid role", "user@example.com", "password123", "manager", ErrInvalidRole},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := usecase.Register(context.Background(), tc.email, tc.password, tc.role)
			if !errors.Is(err, tc.expected) {
				t.Fatalf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestUserUsecaseRegister(t *testing.T) {
	t.Run("returns duplicate email", func(t *testing.T) {
		usecase := NewUserUsecase(&fakeUserRepo{
			createFn: func(ctx context.Context, user *User) (*User, error) {
				return nil, ErrEmailAlreadyExists
			},
		}, logger.NewTestLogger())

		_, err := usecase.Register(context.Background(), "user@example.com", "password123", "user")
		if !errors.Is(err, ErrEmailAlreadyExists) {
			t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
		}
	})

	t.Run("creates user with hashed password", func(t *testing.T) {
		usecase := NewUserUsecase(&fakeUserRepo{
			createFn: func(ctx context.Context, user *User) (*User, error) {
				if user.Email != "admin@example.com" {
					t.Fatalf("expected normalized email, got %s", user.Email)
				}
				if user.Role != "admin" {
					t.Fatalf("expected role admin, got %s", user.Role)
				}
				if user.PasswordHash == "" || user.PasswordHash == "password123" {
					t.Fatal("expected hashed password")
				}
				if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password123")); err != nil {
					t.Fatalf("expected valid bcrypt hash: %v", err)
				}

				return &User{
					ID:        user.ID,
					Email:     user.Email,
					Role:      user.Role,
					CreatedAt: user.CreatedAt,
				}, nil
			},
		}, logger.NewTestLogger())

		createdUser, err := usecase.Register(context.Background(), "  Admin@Example.com ", "password123", "ADMIN")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if createdUser.Email != "admin@example.com" {
			t.Fatalf("expected email admin@example.com, got %s", createdUser.Email)
		}
	})
}

func TestUserUsecaseAuthenticate(t *testing.T) {
	t.Run("returns invalid credentials for unknown user", func(t *testing.T) {
		usecase := NewUserUsecase(&fakeUserRepo{
			getByEmailFn: func(ctx context.Context, email string) (*User, error) {
				return nil, ErrUserNotFound
			},
		}, logger.NewTestLogger())

		_, err := usecase.Authenticate(context.Background(), "user@example.com", "password123")
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("returns invalid credentials for wrong password", func(t *testing.T) {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			t.Fatalf("generate hash: %v", err)
		}

		usecase := NewUserUsecase(&fakeUserRepo{
			getByEmailFn: func(ctx context.Context, email string) (*User, error) {
				return &User{
					Email:        email,
					PasswordHash: string(passwordHash),
					Role:         "user",
				}, nil
			},
		}, logger.NewTestLogger())

		_, err = usecase.Authenticate(context.Background(), "user@example.com", "wrong-password")
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("authenticates user", func(t *testing.T) {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			t.Fatalf("generate hash: %v", err)
		}

		usecase := NewUserUsecase(&fakeUserRepo{
			getByEmailFn: func(ctx context.Context, email string) (*User, error) {
				return &User{
					Email:        email,
					PasswordHash: string(passwordHash),
					Role:         "user",
				}, nil
			},
		}, logger.NewTestLogger())

		foundUser, err := usecase.Authenticate(context.Background(), "User@example.com", "password123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if foundUser.Email != "user@example.com" {
			t.Fatalf("expected normalized email, got %s", foundUser.Email)
		}
	})
}
