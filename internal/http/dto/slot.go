package dto

import "test-backend-1-X1ag/internal/slot"

type GetSlotsResponse struct {
	Slots []slot.Slot `json:"slots"`
}
