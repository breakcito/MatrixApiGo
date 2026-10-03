package dto

import (
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expiresIn"`
	User      string `json:"user"`
}

type UserClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
