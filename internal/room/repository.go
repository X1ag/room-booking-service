package room

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, room *Room) (roomID uuid.UUID, err error)
	GetAll(ctx context.Context) ([]Room, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Room, error)
}
