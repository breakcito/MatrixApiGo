package utils

import (
	"math"
)

// calcular Q y R usando reflexiones de Householder
func CalculateQR(a [][]float64) ([][]float64, [][]float64) {
	m := len(a)
	n := len(a[0])

	// copiar A -> R para no modificar la matriz original
	R := make([][]float64, m)
	for i := 0; i < m; i++ {
		R[i] = make([]float64, n)
		copy(R[i], a[i])
	}

	// Q comienza como la identidad
	Q := make([][]float64, m)
	for i := 0; i < m; i++ {
		Q[i] = make([]float64, m)
		Q[i][i] = 1.0
	}

	// Cantidad de reflectores necesarios
	kMax := m
	if n < m {
		kMax = n
	}

	for k := 0; k < kMax; k++ {
		// Construimos el vector x = R[k:m][k]
		x := make([]float64, m-k)

		normX := 0.0
		for i := k; i < m; i++ {
			x[i-k] = R[i][k]
			normX += R[i][k] * R[i][k]
		}

		normX = math.Sqrt(normX)

		// Si la columna ya es prácticamente cero,
		// no necesitamos aplicar reflector.
		if normX < 1e-12 {
			continue
		}

		// Householder:
		// v = x + sign(x1)*||x||*e1
		//
		// Elegimos el signo para evitar cancelación numérica.
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

		if vNorm2 < 1e-24 {
			continue
		}

		// Aplicar H = I - 2 vv^T / (v^T v)
		// a R desde la izquierda: R <- H R
		for j := k; j < n; j++ {
			// dot = v^T * R[k:m][j]
			dot := 0.0
			for i := k; i < m; i++ {
				dot += v[i-k] * R[i][j]
			}

			factor := 2.0 * dot / vNorm2

			for i := k; i < m; i++ {
				R[i][j] -= factor * v[i-k]
			}
		}

		// Acumulamos Q.
		//
		// Como vamos aplicando:
		// R = H_k ... H_2 H_1 A
		//
		// entonces:
		// Q = H_1 H_2 ... H_k
		//
		// Para acumularla correctamente:
		// Q <- Q H
		for i := 0; i < m; i++ {
			// Calculamos fila i de Q * H
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

	return Q, R
}
