package utils

import (
	"errors"
	"fmt"
	"math"
)

// comprobar que la matriz sea valida
func ValidateMatriz(m [][]float64) error {
	rows := len(m)

	if rows == 0 {
		return errors.New("la matriz no puede estar vacía")
	}

	cols := len(m[0])

	if cols == 0 {
		return errors.New("las filas de la matriz no pueden estar vacías")
	}

	// este limite garantiza que no se llegue al overflow facilmente
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

			// Los valores extremadamente pequeños pueden perder precisión
			if x != 0 && math.Abs(x) < math.Sqrt(math.SmallestNonzeroFloat64) {
				return fmt.Errorf(
					"el valor absoluto en [%d][%d] es demasiado pequeño: %.6e",
					i,
					j,
					x,
				)
			}

			// rechazamos valores no nulos por debajo de la tolerancia
			if x != 0 && math.Abs(x) < 1e-12 {
				return fmt.Errorf(
					"el valor absoluto en [%d][%d] es menor que la tolerancia QR: %.6e < %.6e",
					i,
					j,
					math.Abs(x),
					1e-12,
				)
			}
		}
	}

	return nil
}
