package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"test-backend-1-X1ag/internal/auth"
	"test-backend-1-X1ag/internal/booking"
	"test-backend-1-X1ag/internal/http/dto"
	"test-backend-1-X1ag/internal/http/response"
	"test-backend-1-X1ag/internal/slot"

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

func currentUserID(c *gin.Context) (string, bool) {
	info, ok := auth.AuthInfoFromContext(c.Request.Context())
	if !ok {
		response.JSONError(c, http.StatusUnauthorized, response.ErrorCodeUnauthorized, "unauthorized")
		return "", false
	}

	return info.UserID.String(), true
}

func (h *BookingHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateBookingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid request body")
			return
		}

		userID, ok := currentUserID(c)
		if !ok {
			return
		}

		createdBooking, err := h.usecase.Create(c.Request.Context(), req.SlotID, req.CreateConferenceLink, userID)
		if err != nil {
			if errors.Is(err, booking.ErrInvalidSlotID) {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid slot id format")
				return
			}
			if errors.Is(err, booking.ErrInvalidUserID) {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid user id format")
				return
			}
			if errors.Is(err, slot.ErrSlotNotFound) {
				response.JSONError(c, http.StatusNotFound, response.ErrorCodeSlotNotFound, "slot not found")
				return
			}
			if errors.Is(err, booking.ErrSlotInPast) {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "cannot book a slot in the past")
				return
			}
			if errors.Is(err, booking.ErrSlotAlreadyBooked) {
				response.JSONError(c, http.StatusConflict, response.ErrorCodeSlotAlreadyBooked, "slot is already booked")
				return
			}
			response.JSONError(c, http.StatusInternalServerError, response.ErrorCodeInternal, "failed to create booking")
			return
		}

		c.JSON(http.StatusCreated, dto.CreateBookingResponse{Booking: createdBooking})
	}
}

func (h *BookingHandler) Cancel() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := currentUserID(c)
		if !ok {
			return
		}

		cancelledBooking, err := h.usecase.Cancel(c.Request.Context(), c.Param("bookingId"), userID)
		if err != nil {
			if errors.Is(err, booking.ErrInvalidBookingID) {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid booking id format")
				return
			}
			if errors.Is(err, booking.ErrInvalidUserID) {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid user id format")
				return
			}
			if errors.Is(err, booking.ErrBookingNotFound) {
				response.JSONError(c, http.StatusNotFound, response.ErrorCodeBookingNotFound, "booking not found")
				return
			}
			if errors.Is(err, booking.ErrForbidden) {
				response.JSONError(c, http.StatusForbidden, response.ErrorCodeForbidden, "cannot cancel another user's booking")
				return
			}
			response.JSONError(c, http.StatusInternalServerError, response.ErrorCodeInternal, "failed to cancel booking")
			return
		}

		c.JSON(http.StatusOK, dto.CancelBookingResponse{Booking: cancelledBooking})
	}
}

func (h *BookingHandler) ListBookings() gin.HandlerFunc {
	return func(c *gin.Context) {
		page := 1
		pageSize := 20

		if rawPage := c.Query("page"); rawPage != "" {
			parsedPage, err := strconv.Atoi(rawPage)
			if err != nil {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid page")
				return
			}
			page = parsedPage
		}

		if rawPageSize := c.Query("pageSize"); rawPageSize != "" {
			parsedPageSize, err := strconv.Atoi(rawPageSize)
			if err != nil {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid pageSize")
				return
			}
			pageSize = parsedPageSize
		}

		bookings, total, err := h.usecase.List(c.Request.Context(), page, pageSize)
		if err != nil {
			if errors.Is(err, booking.ErrInvalidPagination) {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid pagination params")
				return
			}
			response.JSONError(c, http.StatusInternalServerError, response.ErrorCodeInternal, "failed to list bookings")
			return
		}

		c.JSON(http.StatusOK, dto.ListBookingsResponse{
			Bookings: bookings,
			Pagination: dto.Pagination{
				Page:     page,
				PageSize: pageSize,
				Total:    total,
			},
		})
	}
}

func (h *BookingHandler) GetUserBookings() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := currentUserID(c)
		if !ok {
			return
		}

		bookings, err := h.usecase.GetUserBookings(c.Request.Context(), userID)
		if err != nil {
			if errors.Is(err, booking.ErrInvalidUserID) {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid user id format")
				return
			}
			response.JSONError(c, http.StatusInternalServerError, response.ErrorCodeInternal, "failed to get user bookings")
			return
		}

		c.JSON(http.StatusOK, dto.GetUserBookingsResponse{Bookings: bookings})
	}
}
