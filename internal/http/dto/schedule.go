package dto

import (
	"test-backend-1-X1ag/internal/schedule"

)

type CreateScheduleRequest struct {
	DaysOfWeek []int  `json:"daysOfWeek" binding:"required"`
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

type CreateScheduleResponse struct {
	Schedule schedule.Schedule `json:"schedule"`
}