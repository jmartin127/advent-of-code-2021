package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

func main() {
	filepath := "input.txt"
	list := helpers.ReadFile(filepath)

	matrix := make([][]int, 0)
	basins := make([][]int, 0)
	for _, line := range list {
		row := make([]int, 0)
		basinRow := make([]int, 0)
		vals := strings.Split(line, "")
		for _, stringVal := range vals {
			v, _ := strconv.Atoi(stringVal)
			row = append(row, v)
			basinRow = append(basinRow, -1)
		}
		matrix = append(matrix, row)
		basins = append(basins, basinRow)
	}

	nextBasinVal := 0
	for i, row := range matrix {
		for j := range row {
			if basins[i][j] == -1 {
				nextBasinVal++
			}
			greedilyFill(matrix, basins, i, j, nextBasinVal)
		}
	}

	result := make(map[int]int, 0)
	for _, row := range basins {
		for _, val := range row {
			fmt.Printf("%d", val)
			if val != -1 {
				result[val]++
			}
		}
		fmt.Println()
	}

	// get the top values
	values := make([]int, 0)
	for _, v := range result {
		values = append(values, v)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(values)))

	fmt.Printf("values %+v", values)

	// top 3
	top1 := values[0]
	top2 := values[1]
	top3 := values[2]

	fmt.Printf("Answer %d %d %d\n", top1, top2, top3)
	fmt.Printf("Final %d\n", top1*top2*top3)
}

func greedilyFill(matrix [][]int, basins [][]int, i, j int, nextBasinVal int) {
	// base case
	if i == -1 || j == -1 || i > len(matrix)-1 || j > len(matrix[0])-1 || valIsNotFillable(matrix[i][j], basins[i][j]) {
		return
	}

	// update the current position
	basins[i][j] = nextBasinVal

	// go right
	greedilyFill(matrix, basins, i, j+1, nextBasinVal)

	// go left
	greedilyFill(matrix, basins, i, j-1, nextBasinVal)

	// go up
	greedilyFill(matrix, basins, i-1, j, nextBasinVal)

	// go down
	greedilyFill(matrix, basins, i+1, j, nextBasinVal)
}

func valIsNotFillable(val int, basin int) bool {
	return val == 9 || basin != -1
}
