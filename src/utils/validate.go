package utils

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
)

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil && i > 0 {
			return i
		}
	}
	return defaultVal
}

// ValidateMatriz comprueba que la matriz sea válida, rectangular, finita y dentro de límites seguros
func ValidateMatriz(m [][]float64) error {
	rows := len(m)

	if rows == 0 {
		return errors.New("la matriz no puede estar vacía")
	}

	cols := len(m[0])

	if cols == 0 {
		return errors.New("las filas de la matriz no pueden estar vacías")
	}

	maxRows := getEnvInt("MATRIX_MAX_ROWS", 100)
	maxCols := getEnvInt("MATRIX_MAX_COLUMNS", 100)
	maxElements := getEnvInt("MATRIX_MAX_ELEMENTS", 10000)

	if rows > maxRows {
		return fmt.Errorf("la matriz excede el número máximo de filas permitido: %d > %d", rows, maxRows)
	}

	if cols > maxCols {
		return fmt.Errorf("la matriz excede el número máximo de columnas permitido: %d > %d", cols, maxCols)
	}

	if rows*cols > maxElements {
		return fmt.Errorf("la matriz excede el número total de elementos permitido: %d > %d", rows*cols, maxElements)
	}

	// este límite garantiza que no se llegue al overflow fácilmente durante el cálculo de normas
	maxAbs := math.Sqrt(math.MaxFloat64 / (8 * float64(rows)))

	for i := 0; i < rows; i++ {
		if len(m[i]) != cols {
			return errors.New(
				"la matriz no es rectangular: todas las filas deben tener el mismo número de columnas",
			)
		}

		for j := 0; j < cols; j++ {
			x := m[i][j]

			// NaN no es válido para operaciones numéricas.
			if math.IsNaN(x) {
				return fmt.Errorf(
					"la matriz contiene NaN en [%d][%d]",
					i,
					j,
				)
			}

			// +Inf y -Inf no son válidos.
			if math.IsInf(x, 0) {
				return fmt.Errorf(
					"la matriz contiene infinito en [%d][%d]",
					i,
					j,
				)
			}

			// evitar overflow en x*x durante el cálculo de la norma.
			if math.Abs(x) > maxAbs {
				return fmt.Errorf(
					"el valor absoluto en [%d][%d] es demasiado grande: %.6e (máximo permitido: %.6e)",
					i,
					j,
					x,
					maxAbs,
				)
			}
		}
	}

	return nil
}
