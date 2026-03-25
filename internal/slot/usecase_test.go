package slot

import (
	"context"
	"errors"
	"testing"
	"time"

	"test-backend-1-X1ag/internal/logger"
	"test-backend-1-X1ag/internal/room"
	"test-backend-1-X1ag/internal/schedule"

	"github.com/google/uuid"
)

type fakeSlotRepo struct {
	getByRoomIDFn func(ctx context.Context, roomID uuid.UUID, startDate time.Time, endDate time.Time) ([]Slot, error)
	createSlotFn  func(ctx context.Context, slot Slot) error
	getSlotByIDFn func(ctx context.Context, slotID uuid.UUID) (Slot, error)
}

func (f *fakeSlotRepo) GetByRoomID(ctx context.Context, roomID uuid.UUID, startDate time.Time, endDate time.Time) ([]Slot, error) {
	return f.getByRoomIDFn(ctx, roomID, startDate, endDate)
}

func (f *fakeSlotRepo) CreateSlot(ctx context.Context, slot Slot) error {
	return f.createSlotFn(ctx, slot)
}

func (f *fakeSlotRepo) GetSlotByID(ctx context.Context, slotID uuid.UUID) (Slot, error) {
	return f.getSlotByIDFn(ctx, slotID)
}

type fakeSlotRoomRepo struct {
	getByIDFn func(ctx context.Context, id uuid.UUID) (*room.Room, error)
}

func (f *fakeSlotRoomRepo) Create(ctx context.Context, room *room.Room) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (f *fakeSlotRoomRepo) GetAll(ctx context.Context) ([]room.Room, error) {
	return nil, nil
}

func (f *fakeSlotRoomRepo) GetByID(ctx context.Context, id uuid.UUID) (*room.Room, error) {
	return f.getByIDFn(ctx, id)
}

type fakeSlotScheduleRepo struct {
	getByRoomIDFn func(ctx context.Context, roomID uuid.UUID) (*schedule.Schedule, error)
}

func (f *fakeSlotScheduleRepo) Create(ctx context.Context, schedule *schedule.Schedule) (*schedule.Schedule, error) {
	return nil, nil
}

func (f *fakeSlotScheduleRepo) GetByRoomID(ctx context.Context, roomID uuid.UUID) (*schedule.Schedule, error) {
	return f.getByRoomIDFn(ctx, roomID)
}

func TestSlotUsecaseGetByRoomID(t *testing.T) {
	t.Run("returns date required", func(t *testing.T) {
		usecase := NewSlotUsecase(&fakeSlotRepo{}, &fakeSlotRoomRepo{}, &fakeSlotScheduleRepo{}, logger.NewTestLogger())

		_, err := usecase.GetByRoomID(context.Background(), uuid.NewString(), "")
		if !errors.Is(err, ErrDateRequired) {
			t.Fatalf("expected ErrDateRequired, got %v", err)
		}
	})

	t.Run("returns invalid date", func(t *testing.T) {
		usecase := NewSlotUsecase(&fakeSlotRepo{}, &fakeSlotRoomRepo{}, &fakeSlotScheduleRepo{}, logger.NewTestLogger())

		_, err := usecase.GetByRoomID(context.Background(), uuid.NewString(), "bad-date")
		if !errors.Is(err, ErrInvalidDate) {
			t.Fatalf("expected ErrInvalidDate, got %v", err)
		}
	})

	t.Run("returns room not found", func(t *testing.T) {
		usecase := NewSlotUsecase(
			&fakeSlotRepo{},
			&fakeSlotRoomRepo{
				getByIDFn: func(ctx context.Context, id uuid.UUID) (*room.Room, error) {
					return nil, room.ErrRoomNotFound
				},
			},
			&fakeSlotScheduleRepo{},
			logger.NewTestLogger(),
		)

		_, err := usecase.GetByRoomID(context.Background(), uuid.NewString(), "2026-03-25")
		if !errors.Is(err, room.ErrRoomNotFound) {
			t.Fatalf("expected room.ErrRoomNotFound, got %v", err)
		}
	})

	t.Run("returns empty list when schedule not found", func(t *testing.T) {
		usecase := NewSlotUsecase(
			&fakeSlotRepo{},
			&fakeSlotRoomRepo{
				getByIDFn: func(ctx context.Context, id uuid.UUID) (*room.Room, error) {
					return &room.Room{ID: id}, nil
				},
			},
			&fakeSlotScheduleRepo{
				getByRoomIDFn: func(ctx context.Context, roomID uuid.UUID) (*schedule.Schedule, error) {
					return nil, schedule.ErrScheduleNotFound
				},
			},
			logger.NewTestLogger(),
		)

		slots, err := usecase.GetByRoomID(context.Background(), uuid.NewString(), "2026-03-25")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(slots) != 0 {
			t.Fatalf("expected empty slots, got %d", len(slots))
		}
	})

	t.Run("creates slots and returns available ones", func(t *testing.T) {
		roomID := uuid.New()
		var createdCount int
		repo := &fakeSlotRepo{
			createSlotFn: func(ctx context.Context, slot Slot) error {
				createdCount++
				return nil
			},
			getByRoomIDFn: func(ctx context.Context, roomID uuid.UUID, startDate time.Time, endDate time.Time) ([]Slot, error) {
				return []Slot{{ID: uuid.New(), RoomID: roomID, StartTime: startDate, EndTime: startDate.Add(30 * time.Minute)}}, nil
			},
		}
		usecase := NewSlotUsecase(
			repo,
			&fakeSlotRoomRepo{
				getByIDFn: func(ctx context.Context, id uuid.UUID) (*room.Room, error) {
					return &room.Room{ID: id}, nil
				},
			},
			&fakeSlotScheduleRepo{
				getByRoomIDFn: func(ctx context.Context, roomID uuid.UUID) (*schedule.Schedule, error) {
					return &schedule.Schedule{
						ID:         uuid.New(),
						RoomID:     roomID,
						DaysOfWeek: []int{3},
						StartTime:  "09:00",
						EndTime:    "10:00",
					}, nil
				},
			},
			logger.NewTestLogger(),
		)

		slots, err := usecase.GetByRoomID(context.Background(), roomID.String(), "2026-03-25")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if createdCount != 2 {
			t.Fatalf("expected 2 created slots, got %d", createdCount)
		}
		if len(slots) != 1 {
			t.Fatalf("expected 1 available slot, got %d", len(slots))
		}
	})
}
