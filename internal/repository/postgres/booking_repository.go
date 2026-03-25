package postgres

import (
	"context"
	"errors"
	"time"

	"test-backend-1-X1ag/internal/booking"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository struct {
	pool *pgxpool.Pool
}

func NewBookingRepository(pool *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{pool: pool}
}

func (r *BookingRepository) Create(ctx context.Context, createdBooking booking.Booking) (booking.Booking, error) {
	query := `INSERT INTO bookings (id, user_id, slot_id, conference_link, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, slot_id, conference_link, status, created_at`
	err := r.pool.QueryRow(ctx, query, createdBooking.ID, createdBooking.UserID, createdBooking.SlotID, createdBooking.ConferenceLink, createdBooking.Status).
		Scan(&createdBooking.ID, &createdBooking.UserID, &createdBooking.SlotID, &createdBooking.ConferenceLink, &createdBooking.Status, &createdBooking.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "idx_bookings_active_slot" {
				return booking.Booking{}, booking.ErrSlotAlreadyBooked
			}
		}
		return booking.Booking{}, err
	}
	return createdBooking, nil
}

func (r *BookingRepository) GetBookingBySlotID(ctx context.Context, slotID uuid.UUID) (booking.Booking, error) {
	query := `SELECT id, user_id, slot_id, conference_link, status, created_at FROM bookings WHERE slot_id = $1 AND status = 'active'`
	var b booking.Booking
	err := r.pool.QueryRow(ctx, query, slotID).Scan(&b.ID, &b.UserID, &b.SlotID, &b.ConferenceLink, &b.Status, &b.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return booking.Booking{}, nil // no active booking for this slot
		}
		return booking.Booking{}, err
	}
	return b, nil
}

func (r *BookingRepository) GetByID(ctx context.Context, bookingID uuid.UUID) (booking.Booking, error) {
	query := `SELECT id, user_id, slot_id, conference_link, status, created_at FROM bookings WHERE id = $1`
	var b booking.Booking
	err := r.pool.QueryRow(ctx, query, bookingID).Scan(&b.ID, &b.UserID, &b.SlotID, &b.ConferenceLink, &b.Status, &b.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return booking.Booking{}, booking.ErrBookingNotFound
		}
		return booking.Booking{}, err
	}
	return b, nil
}

func (r *BookingRepository) Cancel(ctx context.Context, bookingID uuid.UUID) (booking.Booking, error) {
	query := `UPDATE bookings
		SET status = 'cancelled'
		WHERE id = $1
		RETURNING id, user_id, slot_id, conference_link, status, created_at`
	var b booking.Booking
	err := r.pool.QueryRow(ctx, query, bookingID).Scan(&b.ID, &b.UserID, &b.SlotID, &b.ConferenceLink, &b.Status, &b.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return booking.Booking{}, booking.ErrBookingNotFound
		}
		return booking.Booking{}, err
	}
	return b, nil
}

func (r *BookingRepository) List(ctx context.Context, page, pageSize int) ([]booking.Booking, int, error) {
	offset := (page - 1) * pageSize

	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM bookings`).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, user_id, slot_id, conference_link, status, created_at
		FROM bookings
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`
	rows, err := r.pool.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var bookings []booking.Booking
	for rows.Next() {
		var b booking.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.SlotID, &b.ConferenceLink, &b.Status, &b.CreatedAt); err != nil {
			return nil, 0, err
		}
		bookings = append(bookings, b)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
}

func (r *BookingRepository) ListFutureByUser(ctx context.Context, userID uuid.UUID, now time.Time) ([]booking.Booking, error) {
	query := `SELECT b.id, b.user_id, b.slot_id, b.conference_link, b.status, b.created_at
		FROM bookings b
		JOIN slots s ON s.id = b.slot_id
		WHERE b.user_id = $1
		  AND b.status = 'active'
		  AND s.start_at >= $2
		ORDER BY s.start_at ASC`
	rows, err := r.pool.Query(ctx, query, userID, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []booking.Booking
	for rows.Next() {
		var b booking.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.SlotID, &b.ConferenceLink, &b.Status, &b.CreatedAt); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}
