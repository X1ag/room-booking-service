package conference

import (
	"context"

	"github.com/google/uuid"
)

type MockService struct{}

func NewMockService() *MockService {
	return &MockService{}
}

func (s *MockService) CreateLink(ctx context.Context, slotID uuid.UUID, userID uuid.UUID) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	return "https://meet.example.com/" + uuid.NewString(), nil
}
