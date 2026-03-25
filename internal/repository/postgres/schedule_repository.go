package postgres

import (
	"context"
	"errors"
	"test-backend-1-X1ag/internal/schedule"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleRepository struct {
	pool *pgxpool.Pool
}

func NewScheduleRepository(pool *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{pool: pool}
}

func (r *ScheduleRepository) Create(ctx context.Context, schedule *schedule.Schedule) (*schedule.Schedule, error) {
	// Use transaction to ensure atomicity of schedule and schedule_days inserts
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			return
		}
	}()

	err = tx.QueryRow(ctx, `INSERT INTO schedules (id, room_id, start_time, end_time, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`, schedule.ID, schedule.RoomID, schedule.StartTime, schedule.EndTime, schedule.CreatedAt).Scan(&schedule.ID)
	if err != nil {
		return nil, err
	}
	for _, day := range schedule.DaysOfWeek {
		_, err := tx.Exec(ctx, `INSERT INTO schedule_days (schedule_id, day_of_week) VALUES ($1, $2)`, schedule.ID, day)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return schedule, err
}

func (r *ScheduleRepository) GetByRoomID(ctx context.Context, roomID uuid.UUID) (*schedule.Schedule, error) {
	var existing schedule.Schedule
	err := r.pool.QueryRow(ctx, `SELECT id, room_id, to_char(start_time, 'HH24:MI'), to_char(end_time, 'HH24:MI') FROM schedules WHERE room_id = $1`, roomID).Scan(&existing.ID, &existing.RoomID, &existing.StartTime, &existing.EndTime)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, schedule.ErrScheduleNotFound
		}
		return nil, err
	}
	rows, err := r.pool.Query(ctx, `SELECT day_of_week FROM schedule_days WHERE schedule_id = $1`, existing.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var daysOfWeek []int
	for rows.Next() {
		var day int
		err := rows.Scan(&day)
		if err != nil {
			return nil, err
		}
		daysOfWeek = append(daysOfWeek, day)
	}
	existing.DaysOfWeek = daysOfWeek

	return &existing, nil
}
