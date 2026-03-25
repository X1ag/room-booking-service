package room

import (
	"context"
	"errors"
	"testing"

	"test-backend-1-X1ag/internal/logger"

	"github.com/google/uuid"
)

type fakeRoomRepo struct {
	createFn func(ctx context.Context, room *Room) (uuid.UUID, error)
	getAllFn func(ctx context.Context) ([]Room, error)
	getByIDFn func(ctx context.Context, id uuid.UUID) (*Room, error)
}

func (f *fakeRoomRepo) Create(ctx context.Context, room *Room) (uuid.UUID, error) {
	return f.createFn(ctx, room)
}

func (f *fakeRoomRepo) GetAll(ctx context.Context) ([]Room, error) {
	return f.getAllFn(ctx)
}

func (f *fakeRoomRepo) GetByID(ctx context.Context, id uuid.UUID) (*Room, error) {
	return f.getByIDFn(ctx, id)
}

func TestRoomUsecaseCreate(t *testing.T) {
	t.Run("returns error for blank name", func(t *testing.T) {
		usecase := NewRoomUsecase(&fakeRoomRepo{}, logger.NewTestLogger())

		_, err := usecase.Create(context.Background(), "   ", nil, nil)
		if !errors.Is(err, ErrInvalidName) {
			t.Fatalf("expected ErrInvalidName, got %v", err)
		}
	})

	t.Run("returns error for invalid capacity", func(t *testing.T) {
		usecase := NewRoomUsecase(&fakeRoomRepo{}, logger.NewTestLogger())
		capacity := 0

		_, err := usecase.Create(context.Background(), "Room A", nil, &capacity)
		if !errors.Is(err, ErrInvalidCapacity) {
			t.Fatalf("expected ErrInvalidCapacity, got %v", err)
		}
	})

	t.Run("creates room and trims empty description", func(t *testing.T) {
		repo := &fakeRoomRepo{
			createFn: func(ctx context.Context, room *Room) (uuid.UUID, error) {
				if room.Description != nil {
					t.Fatalf("expected nil description, got %v", *room.Description)
				}
				if room.ID == uuid.Nil {
					t.Fatal("expected generated room id")
				}
				return room.ID, nil
			},
		}
		usecase := NewRoomUsecase(repo, logger.NewTestLogger())
		description := "   "
		capacity := 6

		createdRoom, err := usecase.Create(context.Background(), "Room A", &description, &capacity)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if createdRoom.ID == uuid.Nil {
			t.Fatal("expected room id to be set")
		}
		if createdRoom.Description != nil {
			t.Fatalf("expected nil description, got %v", *createdRoom.Description)
		}
	})
}

func TestRoomUsecaseGetByID(t *testing.T) {
	t.Run("returns invalid room id", func(t *testing.T) {
		usecase := NewRoomUsecase(&fakeRoomRepo{}, logger.NewTestLogger())

		_, err := usecase.GetByID(context.Background(), "bad-id")
		if !errors.Is(err, ErrInvalidRoomID) {
			t.Fatalf("expected ErrInvalidRoomID, got %v", err)
		}
	})

	t.Run("returns room from repository", func(t *testing.T) {
		roomID := uuid.New()
		repo := &fakeRoomRepo{
			getByIDFn: func(ctx context.Context, id uuid.UUID) (*Room, error) {
				return &Room{ID: id, Name: "Room A"}, nil
			},
		}
		usecase := NewRoomUsecase(repo, logger.NewTestLogger())

		foundRoom, err := usecase.GetByID(context.Background(), roomID.String())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if foundRoom.ID != roomID {
			t.Fatalf("expected room id %s, got %s", roomID, foundRoom.ID)
		}
	})
}

