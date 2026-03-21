package handlers

import (
	"test-backend-1-X1ag/internal/slot"
)

type SlotHandler struct {
	usecase *slot.SlotUsecase
}

func NewSlotHandler(usecase *slot.SlotUsecase) *SlotHandler {
	return &SlotHandler{
		usecase: usecase,
	}
}