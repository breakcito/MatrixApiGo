package cases

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"matrix-api-go/src/dto"
	"matrix-api-go/src/utils"
)

// ProcessMatrix maneja validación, rotación, factorización QR y llamada al servicio de estadísticas
func ProcessMatrix(c fiber.Ctx) error {
	var req dto.MatrixRequest

	// validar body
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cuerpo de la petición inválido",
		})
	}

	// validar que la matriz sea correcta y esté dentro de los límites
	if err := utils.ValidateMatriz(req.Matrix); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// rotar la matriz
	rotated := utils.Rotate(req.Matrix)

	// obtener las matrices Q y R
	q, r := utils.CalculateQR(req.Matrix)

	// obtener estadísticas de la API de Node
	stats, err := utils.FetchStats(dto.StatsRequest{
		Rotated: rotated,
		Q:       q,
		R:       r,
	})

	if err != nil {
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"error": fiberErr.Message,
			})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":  "error al consultar estadísticas",
			"detail": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.FinalResponse{
		Original: req.Matrix,
		Rotated:  rotated,
		Q:        q,
		R:        r,
		Stats:    *stats,
	})
}
