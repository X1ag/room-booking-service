package postgres

import (
	"context"
	"fmt"

	"test-backend-1-X1ag/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectPool creates and validates a PostgreSQL connection pool.
func ConnectPool(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	poolConfig, err := newPoolConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("build pgx pool config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping pgx pool: %w", err)
	}

	return pool, nil
}

func newPoolConfig(cfg config.DBConfig) (*pgxpool.Config, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse postgres DSN: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxOpenConns
	poolConfig.MinConns = cfg.MinIdleConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime

	return poolConfig, nil
}
