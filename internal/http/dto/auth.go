package dto

type DummyLoginRequest struct {
    Role string `json:"role" binding:"required"`
}

type TokenResponse struct {
    Token string `json:"token"`
}
