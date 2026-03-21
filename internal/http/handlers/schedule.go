package handlers

import (
	"test-backend-1-X1ag/internal/schedule"
)

type ScheduleHandler struct {
	usecase *schedule.ScheduleUsecase
}

func NewScheduleHandler(usecase *schedule.ScheduleUsecase) *ScheduleHandler {
	return &ScheduleHandler{
		usecase: usecase,
	}
}