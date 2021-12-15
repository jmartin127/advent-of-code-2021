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

	fmt.Printf("Answer %d\n", traverseMatrixForCost(solutionMatrix, matrix))
}

func traverseMatrixForCost(solutionMatrix, matrix [][]int) int {
	currentI := len(matrix) - 1
	currentJ := len(matrix) - 1
	var totalCost int
	for true {
		currentCost := matrix[currentI][currentJ]
		totalCost += currentCost

		// see if we are on one of the edges
		if currentJ == 0 {
			currentI = currentI - 1
		} else if currentI == 0 {
			currentJ = currentJ - 1
		} else {
			aboveVal := solutionMatrix[currentI-1][currentJ]
			leftVal := solutionMatrix[currentI][currentJ-1]
			if aboveVal < leftVal {
				currentI = currentI - 1
			} else {
				currentJ = currentJ - 1
			}
		}

		if currentI == 0 && currentJ == 0 {
			break
		}
	}
	return totalCost
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
