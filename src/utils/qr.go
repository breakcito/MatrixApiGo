package utils

import (
	"math"
)

// CalculateQR calcula las matrices Q y R usando reflexiones de Householder
func CalculateQR(a [][]float64) ([][]float64, [][]float64) {
	m := len(a)
	n := len(a[0])

	// copiar A -> R para no modificar la matriz original
	R := make([][]float64, m)
	for i := 0; i < m; i++ {
		R[i] = make([]float64, n)
		copy(R[i], a[i])
	}

	// Q comienza como la identidad m x m
	Q := make([][]float64, m)
	for i := 0; i < m; i++ {
		Q[i] = make([]float64, m)
		Q[i][i] = 1.0
	}

	// Cantidad de reflectores necesarios: min(m-1, n)
	kMax := m - 1
	if n < kMax {
		kMax = n
	}

	for k := 0; k < kMax; k++ {
		// Comprobar la norma del subvector estrictamente debajo de la diagonal: R[k+1:m][k]
		subNorm2 := 0.0
		for i := k + 1; i < m; i++ {
			subNorm2 += R[i][k] * R[i][k]
		}

		// Si ya no hay valores por debajo de la diagonal en esta columna,
		// no se necesita aplicar reflector
		if subNorm2 == 0.0 {
			continue
		}

		// Construimos el vector x = R[k:m][k]
		x := make([]float64, m-k)
		normX := 0.0
		for i := k; i < m; i++ {
			x[i-k] = R[i][k]
			normX += R[i][k] * R[i][k]
		}

		normX = math.Sqrt(normX)
		if normX == 0.0 {
			continue
		}

		// Householder: v = x + sign(x1)*||x||*e1
		sign := 1.0
		if x[0] < 0 {
			sign = -1.0
		}

		v := make([]float64, len(x))
		copy(v, x)
		v[0] += sign * normX

		// ||v||^2
		vNorm2 := 0.0
		for _, value := range v {
			vNorm2 += value * value
		}

		if vNorm2 == 0.0 {
			continue
		}

		// Aplicar H = I - 2 vv^T / (v^T v) a R desde la izquierda: R <- H R
		for j := k; j < n; j++ {
			dot := 0.0
			for i := k; i < m; i++ {
				dot += v[i-k] * R[i][j]
			}

			factor := 2.0 * dot / vNorm2

			for i := k; i < m; i++ {
				R[i][j] -= factor * v[i-k]
			}
		}

		// Acumular Q: Q <- Q H
		for i := 0; i < m; i++ {
			dot := 0.0
			for j := k; j < m; j++ {
				dot += Q[i][j] * v[j-k]
			}

			factor := 2.0 * dot / vNorm2

			for j := k; j < m; j++ {
				Q[i][j] -= factor * v[j-k]
			}
		}
	}

	// Limpiar explícitamente todas las entradas bajo la diagonal de R.
	// Matemáticamente R es triangular superior (i > j => R[i][j] = 0).
	// Esto previene que residuos de coma flotante (ej. -4.3e-15) o columnas
	// con subdiagonales no anuladas contaminen estadísticas (minValue, totalSum, etc.).
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if i > j {
				R[i][j] = 0.0
			}
		}
	}

	return Q, R
}
