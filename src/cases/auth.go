package cases

import (
	"os"

	"github.com/gofiber/fiber/v3"

	"matrix-api-go/src/dto"
	"matrix-api-go/src/utils"
)

// Login valida credenciales y entrega un token JWT
func Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cuerpo de la petición inválido",
		})
	}

	expectedUser := os.Getenv("AUTH_USER")
	if expectedUser == "" {
		expectedUser = "admin"
	}

	expectedPass := os.Getenv("AUTH_PASS")
	if expectedPass == "" {
		expectedPass = "matrix2026"
	}

	if req.Username != expectedUser || req.Password != expectedPass {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "credenciales incorrectas",
		})
	}

	token, expiresIn, err := utils.GenerateUserToken(req.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "no se pudo generar el token de sesión",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.LoginResponse{
		Token:     token,
		ExpiresIn: expiresIn,
		User:      req.Username,
	})
}
