package dto

import "test-backend-1-X1ag/internal/booking"

type CreateBookingRequest struct {
	SlotID               string `json:"slotId" binding:"required,uuid"`
	CreateConferenceLink bool   `json:"createConferenceLink"`
}

type CreateBookingResponse struct {
	Booking booking.Booking `json:"booking"`
}

type CancelBookingResponse struct {
	Booking booking.Booking `json:"booking"`
}

type GetUserBookingsResponse struct {
	Bookings []booking.Booking `json:"bookings"`
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

type ListBookingsResponse struct {
	Bookings   []booking.Booking `json:"bookings"`
	Pagination Pagination        `json:"pagination"`
}
