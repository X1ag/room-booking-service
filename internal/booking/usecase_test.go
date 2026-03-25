package booking

import (
	"context"
	"errors"
	"testing"
	"time"

	"test-backend-1-X1ag/internal/conference"
	"test-backend-1-X1ag/internal/logger"
	"test-backend-1-X1ag/internal/slot"

	"github.com/google/uuid"
)

type fakeBookingRepo struct {
	createFn           func(ctx context.Context, booking Booking) (Booking, error)
	getBookingBySlotFn func(ctx context.Context, slotID uuid.UUID) (Booking, error)
	getByIDFn          func(ctx context.Context, bookingID uuid.UUID) (Booking, error)
	cancelFn           func(ctx context.Context, bookingID uuid.UUID) (Booking, error)
	listFn             func(ctx context.Context, page, pageSize int) ([]Booking, int, error)
	listFutureByUserFn func(ctx context.Context, userID uuid.UUID, now time.Time) ([]Booking, error)
}

func (f *fakeBookingRepo) Create(ctx context.Context, booking Booking) (Booking, error) {
	return f.createFn(ctx, booking)
}

func (f *fakeBookingRepo) GetBookingBySlotID(ctx context.Context, slotID uuid.UUID) (Booking, error) {
	return f.getBookingBySlotFn(ctx, slotID)
}

func (f *fakeBookingRepo) GetByID(ctx context.Context, bookingID uuid.UUID) (Booking, error) {
	return f.getByIDFn(ctx, bookingID)
}

func (f *fakeBookingRepo) Cancel(ctx context.Context, bookingID uuid.UUID) (Booking, error) {
	return f.cancelFn(ctx, bookingID)
}

func (f *fakeBookingRepo) List(ctx context.Context, page, pageSize int) ([]Booking, int, error) {
	return f.listFn(ctx, page, pageSize)
}

func (f *fakeBookingRepo) ListFutureByUser(ctx context.Context, userID uuid.UUID, now time.Time) ([]Booking, error) {
	return f.listFutureByUserFn(ctx, userID, now)
}

type fakeBookingSlotRepo struct {
	getSlotByIDFn func(ctx context.Context, slotID uuid.UUID) (slot.Slot, error)
}

type fakeConferenceService struct {
	createLinkFn func(ctx context.Context, slotID uuid.UUID, userID uuid.UUID) (string, error)
}

func (f *fakeConferenceService) CreateLink(ctx context.Context, slotID uuid.UUID, userID uuid.UUID) (string, error) {
	return f.createLinkFn(ctx, slotID, userID)
}

func (f *fakeBookingSlotRepo) GetByRoomID(ctx context.Context, roomID uuid.UUID, startDate time.Time, endDate time.Time) ([]slot.Slot, error) {
	return nil, nil
}

func (f *fakeBookingSlotRepo) CreateSlot(ctx context.Context, slot slot.Slot) error {
	return nil
}

func (f *fakeBookingSlotRepo) GetSlotByID(ctx context.Context, slotID uuid.UUID) (slot.Slot, error) {
	return f.getSlotByIDFn(ctx, slotID)
}

