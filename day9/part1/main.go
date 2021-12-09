package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

func main() {
	filepath := "input.txt"
	list := helpers.ReadFile(filepath)

	matrix := make([][]int, 0)
	for _, line := range list {
		row := make([]int, 0)
		vals := strings.Split(line, "")
		for _, stringVal := range vals {
			v, _ := strconv.Atoi(stringVal)
			row = append(row, v)
		}
		matrix = append(matrix, row)
	}

	answer := 0
	for i, row := range matrix {
		for j, val := range row {
			if isLowPoint(matrix, i, j, val) {
				//fmt.Printf("Val %d\n", val)
				answer += (val + 1)
				fmt.Printf("%s", ".")
			} else {
				fmt.Printf("%d", val)
			}
		}
		fmt.Println()
	}

	fmt.Printf("Answer %d\n", answer)
}

func isLowPoint(matrix [][]int, i, j int, val int) bool {
	// check above
	if i > 0 && matrix[i-1][j] <= val {
		return false
	}

	// check below
	if i < len(matrix)-1 && matrix[i+1][j] <= val {
		return false
	}

	// check left
	if j > 0 && matrix[i][j-1] <= val {
		return false
	}

	// check right
	if j < len(matrix[0])-1 && matrix[i][j+1] <= val {
		return false
	}

	return true
}
