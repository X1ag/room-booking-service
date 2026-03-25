package response

import "github.com/gin-gonic/gin"

type ErrorCode string

const (
	ErrorCodeInvalidRequest  ErrorCode = "INVALID_REQUEST"
	ErrorCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrorCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrorCodeRoomNotFound    ErrorCode = "ROOM_NOT_FOUND"
	ErrorCodeSlotNotFound    ErrorCode = "SLOT_NOT_FOUND"
	ErrorCodeSlotBooked      ErrorCode = "SLOT_ALREADY_BOOKED"
	ErrorCodeBookingNotFound ErrorCode = "BOOKING_NOT_FOUND"
	ErrorCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrorCodeScheduleExists  ErrorCode = "SCHEDULE_EXISTS"
	ErrorCodeInternal        ErrorCode = "INTERNAL_ERROR"
	ErrorCodeSlotAlreadyBooked ErrorCode = "SLOT_ALREADY_BOOKED"
)

type ErrorDetails struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type ErrorEnvelope struct {
	Error ErrorDetails `json:"error"`
}

func JSONError(c *gin.Context, status int, code ErrorCode, message string) {
	c.JSON(status, ErrorEnvelope{
		Error: ErrorDetails{
			Code:    code,
			Message: message,
		},
	})
}
