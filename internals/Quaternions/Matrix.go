package Quaternions

import (
	"fmt"
)

type Matrix [][]float64

// CreateMatrix Create type matrix for every row that was send as parameter dimantions depand on input user
func CreateMatrix(row, col int, rows ...[]float64) *Matrix {
	if len(rows) != row {
		fmt.Println(fmt.Errorf("wrong number of rows, should be %d and is %d", row, len(rows)))
		return nil
	}

	for ind, r := range rows {
		if len(r) != col {
			fmt.Println(fmt.Errorf("wrong number of cols on %d index, should be %d and is %d", ind, row, len(rows)))
			return nil
		}
	}
	var m Matrix
	m = append(m, rows...)
	return &m
}

func (m *Matrix) Det3x3() (float64, error) {
	matrix := *m
	if len(*m) == len(matrix[0]) && len(*m) == 3 {
		return matrix[0][0]*(matrix[1][1]*matrix[2][2]-matrix[1][2]*matrix[2][1]) -
			matrix[0][1]*(matrix[1][0]*matrix[2][2]-matrix[1][2]*matrix[2][0]) +
			matrix[0][2]*(matrix[1][0]*matrix[2][1]-matrix[1][1]*matrix[2][0]), nil
	}
	return 0, fmt.Errorf("wrong shape of matrix")
}

// MultiplyMatrix multiplies two matrices
func MultiplyMatrix(m1, m2 Matrix) *Matrix {
	rows1, cols1 := len(m1), len(m1[0])
	_, cols2 := len(m2), len(m2[0])

	result := make(Matrix, rows1)
	for i := range result {
		result[i] = make([]float64, cols2)
	}

	for i := 0; i < rows1; i++ {
		for j := 0; j < cols2; j++ {
			for k := 0; k < cols1; k++ {
				result[i][j] += m1[i][k] * m2[k][j]
			}
		}
	}

	return &result
}

// PrintMatrix prints a matrix
func (m Matrix) PrintMatrix() {
	fmt.Println()
	for _, row := range m {
		fmt.Printf("|")
		for ind, val := range row {
			if ind == len(row)-1 {
				fmt.Printf("%.2f\t", val)
				break
			}
			fmt.Printf("%.2f;\t", val)
		}
		fmt.Println("|")
	}
}

// GetCoors Return X, Y, Z position
func (m Matrix) GetCoors() (r []float64) {
	for _, val := range m {
		r = append(r, val[3])
	}
	return r
}

// MultiplyMatrices multiplies multiple matrices
func MultiplyMatrices(matrices ...Matrix) *Matrix {
	result := matrices[0]
	for i := 1; i < len(matrices); i++ {
		result = *MultiplyMatrix(result, matrices[i])
	}
	return &result
}
