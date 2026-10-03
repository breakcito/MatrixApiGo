package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"matrix-api-go/src/dto"

	"github.com/golang-jwt/jwt/v5"
)

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// GenerateUserToken genera un token JWT firmado para el frontend
func GenerateUserToken(username string) (string, int64, error) {
	secret := []byte(getEnv("JWT_SECRET", "super-secret-go-front-key-2026"))
	expiresIn := int64(24 * 3600) // 24 horas

	claims := dto.UserClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username,
			Issuer:    "matrix-api-go",
			Audience:  jwt.ClaimStrings{"matrix-frontend"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", 0, fmt.Errorf("error al firmar token de usuario: %w", err)
	}

	return tokenStr, expiresIn, nil
}

// ValidateUserToken valida el token JWT proveniente del frontend
func ValidateUserToken(tokenStr string) (*dto.UserClaims, error) {
	secret := []byte(getEnv("JWT_SECRET", "super-secret-go-front-key-2026"))

	token, err := jwt.ParseWithClaims(tokenStr, &dto.UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", t.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token inválido: %w", err)
	}

	claims, ok := token.Claims.(*dto.UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("claims de token inválidos")
	}

	return claims, nil
}
// GenerateServiceToken genera un token JWT firmado de servicio para comunicarse con Node.js
func GenerateServiceToken() (string, error) {
	secretStr := getEnv("STATS_JWT_SECRET", "")
	secret := []byte(secretStr)

	claims := jwt.RegisteredClaims{
		Subject:   "matrix-api-go",
		Issuer:    "matrix-api-go",
		Audience:  jwt.ClaimStrings{"matrix-api-node"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("error al firmar token de servicio: %w", err)
	}

	return tokenStr, nil
}
