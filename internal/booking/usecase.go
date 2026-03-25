package booking

import (
	"context"
	"errors"
	"test-backend-1-X1ag/internal/conference"
	"test-backend-1-X1ag/internal/logger"
	"test-backend-1-X1ag/internal/slot"
	"time"

	"github.com/google/uuid"
)

type BookingUsecase struct {
	bookingRepo       Repository
	slotRepo          slot.Repository
	conferenceService conference.Service
	logger            *logger.ZerologLogger
}

func NewBookingUsecase(bookingRepo Repository, slotRepo slot.Repository, conferenceService conference.Service, logger *logger.ZerologLogger) *BookingUsecase {
	return &BookingUsecase{
		bookingRepo:       bookingRepo,
		slotRepo:          slotRepo,
		conferenceService: conferenceService,
		logger:            logger,
	}
}

func (u *BookingUsecase) Create(ctx context.Context, slotID string, createConferenceLink bool, userID string) (Booking, error) {
	slotUUID, err := uuid.Parse(slotID)
	if err != nil {
		u.logger.Error().Err(err).Str("slot_id", slotID).Msg("invalid slot id format")
		return Booking{}, ErrInvalidSlotID
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		u.logger.Error().Err(err).Str("user_id", userID).Msg("invalid user id format")
		return Booking{}, ErrInvalidUserID
	}

	slotExists, err := u.slotRepo.GetSlotByID(ctx, slotUUID)
	if err != nil {
		if errors.Is(err, slot.ErrSlotNotFound) {
			u.logger.Info().Str("slot_id", slotUUID.String()).Msg("slot not found")
			return Booking{}, slot.ErrSlotNotFound
		}
		u.logger.Error().Err(err).Str("slot_id", slotUUID.String()).Msg("failed to get slot by id")
		return Booking{}, err
	}
	if slotExists.StartTime.Before(time.Now().UTC()) {
		u.logger.Info().Str("slot_id", slotUUID.String()).Msg("slot is in the past")
		return Booking{}, ErrSlotInPast
	}

	bookingExists, err := u.bookingRepo.GetBookingBySlotID(ctx, slotUUID)
	if err != nil {
		u.logger.Error().Err(err).Str("slot_id", slotUUID.String()).Msg("failed to get booking by slot id")
		return Booking{}, err
	}
	if bookingExists != (Booking{}) {
		u.logger.Info().Str("slot_id", slotUUID.String()).Msg("slot is already booked")
		return Booking{}, ErrSlotAlreadyBooked
	}

	var link *string
	if createConferenceLink {
		if u.conferenceService == nil {
			u.logger.Error().Str("slot_id", slotUUID.String()).Msg("conference service is not configured")
			return Booking{}, ErrConferenceUnavailable
		}

		generated, err := u.conferenceService.CreateLink(ctx, slotUUID, userUUID)
		if err != nil {
			u.logger.Error().Err(err).Str("slot_id", slotUUID.String()).Str("user_id", userUUID.String()).Msg("failed to create conference link")
			return Booking{}, ErrConferenceUnavailable
		}
		link = &generated
	}

	bookingToCreate := Booking{
		ID:             uuid.New(),
		UserID:         userUUID,
		SlotID:         slotUUID,
		ConferenceLink: link,
		Status:         "active",
	}

	createdBooking, err := u.bookingRepo.Create(ctx, bookingToCreate)
	if err != nil {
		u.logger.Error().Err(err).Str("slot_id", slotUUID.String()).Str("user_id", userUUID.String()).Msg("failed to create booking")
		return Booking{}, err
	}

	u.logger.Info().Str("booking_id", createdBooking.ID.String()).Str("slot_id", slotUUID.String()).Str("user_id", userUUID.String()).Msg("booking created successfully")
	return createdBooking, nil
}

func (u *BookingUsecase) Cancel(ctx context.Context, bookingID string, userID string) (Booking, error) {
	bookingUUID, err := uuid.Parse(bookingID)
	if err != nil {
		u.logger.Error().Err(err).Str("booking_id", bookingID).Msg("invalid booking id format")
		return Booking{}, ErrInvalidBookingID
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		u.logger.Error().Err(err).Str("user_id", userID).Msg("invalid user id format")
		return Booking{}, ErrInvalidUserID
	}

	existingBooking, err := u.bookingRepo.GetByID(ctx, bookingUUID)
	if err != nil {
		return Booking{}, err
	}

	if existingBooking.UserID != userUUID {
		return Booking{}, ErrForbidden
	}

	cancelledBooking, err := u.bookingRepo.Cancel(ctx, bookingUUID)
	if err != nil {
		return Booking{}, err
	}

	return cancelledBooking, nil
}

func (u *BookingUsecase) List(ctx context.Context, page, pageSize int) ([]Booking, int, error) {
	if page < 1 || pageSize < 1 || pageSize > 100 {
		return nil, 0, ErrInvalidPagination
	}

	return u.bookingRepo.List(ctx, page, pageSize)
}

func (u *BookingUsecase) GetUserBookings(ctx context.Context, userID string) ([]Booking, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		u.logger.Error().Err(err).Str("user_id", userID).Msg("invalid user id format")
		return nil, ErrInvalidUserID
	}

	return u.bookingRepo.ListFutureByUser(ctx, userUUID, time.Now().UTC())
}
