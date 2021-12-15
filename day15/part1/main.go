package main

import (
	"fmt"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

func main() {
	filepath := "input.txt"
	matrix := helpers.ReadFileAsMatrixOfInts(filepath)
	helpers.PrintIntMatrix(matrix)

	solutionMatrix := helpers.NewIntMatrixOfSize(len(matrix), len(matrix), 0)
	fillSolutionMatrix(solutionMatrix, matrix)
	fmt.Println("Solution")
	helpers.PrintIntMatrix(solutionMatrix)

	firstCell := solutionMatrix[0][0]
	lastCell := solutionMatrix[len(matrix)-1][len(matrix)-1]
	fmt.Printf("First %d, Last %d\n", firstCell, lastCell)
	fmt.Printf("Answer %d\n", lastCell-firstCell)
}

func fillSolutionMatrix(solutionMatrix [][]int, matrix [][]int) {
	solutionMatrix[0][0] = matrix[0][0]

	for j := 1; j < len(matrix); j++ {
		solutionMatrix[0][j] = matrix[0][j] + solutionMatrix[0][j-1]
	}

	for i := 1; i < len(matrix); i++ {
		solutionMatrix[i][0] = matrix[i][0] + solutionMatrix[i-1][0]
	}

	for i := 1; i < len(matrix); i++ {
		for j := 1; j < len(matrix); j++ {
			solutionMatrix[i][j] = matrix[i][j] + minimum(solutionMatrix[i-1][j], solutionMatrix[i][j-1])
		}
	}
}

func minimum(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}
