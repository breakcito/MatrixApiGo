package dto

// representa la matriz que envia el front
type MatrixRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

// representa la matriz ya rotada, la matriz q
// y la matriz r que enviamos a la api de node
// para sacar indicadores
type StatsRequest struct {
	Rotated [][]float64 `json:"rotated"`
	Q       [][]float64 `json:"q"`
	R       [][]float64 `json:"r"`
}

// representa lo que la api de node responde
// con los valores de las 5 operaciones requeridas
type StatsResponse struct {
	MaxValue       float64 `json:"maxValue"`
	MinValue       float64 `json:"minValue"`
	Average        float64 `json:"average"`
	TotalSum       float64 `json:"totalSum"`
	HasDiagonalMat bool    `json:"hasDiagonalMat"`
}

// representa la respuesta que le enviamos al front
// para que lo renderice
type FinalResponse struct {
	Original [][]float64       `json:"original"`
	Rotated  [][]float64       `json:"rotated"`
	Q        [][]float64       `json:"q"`
	R        [][]float64       `json:"r"`
	Stats    StatsResponse `json:"stats"`
}
