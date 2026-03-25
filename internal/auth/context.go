package auth

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const authInfoKey contextKey = "auth_info"

type AuthInfo struct {
	UserID uuid.UUID
	Role   string
}

func WithAuthInfo(ctx context.Context, info AuthInfo) context.Context {
	return context.WithValue(ctx, authInfoKey, info)
}

func AuthInfoFromContext(ctx context.Context) (AuthInfo, bool) {
	info, ok := ctx.Value(authInfoKey).(AuthInfo)
	return info, ok
}
