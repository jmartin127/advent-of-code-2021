package main

import (
	"fmt"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

func main() {
	filepath := "day25/input.txt"
	list := helpers.ReadFile(filepath)
	ocean := readInput(list)

	var numSteps int
	for true {
		numSteps++
		newOcean, numMoved := applyStep(ocean)
		ocean = newOcean
		if numMoved == 0 {
			break
		}
	}
	fmt.Printf("Answer %d\n", numSteps)
}

func printMatrix(matrix [][]string) {
	for _, row := range matrix {
		for _, v := range row {
			fmt.Printf("%s", v)
		}
		fmt.Println("")
	}
	fmt.Println()
}

func copyMatrix(matrix [][]string) [][]string {
	result := make([][]string, 0)
	for _, row := range matrix {
		newRow := make([]string, 0)
		for _, v := range row {
			newRow = append(newRow, v)
		}
		result = append(result, newRow)
	}
	return result
}

func applyStep(matrix [][]string) ([][]string, int) {
	new := copyMatrix(matrix)

	var numMoved int

	// move east facing
	for i, row := range matrix {
		for j, v := range row {
			if v == ">" {
				// check if the one to the right is open, including wrapping
				if matrix[i][nextIndexToCheck(j+1, len(row))] == "." {
					new[i][nextIndexToCheck(j+1, len(row))] = ">"
					new[i][j] = "."
					numMoved++
				}
			}
		}
	}
	matrix = copyMatrix(new)

	// move south facing
	for i, row := range matrix {
		for j, v := range row {
			if v == "v" {
				// check if the one below is open, including wrapping
				if matrix[nextIndexToCheck(i+1, len(matrix))][j] == "." {
					new[nextIndexToCheck(i+1, len(matrix))][j] = "v"
					new[i][j] = "."
					numMoved++
				}
			}
		}
	}

	return new, numMoved
}

// handle wrapping
func nextIndexToCheck(i int, matrixSize int) int {
	if i > matrixSize-1 {
		return 0
	} else {
		return i
	}
}

/*
...>...
.......
......>
v.....>
......>
.......
..vvv..
*/
func readInput(list []string) [][]string {
	result := make([][]string, 0)
	for _, line := range list {
		vals := strings.Split(line, "")
		result = append(result, vals)
	}
	return result
}
