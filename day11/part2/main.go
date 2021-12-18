package main

import (
	"fmt"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

type cell struct {
	val     int
	flashed bool
}

func main() {
	matrixOfInts := helpers.ReadFileAsMatrixOfInts("input.txt")
	matrix := convertMatrixToCells(matrixOfInts)
	answer := runNSteps(matrix)
	fmt.Printf("Answer %d\n", answer)
}

func runNSteps(matrix [][]*cell) int {
	for i := 0; i < 100000; i++ {
		currentMatrix := matrix
		newMatrix, numFlashed := runStep(currentMatrix)
		currentMatrix = newMatrix
		if numFlashed == len(matrix)*len(matrix) {
			return i + 1
		}
	}
	return -1
}

func runStep(matrix [][]*cell) ([][]*cell, int) {
	// First, the energy level of each octopus increases by 1. (let's also reset the flashers)
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix); j++ {
			matrix[i][j].val = matrix[i][j].val + 1
			matrix[i][j].flashed = false
		}
	}

	newMatrix, numFlashed := applyFlash(matrix)
	resetFlashers(newMatrix)
	return newMatrix, numFlashed
}

// Finally, any octopus that flashed during this step has its energy level set to 0, as it used all of its energy to flash.
func resetFlashers(matrix [][]*cell) {
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix); j++ {
			if matrix[i][j].flashed {
				matrix[i][j].val = 0
			}
			matrix[i][j].flashed = false
		}
	}
}

func applyFlash(matrix [][]*cell) ([][]*cell, int) {
	numFlashed := 0

	for true {
		var octopusFlashed bool
		for i := 0; i < len(matrix); i++ {
			for j := 0; j < len(matrix); j++ {
				if matrix[i][j].val > 9 && !matrix[i][j].flashed { // Then, any octopus with an energy level greater than 9 flashes. (An octopus can only flash at most once per step.)
					increaseAdjacent(i, j, matrix)
					matrix[i][j].flashed = true
					octopusFlashed = true
					numFlashed++
				}
			}
		}

		// This process continues as long as new octopuses keep having their energy level increased beyond 9.
		if !octopusFlashed {
			break
		}
	}

	return matrix, numFlashed
}

// This increases the energy level of all adjacent octopuses by 1, including
// octopuses that are diagonally adjacent.
func increaseAdjacent(i, j int, m [][]*cell) {
	incrementPos(i+1, j, m)   // up
	incrementPos(i-1, j, m)   // down
	incrementPos(i, j-1, m)   // left
	incrementPos(i, j+1, m)   // right
	incrementPos(i+1, j-1, m) // up-left
	incrementPos(i+1, j+1, m) // up-right
	incrementPos(i-1, j-1, m) // down-left
	incrementPos(i-1, j+1, m) // down-right
}

func incrementPos(i, j int, m [][]*cell) {
	if isInBounds(i, j, len(m)) {
		m[i][j].val = m[i][j].val + 1
	}
}

func isInBounds(i, j, size int) bool {
	return i >= 0 && i <= size-1 && j >= 0 && j <= size-1
}

func newMatrix(size int) [][]*cell {
	matrix := make([][]*cell, 0)
	for i := 0; i < size; i++ {
		row := make([]*cell, 0)
		for j := 0; j < size; j++ {
			row = append(row, &cell{})
		}
		matrix = append(matrix, row)
	}
	return matrix
}

func convertMatrixToCells(m [][]int) [][]*cell {
	matrix := newMatrix(len(m))
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m); j++ {
			matrix[i][j].val = m[i][j]
		}
	}
	return matrix
}

func copyMatrix(m [][]*cell) [][]*cell {
	matrix := newMatrix(len(m))
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m); j++ {
			matrix[i][j].val = m[i][j].val
		}
	}
	return matrix
}

func printMatrix(matrix [][]*cell) {
	for i := 0; i < len(matrix); i++ {
		row := matrix[i]
		for _, v := range row {
			fmt.Printf("%d", v.val)
		}
		fmt.Println()
	}
}
