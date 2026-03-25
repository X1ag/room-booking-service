package postgres

import (
	"context"
	"errors"
	"test-backend-1-X1ag/internal/slot"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SlotRepository struct {
	pool *pgxpool.Pool
}

func NewSlotRepository(pool *pgxpool.Pool) *SlotRepository {
	return &SlotRepository{pool: pool}
}

// GetByRoomID returns available slots for a given room and date
func (r *SlotRepository) GetByRoomID(ctx context.Context, roomID uuid.UUID, startDate time.Time, endDate time.Time) ([]slot.Slot, error) {
	query := `SELECT s.id, s.room_id, s.start_at, s.end_at 
		FROM slots s
		LEFT JOIN bookings b
			ON b.slot_id = s.id
		AND b.status = 'active'
		WHERE s.room_id = $1
			AND s.start_at >= $2
			AND s.start_at < $3
			AND b.id IS NULL
		ORDER BY s.start_at;
	`
	rows, err := r.pool.Query(ctx, query, roomID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []slot.Slot
	for rows.Next() {
		var s slot.Slot
		if err := rows.Scan(&s.ID, &s.RoomID, &s.StartTime, &s.EndTime); err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return slots, nil
}

func (r *SlotRepository) CreateSlot(ctx context.Context, slot slot.Slot) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO slots (id, room_id, start_at, end_at) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`, slot.ID, slot.RoomID, slot.StartTime, slot.EndTime)
	return err
}

func (r *SlotRepository) GetSlotByID(ctx context.Context, slotID uuid.UUID) (slot.Slot, error) {
	var s slot.Slot
	err := r.pool.QueryRow(ctx, `SELECT id, room_id, start_at, end_at FROM slots WHERE id = $1`, slotID).Scan(&s.ID, &s.RoomID, &s.StartTime, &s.EndTime)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return slot.Slot{}, slot.ErrSlotNotFound
		}
		return slot.Slot{}, err
	}
	return s, nil
}
