package middleware

import (
	"net/http"
	"strings"
	"test-backend-1-X1ag/internal/auth"
	"test-backend-1-X1ag/internal/http/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtManager auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")	
		if header == "" {
			response.JSONError(c, http.StatusUnauthorized, response.ErrorCodeUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")

		if token == "" {
			response.JSONError(c, http.StatusUnauthorized, response.ErrorCodeUnauthorized, "missing token")
			c.Abort()
			return
		}

		claims, err := jwtManager.Parse(token)
		if err != nil {
			response.JSONError(c, http.StatusUnauthorized, response.ErrorCodeUnauthorized, err.Error())
			c.Abort()	
			return
		}

		info := auth.AuthInfo{
			UserID: claims.UserID,
			Role: claims.Role,
		}

		ctx := auth.WithAuthInfo(c.Request.Context(), info)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}