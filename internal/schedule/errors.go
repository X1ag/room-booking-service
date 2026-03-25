package schedule

import "errors"

var (
	ErrInvalidDaysOfWeek     = errors.New("invalid days of week")
	ErrInvalidTime           = errors.New("invalid time format")
	ErrStartTimeAfterEndTime = errors.New("start time must be before end time")
	ErrInvalidRoomID         = errors.New("invalid room id format")
	ErrScheduleAlreadyExists = errors.New("schedule already exists")
	ErrScheduleNotFound      = errors.New("schedule not found")
	ErrRoomNotFound          = errors.New("room not found")
)
