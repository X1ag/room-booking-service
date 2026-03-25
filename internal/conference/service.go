package conference

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	CreateLink(ctx context.Context, slotID uuid.UUID, userID uuid.UUID) (string, error)
}
