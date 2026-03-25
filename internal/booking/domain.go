package booking

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"userId"`
	SlotID         uuid.UUID `json:"slotId"`
	ConferenceLink *string   `json:"conferenceLink,omitempty"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
}
