package schedule

import (
	"context"
	"errors"
	"test-backend-1-X1ag/internal/logger"
	"test-backend-1-X1ag/internal/room"
	"time"

	"github.com/google/uuid"
)

type ScheduleUsecase struct {
	scheduleRepo Repository
	roomRepo     room.Repository
	logger       *logger.ZerologLogger
}

func NewSheduleUsecase(scheduleRepo Repository, roomRepo room.Repository, logger *logger.ZerologLogger) *ScheduleUsecase {
	return &ScheduleUsecase{
		scheduleRepo: scheduleRepo,
		roomRepo:     roomRepo,
		logger:       logger,
	}
}

func (u *ScheduleUsecase) Create(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*Schedule, error) {
	roomUUID, err := uuid.Parse(roomID)
	if err != nil {
		u.logger.Error().Err(err).Msg("invalid room id format")
		return nil, ErrInvalidRoomID
	}
	start, err := time.Parse("15:04", startTime)
	if err != nil {
		u.logger.Error().Err(err).Msg("invalid start time format")
		return nil, ErrInvalidTime
	}

	end, err := time.Parse("15:04", endTime)
	if err != nil {
		u.logger.Error().Err(err).Msg("invalid end time format")
		return nil, ErrInvalidTime
	}

	if time.Time(start).After(time.Time(end)) || time.Time(start).Equal(time.Time(end)) {
		u.logger.Error().Msg("start time must be before end time")
		return nil, ErrStartTimeAfterEndTime
	}

	if len(daysOfWeek) == 0 {
		u.logger.Error().Int("count", len(daysOfWeek)).Msg("days of week must not be empty")
		return nil, ErrInvalidDaysOfWeek
	}

	uniqueDays := make(map[int]struct{})
	for _, day := range daysOfWeek {
		uniqueDays[day] = struct{}{}
	}
	if len(uniqueDays) != len(daysOfWeek) {
		u.logger.Error().Int("count", len(daysOfWeek)).Msg("days of week must be unique")
		return nil, ErrInvalidDaysOfWeek
	}
	for _, day := range daysOfWeek {
		if day <= 0 || day > 7 {
			u.logger.Error().Int("day", day).Msg("invalid day of week")
			return nil, ErrInvalidDaysOfWeek
		}
	}

	_, err = u.roomRepo.GetByID(ctx, roomUUID)
	if err != nil {
		if errors.Is(err, room.ErrRoomNotFound) {
			u.logger.Error().Err(err).Msg("failed to get room")
			return nil, ErrRoomNotFound
		}
		u.logger.Error().Err(err).Msg("failed to get room")
		return nil, err
	}

	existing, err := u.scheduleRepo.GetByRoomID(ctx, roomUUID)
	if err != nil {
		if !errors.Is(err, ErrScheduleNotFound) {
			u.logger.Error().Err(err).Msg("failed to get schedule by room id")
			return nil, err
		}
	}
	if existing != nil {
		u.logger.Error().Str("room_id", roomID).Msg("schedule already exists for room")
		return nil, ErrScheduleAlreadyExists
	}

	schedule := &Schedule{
		ID:         uuid.New(),
		RoomID:     roomUUID,
		DaysOfWeek: daysOfWeek,
		StartTime:  startTime,
		EndTime:    endTime,
		CreatedAt:  time.Now().UTC(),
	}

	createdSchedule, err := u.scheduleRepo.Create(ctx, schedule)
	if err != nil {
		u.logger.Error().Err(err).Msg("failed to create schedule")
		return nil, err
	}

	u.logger.Info().Str("schedule_id", createdSchedule.ID.String()).Msg("schedule created successfully")

	return createdSchedule, nil
}
