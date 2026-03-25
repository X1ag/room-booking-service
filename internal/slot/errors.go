package slot

import "errors"

var (
	ErrDateRequired = errors.New("date query parameter is required")
	ErrInvalidDate = errors.New("invalid date format")
	ErrDayDoesNotApply = errors.New("schedule does not apply to this date") 
)