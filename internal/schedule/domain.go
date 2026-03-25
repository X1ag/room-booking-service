package schedule

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID   uuid.UUID `json:"id"`
	RoomID uuid.UUID `json:"roomId"`
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"`
	DaysOfWeek []int `json:"daysOfWeek"`
	CreatedAt time.Time `json:"createdAt"`
}