func TestBookingUsecaseCreate(t *testing.T) {
	t.Run("returns invalid slot id", func(t *testing.T) {
		usecase := NewBookingUsecase(&fakeBookingRepo{}, &fakeBookingSlotRepo{}, conference.NewMockService(), logger.NewTestLogger())

		_, err := usecase.Create(context.Background(), "bad-id", false, uuid.NewString())
		if !errors.Is(err, ErrInvalidSlotID) {
			t.Fatalf("expected ErrInvalidSlotID, got %v", err)
		}
	})

	t.Run("returns invalid user id", func(t *testing.T) {
		usecase := NewBookingUsecase(&fakeBookingRepo{}, &fakeBookingSlotRepo{}, conference.NewMockService(), logger.NewTestLogger())

		_, err := usecase.Create(context.Background(), uuid.NewString(), false, "bad-id")
		if !errors.Is(err, ErrInvalidUserID) {
			t.Fatalf("expected ErrInvalidUserID, got %v", err)
		}
	})

	t.Run("returns slot not found", func(t *testing.T) {
		usecase := NewBookingUsecase(
			&fakeBookingRepo{},
			&fakeBookingSlotRepo{
				getSlotByIDFn: func(ctx context.Context, slotID uuid.UUID) (slot.Slot, error) {
					return slot.Slot{}, slot.ErrSlotNotFound
				},
			},
			conference.NewMockService(),
			logger.NewTestLogger(),
		)

		_, err := usecase.Create(context.Background(), uuid.NewString(), false, uuid.NewString())
		if !errors.Is(err, slot.ErrSlotNotFound) {
			t.Fatalf("expected slot.ErrSlotNotFound, got %v", err)
		}
	})

	t.Run("returns slot in past", func(t *testing.T) {
		usecase := NewBookingUsecase(
			&fakeBookingRepo{},
			&fakeBookingSlotRepo{
				getSlotByIDFn: func(ctx context.Context, slotID uuid.UUID) (slot.Slot, error) {
					return slot.Slot{ID: slotID, StartTime: time.Now().UTC().Add(-time.Hour)}, nil
				},
			},
			conference.NewMockService(),
			logger.NewTestLogger(),
		)

		_, err := usecase.Create(context.Background(), uuid.NewString(), false, uuid.NewString())
		if !errors.Is(err, ErrSlotInPast) {
			t.Fatalf("expected ErrSlotInPast, got %v", err)
		}
	})

	t.Run("returns already booked", func(t *testing.T) {
		usecase := NewBookingUsecase(
			&fakeBookingRepo{
				getBookingBySlotFn: func(ctx context.Context, slotID uuid.UUID) (Booking, error) {
					return Booking{ID: uuid.New(), SlotID: slotID}, nil
				},
			},
			&fakeBookingSlotRepo{
				getSlotByIDFn: func(ctx context.Context, slotID uuid.UUID) (slot.Slot, error) {
					return slot.Slot{ID: slotID, StartTime: time.Now().UTC().Add(time.Hour)}, nil
				},
			},
			conference.NewMockService(),
			logger.NewTestLogger(),
		)

		_, err := usecase.Create(context.Background(), uuid.NewString(), false, uuid.NewString())
		if !errors.Is(err, ErrSlotAlreadyBooked) {
			t.Fatalf("expected ErrSlotAlreadyBooked, got %v", err)
		}
	})

	t.Run("creates booking with conference link", func(t *testing.T) {
		repo := &fakeBookingRepo{
			getBookingBySlotFn: func(ctx context.Context, slotID uuid.UUID) (Booking, error) {
				return Booking{}, nil
			},
			createFn: func(ctx context.Context, booking Booking) (Booking, error) {
				if booking.ID == uuid.Nil {
					t.Fatal("expected generated booking id")
				}
				if booking.ConferenceLink == nil || *booking.ConferenceLink == "" {
					t.Fatal("expected conference link")
				}
				return booking, nil
			},
		}
		usecase := NewBookingUsecase(
			repo,
			&fakeBookingSlotRepo{
				getSlotByIDFn: func(ctx context.Context, slotID uuid.UUID) (slot.Slot, error) {
					return slot.Slot{ID: slotID, StartTime: time.Now().UTC().Add(time.Hour)}, nil
				},
			},
			&fakeConferenceService{
				createLinkFn: func(ctx context.Context, slotID uuid.UUID, userID uuid.UUID) (string, error) {
					return "https://meet.example.com/test-link", nil
				},
			},
			logger.NewTestLogger(),
		)

		createdBooking, err := usecase.Create(context.Background(), uuid.NewString(), true, uuid.NewString())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if createdBooking.Status != "active" {
			t.Fatalf("expected status active, got %s", createdBooking.Status)
		}
	})

	t.Run("returns conference service unavailable", func(t *testing.T) {
		usecase := NewBookingUsecase(
			&fakeBookingRepo{
				getBookingBySlotFn: func(ctx context.Context, slotID uuid.UUID) (Booking, error) {
					return Booking{}, nil
				},
			},
			&fakeBookingSlotRepo{
				getSlotByIDFn: func(ctx context.Context, slotID uuid.UUID) (slot.Slot, error) {
					return slot.Slot{ID: slotID, StartTime: time.Now().UTC().Add(time.Hour)}, nil
				},
			},
			&fakeConferenceService{
				createLinkFn: func(ctx context.Context, slotID uuid.UUID, userID uuid.UUID) (string, error) {
					return "", errors.New("conference service down")
				},
			},
			logger.NewTestLogger(),
		)

		_, err := usecase.Create(context.Background(), uuid.NewString(), true, uuid.NewString())
		if !errors.Is(err, ErrConferenceUnavailable) {
			t.Fatalf("expected ErrConferenceUnavailable, got %v", err)
		}
	})
}

func TestBookingUsecaseCancelAndList(t *testing.T) {
	t.Run("returns forbidden for чужая бронь", func(t *testing.T) {
		ownerID := uuid.New()
		usecase := NewBookingUsecase(
			&fakeBookingRepo{
				getByIDFn: func(ctx context.Context, bookingID uuid.UUID) (Booking, error) {
					return Booking{ID: bookingID, UserID: ownerID}, nil
				},
			},
			&fakeBookingSlotRepo{},
			conference.NewMockService(),
			logger.NewTestLogger(),
		)

		_, err := usecase.Cancel(context.Background(), uuid.NewString(), uuid.NewString())
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("cancels booking", func(t *testing.T) {
		userID := uuid.New()
		repo := &fakeBookingRepo{
			getByIDFn: func(ctx context.Context, bookingID uuid.UUID) (Booking, error) {
				return Booking{ID: bookingID, UserID: userID}, nil
			},
			cancelFn: func(ctx context.Context, bookingID uuid.UUID) (Booking, error) {
				return Booking{ID: bookingID, UserID: userID, Status: "cancelled"}, nil
			},
		}
		usecase := NewBookingUsecase(repo, &fakeBookingSlotRepo{}, conference.NewMockService(), logger.NewTestLogger())

		cancelledBooking, err := usecase.Cancel(context.Background(), uuid.NewString(), userID.String())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cancelledBooking.Status != "cancelled" {
			t.Fatalf("expected status cancelled, got %s", cancelledBooking.Status)
		}
	})

	t.Run("returns invalid pagination", func(t *testing.T) {
		usecase := NewBookingUsecase(&fakeBookingRepo{}, &fakeBookingSlotRepo{}, conference.NewMockService(), logger.NewTestLogger())

		_, _, err := usecase.List(context.Background(), 0, 20)
		if !errors.Is(err, ErrInvalidPagination) {
			t.Fatalf("expected ErrInvalidPagination, got %v", err)
		}
	})

	t.Run("returns invalid user id for user bookings", func(t *testing.T) {
		usecase := NewBookingUsecase(&fakeBookingRepo{}, &fakeBookingSlotRepo{}, conference.NewMockService(), logger.NewTestLogger())

		_, err := usecase.GetUserBookings(context.Background(), "bad-id")
		if !errors.Is(err, ErrInvalidUserID) {
			t.Fatalf("expected ErrInvalidUserID, got %v", err)
		}
	})
}
