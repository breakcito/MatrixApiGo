package utils

// rotar una matriz en sentido horario/derecha
// si es un rectangulo cambia su dimension de MxN a NxM
func Rotate(m [][]float64) [][]float64 {
	rows := len(m)
	if rows == 0 {
		return [][]float64{}
	}
	cols := len(m[0])

	// la primera columna de la matriz original leida de abajo
	// hacia arriba se convierte en la primera fila de la nueva matriz
	// la segunda columna original pasa a ser la segunda fila y asi sucesivamente
	res := make([][]float64, cols)
	for i := range res {
		res[i] = make([]float64, rows)
		for j := 0; j < rows; j++ {
			res[i][j] = m[rows-1-j][i]
		}
	}
	return res
}
