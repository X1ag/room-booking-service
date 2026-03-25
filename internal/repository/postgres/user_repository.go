package postgres

import (
	"context"
	"errors"

	"test-backend-1-X1ag/internal/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	query := `INSERT INTO users (id, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id, email, role, created_at`
	var createdUser user.User
	err := r.pool.QueryRow(ctx, query, u.ID, u.Email, u.PasswordHash, u.Role).Scan(&createdUser.ID, &createdUser.Email, &createdUser.Role, &createdUser.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
			return nil, user.ErrEmailAlreadyExists
		}
		return nil, err
	}
	return &createdUser, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `SELECT id, email, password_hash, role, created_at FROM users WHERE email = $1`
	var u user.User
	err := r.pool.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}
