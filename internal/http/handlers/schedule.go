package handlers

import (
	"errors"
	"net/http"

	"test-backend-1-X1ag/internal/http/dto"
	"test-backend-1-X1ag/internal/http/response"
	"test-backend-1-X1ag/internal/schedule"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	usecase *schedule.ScheduleUsecase
}

func NewScheduleHandler(usecase *schedule.ScheduleUsecase) *ScheduleHandler {
	return &ScheduleHandler{
		usecase: usecase,
	}
}

func (h *ScheduleHandler) Create() func(c *gin.Context) {
	return func(c *gin.Context) {
		var req dto.CreateScheduleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid request")
			return
		}

		roomID := c.Param("roomId")
		if roomID == "" {
			response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "missing room id")
			return
		}

		createdSchedule, err := h.usecase.Create(c.Request.Context(), roomID, req.DaysOfWeek, req.StartTime, req.EndTime)
		if err != nil {
			if errors.Is(err, schedule.ErrInvalidDaysOfWeek) {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid days of week")
				return
			}
			if errors.Is(err, schedule.ErrInvalidTime) {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid time format")
				return
			}
			if errors.Is(err, schedule.ErrStartTimeAfterEndTime) {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "start time must be before end time")
				return
			}
			if errors.Is(err, schedule.ErrInvalidRoomID) {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid room id format")
				return
			}
			if errors.Is(err, schedule.ErrScheduleAlreadyExists) {
				response.JSONError(c, http.StatusConflict, response.ErrorCodeScheduleExists, "schedule already exists")
				return
			}
			if errors.Is(err, schedule.ErrRoomNotFound) {
				response.JSONError(c, http.StatusNotFound, response.ErrorCodeRoomNotFound, "room not found")
				return
			}
			response.JSONError(c, http.StatusInternalServerError, response.ErrorCodeInternal, "failed to create schedule")
			return
		}

		c.JSON(http.StatusCreated, dto.CreateScheduleResponse{Schedule: *createdSchedule})
	}
}
