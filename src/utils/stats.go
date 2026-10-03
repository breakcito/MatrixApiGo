package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"matrix-api-go/src/dto"
)

// enviar las matrices a la API de Node y retorna sus metricas
func FetchStats(data dto.StatsRequest) (*dto.StatsResponse, error) {
	nodeURL := os.Getenv("NODE_API_URL")

	if nodeURL == "" {
		return nil, errors.New("NODE_API_URL no está configurada")
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error serializando payload: %w", err)
	}

	resp, err := (&http.Client{Timeout: 5 * time.Second}).Post(nodeURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error de conexión con Node.js: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Node.js respondió con status %d", resp.StatusCode)
	}

	var stats dto.StatsResponse
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta de Node.js: %w", err)
	}

	return &stats, nil
}
