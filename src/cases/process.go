package cases

import (
	"github.com/gofiber/fiber/v3"

	"matrix-api-go/src/dto"
	"matrix-api-go/src/utils"
)

// manejar todo, la validacion, rotacion, calculo de matrices q y r y
// enviarle esto a la api con node
func ProcessMatrix(c fiber.Ctx) error {
	var req dto.MatrixRequest

	// validar body
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cuerpo de la petición inválido",
		})
	}

	// validar que la matriz sea correcta y este dentro de los limites
	if err := utils.ValidateMatriz(req.Matrix); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// rotar la matriz
	rotated := utils.Rotate(req.Matrix)

	// obteenr las matrices q y r
	q, r := utils.CalculateQR(req.Matrix)

	// obtener estadisticas
	stats, err := utils.FetchStats(dto.StatsRequest{
		Rotated: rotated,
		Q:       q,
		R:       r,
	})

	if err != nil {
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
