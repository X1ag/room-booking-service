package handlers

import (
	"errors"
	"net/http"
	"test-backend-1-X1ag/internal/auth"
	"test-backend-1-X1ag/internal/http/dto"
	"test-backend-1-X1ag/internal/http/response"

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
		response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid request")
		return
	}

	token, err := h.usecase.DummyLogin(c.Request.Context(), req.Role)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidRole) {
			response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, err.Error())
			return
		}

		response.JSONError(c, http.StatusInternalServerError, response.ErrorCodeInternal, "internal server error")
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{
		Token: token,
	})
}
