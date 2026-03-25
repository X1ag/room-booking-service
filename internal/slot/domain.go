package slot

import (
	"time"

	"github.com/google/uuid"
)

type Slot struct {
	ID        uuid.UUID `json:"id"`
	RoomID    uuid.UUID `json:"roomId"`
	StartTime time.Time `json:"start"`
	EndTime   time.Time `json:"end"`
}
