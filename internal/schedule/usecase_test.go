package schedule

import (
	"context"
	"errors"
	"testing"

	"test-backend-1-X1ag/internal/logger"
	"test-backend-1-X1ag/internal/room"

	"github.com/google/uuid"
)

type fakeScheduleRepo struct {
	createFn func(ctx context.Context, schedule *Schedule) (*Schedule, error)
	getByRoomIDFn func(ctx context.Context, roomID uuid.UUID) (*Schedule, error)
}

func (f *fakeScheduleRepo) Create(ctx context.Context, schedule *Schedule) (*Schedule, error) {
	return f.createFn(ctx, schedule)
}

func (f *fakeScheduleRepo) GetByRoomID(ctx context.Context, roomID uuid.UUID) (*Schedule, error) {
	return f.getByRoomIDFn(ctx, roomID)
}

type fakeScheduleRoomRepo struct {
	getByIDFn func(ctx context.Context, id uuid.UUID) (*room.Room, error)
}

func (f *fakeScheduleRoomRepo) Create(ctx context.Context, room *room.Room) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (f *fakeScheduleRoomRepo) GetAll(ctx context.Context) ([]room.Room, error) {
	return nil, nil
}

func (f *fakeScheduleRoomRepo) GetByID(ctx context.Context, id uuid.UUID) (*room.Room, error) {
	return f.getByIDFn(ctx, id)
}

func TestScheduleUsecaseCreateValidation(t *testing.T) {
	usecase := NewSheduleUsecase(&fakeScheduleRepo{}, &fakeScheduleRoomRepo{}, logger.NewTestLogger())

	testCases := []struct {
		name       string
		roomID     string
		daysOfWeek []int
		startTime  string
		endTime    string
		expected   error
	}{
		{"invalid room id", "bad-id", []int{1}, "09:00", "10:00", ErrInvalidRoomID},
		{"invalid start time", uuid.NewString(), []int{1}, "bad", "10:00", ErrInvalidTime},
		{"invalid end time", uuid.NewString(), []int{1}, "09:00", "bad", ErrInvalidTime},
		{"start equals end", uuid.NewString(), []int{1}, "09:00", "09:00", ErrStartTimeAfterEndTime},
		{"empty days", uuid.NewString(), []int{}, "09:00", "10:00", ErrInvalidDaysOfWeek},
		{"duplicate days", uuid.NewString(), []int{1, 1}, "09:00", "10:00", ErrInvalidDaysOfWeek},
		{"day out of range", uuid.NewString(), []int{8}, "09:00", "10:00", ErrInvalidDaysOfWeek},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := usecase.Create(context.Background(), tc.roomID, tc.daysOfWeek, tc.startTime, tc.endTime)
			if !errors.Is(err, tc.expected) {
				t.Fatalf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestScheduleUsecaseCreate(t *testing.T) {
	t.Run("returns room not found", func(t *testing.T) {
		usecase := NewSheduleUsecase(
			&fakeScheduleRepo{},
			&fakeScheduleRoomRepo{
				getByIDFn: func(ctx context.Context, id uuid.UUID) (*room.Room, error) {
					return nil, room.ErrRoomNotFound
				},
			},
			logger.NewTestLogger(),
		)

		_, err := usecase.Create(context.Background(), uuid.NewString(), []int{1}, "09:00", "10:00")
		if !errors.Is(err, ErrRoomNotFound) {
			t.Fatalf("expected ErrRoomNotFound, got %v", err)
		}
	})

	t.Run("returns schedule already exists", func(t *testing.T) {
		roomID := uuid.New()
		usecase := NewSheduleUsecase(
			&fakeScheduleRepo{
				getByRoomIDFn: func(ctx context.Context, roomID uuid.UUID) (*Schedule, error) {
					return &Schedule{ID: uuid.New(), RoomID: roomID}, nil
				},
			},
			&fakeScheduleRoomRepo{
				getByIDFn: func(ctx context.Context, id uuid.UUID) (*room.Room, error) {
					return &room.Room{ID: id}, nil
				},
			},
			logger.NewTestLogger(),
		)

		_, err := usecase.Create(context.Background(), roomID.String(), []int{1}, "09:00", "10:00")
		if !errors.Is(err, ErrScheduleAlreadyExists) {
			t.Fatalf("expected ErrScheduleAlreadyExists, got %v", err)
		}
	})

	t.Run("creates schedule", func(t *testing.T) {
		roomID := uuid.New()
		repo := &fakeScheduleRepo{
			getByRoomIDFn: func(ctx context.Context, roomID uuid.UUID) (*Schedule, error) {
				return nil, ErrScheduleNotFound
			},
			createFn: func(ctx context.Context, schedule *Schedule) (*Schedule, error) {
				if schedule.ID == uuid.Nil {
					t.Fatal("expected generated schedule id")
				}
				return schedule, nil
			},
		}
		usecase := NewSheduleUsecase(
			repo,
			&fakeScheduleRoomRepo{
				getByIDFn: func(ctx context.Context, id uuid.UUID) (*room.Room, error) {
					return &room.Room{ID: id}, nil
				},
			},
			logger.NewTestLogger(),
		)

		createdSchedule, err := usecase.Create(context.Background(), roomID.String(), []int{1, 2}, "09:00", "10:00")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if createdSchedule.RoomID != roomID {
			t.Fatalf("expected room id %s, got %s", roomID, createdSchedule.RoomID)
		}
	})
}

