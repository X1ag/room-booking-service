package schedule

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, schedule *Schedule) (*Schedule, error)
	GetByRoomID(ctx context.Context, roomID uuid.UUID) (*Schedule, error)
}
