package handlers

import (
	"net/http"
	"test-backend-1-X1ag/internal/auth"
	"test-backend-1-X1ag/internal/http/dto"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	usecase *auth.AuthUsecase	
}

func NewAuthHandler(usecase *auth.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		usecase: usecase,
	}
}

func (h *AuthHandler) DummyLogin(c *gin.Context) {
	var req dto.DummyLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}

	token, err := h.usecase.DummyLogin(c.Request.Context(), req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, dto.TokenResponse{
		Token: token,
	})
}