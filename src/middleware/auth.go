package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"matrix-api-go/src/utils"
)

// RequireAuth valida que la petición contenga un JWT válido en el encabezado Authorization
func RequireAuth(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "token de autenticación no proporcionado",
		})
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "formato de autorización inválido. Debe ser 'Bearer <token>'",
		})
	}

	tokenStr := strings.TrimSpace(parts[1])
	claims, err := utils.ValidateUserToken(tokenStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":  "token de autenticación inválido o expirado",
			"detail": err.Error(),
		})
	}

	c.Locals("user", claims.Username)
	return c.Next()
}
