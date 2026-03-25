package handlers

import (
	"errors"
	"net/http"
	"test-backend-1-X1ag/internal/http/dto"
	"test-backend-1-X1ag/internal/http/response"
	"test-backend-1-X1ag/internal/room"
	"test-backend-1-X1ag/internal/slot"

	"github.com/gin-gonic/gin"
)

type SlotHandler struct {
	usecase *slot.SlotUsecase
}

func NewSlotHandler(usecase *slot.SlotUsecase) *SlotHandler {
	return &SlotHandler{
		usecase: usecase,
	}
}

func (h *SlotHandler) GetSlotsByRoomID() gin.HandlerFunc {
	return func(c *gin.Context) {
		slots, err := h.usecase.GetByRoomID(c.Request.Context(), c.Param("roomId"), c.Query("date"))
		if err != nil {
			if errors.Is(err, room.ErrInvalidRoomID) {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid room id format")
				return
			}
			if errors.Is(err, room.ErrRoomNotFound) {
				response.JSONError(c, http.StatusNotFound, response.ErrorCodeRoomNotFound, "room not found")
				return
			}
			if errors.Is(err, slot.ErrDateRequired) {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "date query parameter is required")
				return
			}
			if errors.Is(err, slot.ErrInvalidDate) {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid date format, expected YYYY-MM-DD")
				return
			}
			
			response.JSONError(c, http.StatusInternalServerError, response.ErrorCodeInternal, "failed to get slots")
			return
		}
		c.JSON(http.StatusOK, dto.GetSlotsResponse{Slots: slots})
	}
}
