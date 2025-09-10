package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/bagasunix/transnovasi/internal/dtos/responses"
)

type Claims struct {
	User *responses.UserResponse `json:"user,omitempty"`
	jwt.StandardClaims
}

// Fungsi untuk membuat Claims langsung
func NewClaims(user *responses.UserResponse, customer *responses.CustomerResponse, expiresAt time.Time) *Claims {
	return &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}
}
