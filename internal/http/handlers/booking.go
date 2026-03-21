package handlers

import "test-backend-1-X1ag/internal/booking"

type BookingHandler struct {
	usecase *booking.BookingUsecase
}

func NewBookingHandler(usecase *booking.BookingUsecase) *BookingHandler {
	return &BookingHandler{
		usecase: usecase,
	}
}