package handlers

import (
	"errors"
	"net/http"

	"test-backend-1-X1ag/internal/auth"
	"test-backend-1-X1ag/internal/http/dto"
	"test-backend-1-X1ag/internal/http/response"
	"test-backend-1-X1ag/internal/user"

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

func (h *AuthHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid request")
			return
		}

		createdUser, err := h.usecase.Register(c.Request.Context(), req.Email, req.Password, req.Role)
		if err != nil {
			if errors.Is(err, user.ErrInvalidEmail) || errors.Is(err, user.ErrInvalidPassword) || errors.Is(err, user.ErrInvalidRole) || errors.Is(err, user.ErrEmailAlreadyExists) {
				response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, err.Error())
				return
			}

			response.JSONError(c, http.StatusInternalServerError, response.ErrorCodeInternal, "internal server error")
			return
		}

		c.JSON(http.StatusCreated, dto.RegisterResponse{User: *createdUser})
	}
}

func (h *AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.JSONError(c, http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid request")
			return
		}

		token, err := h.usecase.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			if errors.Is(err, user.ErrInvalidCredentials) {
				response.JSONError(c, http.StatusUnauthorized, response.ErrorCodeUnauthorized, "invalid email or password")
				return
			}

			response.JSONError(c, http.StatusInternalServerError, response.ErrorCodeInternal, "internal server error")
			return
		}

		c.JSON(http.StatusOK, dto.TokenResponse{Token: token})
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

	c.JSON(http.StatusOK, dto.TokenResponse{Token: token})
}
