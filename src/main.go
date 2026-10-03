package main

import (
	"errors"
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"

	"matrix-api-go/src/cases"
)

var (
	Version = "dev"
	Commit  = "none"
)

func main() {
	// garantizar respuestas json
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// exigir json en el body
	app.Use(func(c fiber.Ctx) error {
		method := c.Method()
		if method == fiber.MethodPost || method == fiber.MethodPut || method == fiber.MethodPatch {
			contentType := c.Get("Content-Type")
			if !strings.HasPrefix(contentType, "application/json") {
				return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{
					"error": "El encabezado Content-Type debe ser 'application/json'",
				})
			}
		}
		return c.Next()
	})

	// Casos de uso
	app.Get("/health", cases.HealthCheck)
	app.Post("/api/matrix", cases.ProcessMatrix)

	log.Fatal(app.Listen(":3000"))
}
