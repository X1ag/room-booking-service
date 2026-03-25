package middleware

import (
	"net/http"

	"test-backend-1-X1ag/internal/auth"
	"test-backend-1-X1ag/internal/http/response"

	"github.com/gin-gonic/gin"
)

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		info, ok := auth.AuthInfoFromContext(c.Request.Context())
		if !ok {
			response.JSONError(c, http.StatusUnauthorized, response.ErrorCodeUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		if info.Role != role {
			response.JSONError(c, http.StatusForbidden, response.ErrorCodeForbidden, "forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}
