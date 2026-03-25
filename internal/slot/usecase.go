package slot

import (
	"context"
	"errors"
	"time"

	"test-backend-1-X1ag/internal/logger"
	"test-backend-1-X1ag/internal/room"
	"test-backend-1-X1ag/internal/schedule"

	"github.com/google/uuid"
)

type SlotUsecase struct {
	slotRepo     Repository
	roomRepo     room.Repository
	scheduleRepo schedule.Repository
	logger       *logger.ZerologLogger
}

func NewSlotUsecase(slotRepo Repository, roomRepo room.Repository, scheduleRepo schedule.Repository, logger *logger.ZerologLogger) *SlotUsecase {
	return &SlotUsecase{
		slotRepo:     slotRepo,
		roomRepo:     roomRepo,
		scheduleRepo: scheduleRepo,
		logger:       logger,
	}
}

func (u *SlotUsecase) GetByRoomID(ctx context.Context, roomID string, date string) ([]Slot, error) {
	if date == "" {
		return nil, ErrDateRequired
	}

	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		u.logger.Error().Err(err).Msg("invalid date format")
		return nil, ErrInvalidDate
	}

	roomUUID, err := uuid.Parse(roomID)
	if err != nil {
		u.logger.Error().Err(err).Msg("invalid room id format")
		return nil, room.ErrInvalidRoomID
	}

	_, err = u.roomRepo.GetByID(ctx, roomUUID) // check if room exists
	if err != nil {
		if errors.Is(err, room.ErrRoomNotFound) {
			u.logger.Error().Err(err).Msg("failed to get room")
			return nil, room.ErrRoomNotFound
		}
		u.logger.Error().Err(err).Msg("failed to get room")
		return nil, err
	}

	roomSchedule, err := u.scheduleRepo.GetByRoomID(ctx, roomUUID) // check schedule for this room
	if err != nil {
		if errors.Is(err, schedule.ErrScheduleNotFound) {
			u.logger.Error().Err(err).Msg("room schedule not found")
			return []Slot{}, nil
		}
		u.logger.Error().Err(err).Msg("failed to get schedule")
		return nil, err
	}

	if roomSchedule == nil {
		u.logger.Error().Msg("room schedule not found")
		return []Slot{}, nil
	}

	found := false
	for _, day := range roomSchedule.DaysOfWeek {
		weekDay := int(dateTime.Weekday())
		if weekDay == 0 {
			weekDay = 7 // convert Sunday from 0 to 7
		}
		if weekDay == int(day) { // time.Weekday starts from Sunday(0), our schedule starts from Monday(1)
			found = true
			break
		}
	}

	if !found {
		u.logger.Info().Msg("schedule does not apply to this date")
		return []Slot{}, nil
	}

	startClock, err := time.Parse("15:04", roomSchedule.StartTime)
	if err != nil {
		u.logger.Error().Err(err).Msg("failed to parse start time")
		return []Slot{}, err
	}
	startAt := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), startClock.Hour(), startClock.Minute(), 0, 0, time.UTC)

	endClock, err := time.Parse("15:04", roomSchedule.EndTime)
	if err != nil {
		u.logger.Error().Err(err).Msg("failed to parse end time")
		return []Slot{}, err
	}
	endAt := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), endClock.Hour(), endClock.Minute(), 0, 0, time.UTC)

	// iterate from startTime to endTime with step 30 minutes and create slots
	slotStart := startAt
	for slotStart.Before(endAt) {
		slotEnd := slotStart.Add(30 * time.Minute)
		if slotEnd.After(endAt) {
			break
		}
		err := u.slotRepo.CreateSlot(ctx, Slot{
			ID:        uuid.New(),
			RoomID:    roomUUID,
			StartTime: slotStart,
			EndTime:   slotEnd,
		})
		if err != nil {
			u.logger.Error().Err(err).Msg("failed to create slot")
			return nil, err
		}
		slotStart = slotEnd
	}

	availableSlots, err := u.slotRepo.GetByRoomID(ctx, roomUUID, startAt, endAt)
	if err != nil {
		u.logger.Error().Err(err).Str("room_id", roomID).Msg("Failed to get slots by room ID")
		return nil, err
	}

	u.logger.Info().Int("count", len(availableSlots)).Str("room_id", roomID).Msg("Slots retrieved successfully")
	return availableSlots, nil
}
