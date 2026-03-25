package room

import (
	"context"
	"strings"
	"time"

	"test-backend-1-X1ag/internal/logger"

	"github.com/google/uuid"
)

type RoomUsecase struct {
	repo   Repository
	logger *logger.ZerologLogger
}

func NewRoomUsecase(repo Repository, logger *logger.ZerologLogger) *RoomUsecase {
	return &RoomUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (u *RoomUsecase) Create(ctx context.Context, name string, description *string, capacity *int) (*Room, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrInvalidName
	}

	if capacity != nil && *capacity <= 0 {
		return nil, ErrInvalidCapacity
	}

	if description != nil {
		desc := strings.TrimSpace(*description)
		if desc == "" {
			description = nil
		} else {
			description = &desc
		}
	}

	room := &Room{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Capacity:    capacity,
		CreatedAt:   time.Now().UTC(),
	}
	dbContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	id, err := u.repo.Create(dbContext, room)
	if err != nil {
		u.logger.Error().Err(err).Msg("failed to create room")
		return nil, err
	}
	room.ID = id

	u.logger.Info().Str("room_id", room.ID.String()).Msg("room created successfully")
	return room, nil
}

func (u *RoomUsecase) GetAll(ctx context.Context) ([]Room, error) {
	dbContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	rooms, err := u.repo.GetAll(dbContext)
	if err != nil {
		u.logger.Error().Err(err).Msg("failed to get all rooms")
		return nil, err
	}

	u.logger.Info().Int("count", len(rooms)).Msg("rooms retrieved successfully")
	return rooms, nil
}

func (u *RoomUsecase) GetByID(ctx context.Context, id string) (*Room, error) {
	roomUUID, err := uuid.Parse(id)
	if err != nil {
		u.logger.Error().Err(err).Msg("invalid room id format")
		return nil, ErrInvalidRoomID
	}
	dbContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	room, err := u.repo.GetByID(dbContext, roomUUID)
	if err != nil {
		u.logger.Error().Err(err).Msg("failed to get room by id")
		return nil, err
	}

	u.logger.Info().Str("room_id", room.ID.String()).Msg("room retrieved successfully")
	return room, nil
}
