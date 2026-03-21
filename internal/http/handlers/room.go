package handlers

import "test-backend-1-X1ag/internal/room"

type RoomHandler struct {
	usecase *room.RoomUsecase
}

func NewRoomHandler(usecase *room.RoomUsecase) *RoomHandler {
	return &RoomHandler{
		usecase: usecase,
	}
}