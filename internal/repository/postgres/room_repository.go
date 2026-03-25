package postgres

import (
	"context"
	"errors"

	"test-backend-1-X1ag/internal/room"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomRepository struct {
	pool *pgxpool.Pool
}

func NewRoomRepository(pool *pgxpool.Pool) *RoomRepository {
	return &RoomRepository{pool: pool}
}

func (r *RoomRepository) Create(ctx context.Context, room *room.Room) (roomID uuid.UUID, err error) {
	query := `INSERT INTO rooms (id, name, description, capacity, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = r.pool.QueryRow(ctx, query, room.ID, room.Name, room.Description, room.Capacity, room.CreatedAt).Scan(&room.ID)
	return room.ID, err
}

func (r *RoomRepository) GetAll(ctx context.Context) ([]room.Room, error) {
	query := `SELECT id, name, description, capacity, created_at FROM rooms ORDER BY created_at DESC, id`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []room.Room
	for rows.Next() {
		var rm room.Room
		err := rows.Scan(&rm.ID, &rm.Name, &rm.Description, &rm.Capacity, &rm.CreatedAt)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, rm)
	}

	return rooms, rows.Err()
}

func (r *RoomRepository) GetByID(ctx context.Context, id uuid.UUID) (*room.Room, error) {
	query := `SELECT id, name, description, capacity, created_at FROM rooms WHERE id = $1`
	var rm room.Room
	err := r.pool.QueryRow(ctx, query, id).Scan(&rm.ID, &rm.Name, &rm.Description, &rm.Capacity, &rm.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, room.ErrRoomNotFound
		}
		return nil, err
	}
	return &rm, nil
}
