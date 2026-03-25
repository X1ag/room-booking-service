package handlers

import (
	"test-backend-1-X1ag/internal/booking"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	usecase *booking.BookingUsecase
}

func NewBookingHandler(usecase *booking.BookingUsecase) *BookingHandler {
	return &BookingHandler{
		usecase: usecase,
	}
}

func (h *BookingHandler) Create() gin.HandlerFunc {
	return func (c *gin.Context) {
		// TODO: implement booking creation handler
		// need to check: 
		// - if slot is available for booking, 
		// - if user has no booking for this slot, 
		// - slot is not in the past,
		// - check if slot have no bookings for this room and time  
		// - create booking
	}
}

func (h *BookingHandler) Cancel() gin.HandlerFunc {
	return func (c *gin.Context) {
		// TODO: implement booking cancellation
		// need to check:
		// - if canceled booking owner is the same as user who tries to cancel, 
		// - if canceled booking is not already canceled
	}
}

func (h *BookingHandler) ListBookings() gin.HandlerFunc {
	return func (c *gin.Context) {
		// TODO: implement bookings list handler(by page or limit-offset)
	}
}

func (h *BookingHandler) GetUserBookings() gin.HandlerFunc {
	return func (c *gin.Context) {
		// TODO: implement user bookings list handler
	}
}