package slot

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	GetByRoomID(ctx context.Context, roomID uuid.UUID, startDate time.Time, endDate time.Time) ([]Slot, error)
	CreateSlot(ctx context.Context, slot Slot) error
	GetSlotByID(ctx context.Context, slotID uuid.UUID) (Slot, error)
}