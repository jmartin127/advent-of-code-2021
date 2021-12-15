package main

import (
	"fmt"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

const MAX_VAL = 4294967295 // max of unit32

type cell struct {
	x        int
	y        int
	distance int
}

func main() {
	filepath := "input.txt"
	matrix := helpers.ReadFileAsMatrixOfInts(filepath)

	largerMatrix := makeMatrixLarger(matrix)

	fmt.Printf("Answer %d\n", findMinimumCost(largerMatrix)-largerMatrix[0][0])
}

func makeMatrixLarger(matrix [][]int) [][]int {
	numTimesLarger := 5
	newSize := len(matrix) * numTimesLarger
	result := helpers.NewIntMatrixOfSize(newSize, newSize, 0)

	// first initialize with existing values
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix); j++ {
			result[i][j] = matrix[i][j]
		}
	}

	// fill in the remaining values
	for i := 0; i < len(result); i++ {
		for j := 0; j < len(result); j++ {
			val := result[i][j]
			if val == 0 { // some are already filled in
				var priorVal int
				if j-len(matrix) >= 0 {
					priorVal = result[i][j-len(matrix)]
				} else {
					priorVal = result[i-len(matrix)][j]
				}
				newVal := priorVal + 1
				if newVal == 10 {
					newVal = 1
				}
				result[i][j] = newVal
			}
		}
	}

	return result
}

func isInside(i, j, size int) bool {
	return i >= 0 && i < size && j >= 0 && j < size
}

// Dikjkstra's algorithm
// Reference: https://www.geeksforgeeks.org/minimum-cost-path-left-right-bottom-moves-allowed/
func findMinimumCost(matrix [][]int) int {
	distanceMatrix := helpers.NewIntMatrixOfSize(len(matrix), len(matrix), MAX_VAL)
	dx := []int{-1, 0, 1, 0}
	dy := []int{0, 1, 0, -1}
	set := []*cell{&cell{x: 0, y: 0, distance: 0}}
	distanceMatrix[0][0] = matrix[0][0]

	for len(set) > 0 {
		var k *cell
		k, set = popFirst(set)

		for i := 0; i < 4; i++ {
			x := k.x + dx[i]
			y := k.y + dy[i]

			if !isInside(x, y, len(matrix)) {
				continue
			}

			if distanceMatrix[x][y] > (distanceMatrix[k.x][k.y] + matrix[x][y]) {
				if distanceMatrix[x][y] != MAX_VAL {
					set = removeFromSet(set, x, y, distanceMatrix[x][y])
				}

				distanceMatrix[x][y] = distanceMatrix[k.x][k.y] + matrix[x][y]
				set = append(set, &cell{x: x, y: y, distance: distanceMatrix[x][y]})
			}
		}
	}

	return distanceMatrix[len(matrix)-1][len(matrix)-1]
}

func removeFromSet(set []*cell, x, y, distance int) []*cell {
	for i, c := range set {
		if c.x == x && c.y == y && c.distance == distance {
			return removeFromSlice(set, i)
		}
	}

	return set
}

func removeFromSlice(set []*cell, pos int) []*cell {
	return append(set[:pos], set[pos+1:]...)
}

func popFirst(set []*cell) (*cell, []*cell) {
	return set[0], set[1:]
}

func minimum(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}
