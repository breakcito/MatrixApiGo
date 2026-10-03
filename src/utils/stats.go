package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"

	"matrix-api-go/src/dto"
)

// FetchStats envía las matrices a la API de Node y retorna sus métricas
func FetchStats(data dto.StatsRequest) (*dto.StatsResponse, error) {
	nodeURL := os.Getenv("API_STATS_URL")

	if nodeURL == "" {
		return nil, errors.New("API_STATS_URL no está configurada")
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error serializando payload: %w", err)
	}

	resp, err := (&http.Client{Timeout: 10 * time.Second}).Post(nodeURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error de conexión con Node.js: %w", err)
	}
	defer resp.Body.Close()

	// Si Node devuelve 413, preservamos el 413 en vez de traducir a 502
	if resp.StatusCode == http.StatusRequestEntityTooLarge {
		return nil, fiber.NewError(fiber.StatusRequestEntityTooLarge, "el payload de las matrices excede el límite permitido por el servicio de estadísticas (413)")
	}

	// Si Node responde con cualquier otro error HTTP
	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Error != "" {
			return nil, fiber.NewError(resp.StatusCode, errResp.Error)
		}
		return nil, fiber.NewError(resp.StatusCode, fmt.Sprintf("Node.js respondió con status %d", resp.StatusCode))
	}

	var stats dto.StatsResponse
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta de Node.js: %w", err)
	}

	return &stats, nil
}
