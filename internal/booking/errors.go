package booking

import "errors"

var (
	ErrInvalidUserID         = errors.New("invalid user id format")
	ErrSlotAlreadyBooked     = errors.New("slot is already booked")
	ErrInvalidSlotID         = errors.New("invalid slot id format")
	ErrSlotInPast            = errors.New("slot is in the past")
	ErrConferenceUnavailable = errors.New("conference service unavailable")
	ErrInvalidBookingID      = errors.New("invalid booking id format")
	ErrBookingNotFound       = errors.New("booking not found")
	ErrForbidden             = errors.New("forbidden")
	ErrInvalidPagination     = errors.New("invalid pagination params")
)
