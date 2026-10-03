package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"matrix-api-go/src/cases"
)

var (
	Version = "dev"
	Commit  = "none"
)

func main() {
	// Configuración de Fiber con límite de body simétrico a Node.js (10MB)
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10MB
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

	// Configuración de CORS según APP_CLIENT
	appClient := os.Getenv("APP_CLIENT")
	corsCfg := cors.Config{
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}
	if appClient != "" {
		var origins []string
		for _, o := range strings.Split(appClient, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				origins = append(origins, o)
			}
		}
		corsCfg.AllowOrigins = origins
	} else {
		corsCfg.AllowOrigins = []string{"*"}
	}
	app.Use(cors.New(corsCfg))

	// Exigir application/json en el body para peticiones de mutación
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	log.Printf("Go server listening on port %s\n", port)
	log.Fatal(app.Listen(port))
}
