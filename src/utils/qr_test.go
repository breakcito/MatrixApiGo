package utils

import (
	"math"
	"testing"
)

func TestValidateMatriz_ValidSmallNumbers(t *testing.T) {
	// 1e-15 debe ser aceptado y no rechazado arbitrariamente
	m := [][]float64{
		{1e-15, 2.0},
		{0.0, 3.0},
	}
	if err := ValidateMatriz(m); err != nil {
		t.Fatalf("se esperaba que 1e-15 fuera aceptado, pero dio error: %v", err)
	}
}

func TestValidateMatriz_Empty(t *testing.T) {
	if err := ValidateMatriz([][]float64{}); err == nil {
		t.Fatal("se esperaba error para matriz vacía")
	}

	if err := ValidateMatriz([][]float64{{}}); err == nil {
		t.Fatal("se esperaba error para filas vacías [[]]")
	}
}

func TestValidateMatriz_NonRectangular(t *testing.T) {
	m := [][]float64{
		{1.0, 2.0},
		{3.0},
	}
	if err := ValidateMatriz(m); err == nil {
		t.Fatal("se esperaba error para matriz no rectangular")
	}
}

func TestValidateMatriz_NaNAndInf(t *testing.T) {
	m1 := [][]float64{{math.NaN()}}
	if err := ValidateMatriz(m1); err == nil {
		t.Fatal("se esperaba error para NaN")
	}

	m2 := [][]float64{{math.Inf(1)}}
	if err := ValidateMatriz(m2); err == nil {
		t.Fatal("se esperaba error para +Inf")
	}
}

func TestValidateMatriz_DimensionGuards(t *testing.T) {
	// Crear matriz que excede 100 filas
	big := make([][]float64, 101)
	for i := range big {
		big[i] = []float64{1.0}
	}
	if err := ValidateMatriz(big); err == nil {
		t.Fatal("se esperaba error al exceder máximo de filas")
	}
}

func TestCalculateQR_SubdiagonalClean(t *testing.T) {
	// Matriz 3x3 de prueba
	a := [][]float64{
		{12, -51, 4},
		{6, 167, -68},
		{-4, 24, -41},
	}

	q, r := CalculateQR(a)

	// Verificar que bajo la diagonal de R todos los elementos sean estrictamente 0.0
	for i := 0; i < len(r); i++ {
		for j := 0; j < len(r[0]); j++ {
			if i > j && r[i][j] != 0.0 {
				t.Fatalf("R[%d][%d] = %e, debería ser exactamente 0.0", i, j, r[i][j])
			}
		}
	}

	// Verificar ortogonalidad de Q: Q^T * Q = I
	m := len(q)
	for i := 0; i < m; i++ {
		for j := 0; j < m; j++ {
			dot := 0.0
			for k := 0; k < m; k++ {
				dot += q[k][i] * q[k][j]
			}
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			if math.Abs(dot-expected) > 1e-9 {
				t.Fatalf("Q no es ortogonal en [%d][%d]: %f != %f", i, j, dot, expected)
			}
		}
	}

	// Verificar reconstrucción: Q * R ≈ A
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			dot := 0.0
			for k := 0; k < len(r); k++ {
				dot += q[i][k] * r[k][j]
			}
			if math.Abs(dot-a[i][j]) > 1e-9 {
				t.Fatalf("Q*R != A en [%d][%d]: %f != %f", i, j, dot, a[i][j])
			}
		}
	}
}

func TestCalculateQR_ZeroSubdiagonalResidue(t *testing.T) {
	// Caso con valores donde el reflector podría dejar residuos diminutos
	a := [][]float64{
		{1, 2},
		{0, 3},
		{0, 0},
	}

	_, r := CalculateQR(a)
	for i := 0; i < len(r); i++ {
		for j := 0; j < len(r[0]); j++ {
			if i > j && r[i][j] != 0.0 {
				t.Fatalf("R[%d][%d] = %e, residuo detectado bajo la diagonal", i, j, r[i][j])
			}
		}
	}
}
