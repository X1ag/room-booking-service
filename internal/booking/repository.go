package booking

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, booking Booking) (Booking, error)
	GetBookingBySlotID(ctx context.Context, slotID uuid.UUID) (Booking, error)
	GetByID(ctx context.Context, bookingID uuid.UUID) (Booking, error)
	Cancel(ctx context.Context, bookingID uuid.UUID) (Booking, error)
	List(ctx context.Context, page, pageSize int) ([]Booking, int, error)
	ListFutureByUser(ctx context.Context, userID uuid.UUID, now time.Time) ([]Booking, error)
}
