package auth

import (
	"errors"
	"test-backend-1-X1ag/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	secret []byte
	ttl time.Duration 
}
func NewJWTManager(auth config.AuthConfig) *JWTManager {
	return &JWTManager{
		secret: []byte(auth.JWTSecret),	
		ttl:    auth.TokenTTL,
	}	
}

func (m *JWTManager) Generate(userID uuid.UUID, role string) (string, error) {
	if userID == uuid.Nil {
		return "", errors.New("empty user id") 
	}
	if role == "" {
		return "", errors.New("empty role") 
	}

	now := time.Now().UTC()

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) Parse(token string) (*Claims, error) {
	claims := Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	if claims.UserID == uuid.Nil {
		return nil, ErrInvalidToken
	}

	if claims.Role == "" {
		return nil, ErrInvalidToken
	}

	return &claims, nil
}
