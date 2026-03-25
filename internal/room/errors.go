package room

import "errors"

var (
	ErrInvalidName     = errors.New("invalid room name")
	ErrInvalidCapacity = errors.New("invalid room capacity")
	ErrInvalidRoomID   = errors.New("invalid room id")
	ErrRoomNotFound    = errors.New("room not found")
)
