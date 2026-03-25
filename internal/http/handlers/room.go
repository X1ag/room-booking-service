package handlers

import (
	"errors"
	"net/http"
	"test-backend-1-X1ag/internal/http/dto"
	"test-backend-1-X1ag/internal/http/response"
	"test-backend-1-X1ag/internal/room"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	usecase *room.RoomUsecase
}

func NewRoomHandler(usecase *room.RoomUsecase) *RoomHandler {
	return &RoomHandler{
		usecase: usecase,
	}
}

func (h *RoomHandler) Create() func(c *gin.Context) {
	return func(c *gin.Context) {
		var req dto.CreateRoomRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid request")
			return
		}

		createdRoom, err := h.usecase.Create(c.Request.Context(), req.Name, req.Description, req.Capacity)
		if err != nil {
			if errors.Is(err, room.ErrInvalidName) {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid room name")
				return
			}
			if errors.Is(err, room.ErrInvalidCapacity) {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid room capacity")
				return
			}
			response.JSONError(c, http.StatusInternalServerError, response.ErrorCodeInternal, "failed to create room")
			return
		}

		c.JSON(http.StatusCreated, 	dto.CreateRoomResponse{Room: *createdRoom})
	}
}

func (h *RoomHandler) GetRooms() func(c *gin.Context) {
	return func(c *gin.Context) {
		rooms, err := h.usecase.GetAll(c.Request.Context())
		if err != nil {
			response.JSONError(c, http.StatusInternalServerError, response.ErrorCodeInternal, "failed to get rooms")
			return
		}
		c.JSON(http.StatusOK, dto.GetRoomsResponse{Rooms: rooms})
	}
}