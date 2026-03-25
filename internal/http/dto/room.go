package dto

import (
	"test-backend-1-X1ag/internal/room"
)

type CreateRoomRequest struct {
	Name        string `json:"name" binding:"required"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity"`
}

type CreateRoomResponse struct {
	Room room.Room `json:"room"` 
}

type GetRoomsResponse struct {
	Rooms []room.Room `json:"rooms"`
}