package dto

import "test-backend-1-X1ag/internal/user"

type DummyLoginRequest struct {
	Role string `json:"role" binding:"required,oneof=admin user"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,max=72"`
}

type RegisterResponse struct {
	User user.User `json:"user"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
