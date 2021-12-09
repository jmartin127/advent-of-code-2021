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
			basinRow = append(basinRow, 9)
		}
		matrix = append(matrix, row)
		basins = append(basins, basinRow)
	}

	nextBasinVal := 0
	for i, row := range matrix {
		for j, val := range row {
			newBasin := setVal(matrix, basins, i, j, val, nextBasinVal)
			if newBasin {
				nextBasinVal++
				//fmt.Println("new basin")
			}
		}
	}

	result := make(map[int]int, 0)
	for _, row := range basins {
		for _, val := range row {
			fmt.Printf("%d", val)
			if val != 9 {
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

func setVal(matrix [][]int, basins [][]int, i, j int, val int, nextBasinVal int) bool {
	if matrix[i][j] == 9 {
		return false
	}

	// check above
	if i > 0 {
		if matrix[i-1][j] != 9 {
			basins[i][j] = basins[i-1][j]
			return false
		}
	}

	// check left
	if j > 0 {
		if matrix[i][j-1] != 9 {
			basins[i][j] = basins[i][j-1]
			return false
		}
	}

	// check upper-right diagonal
	if i > 0 && j < len(matrix[0])-1 {
		if matrix[i-1][j+1] != 9 {
			if matrix[i-1][j] == 9 && matrix[i][j+1] == 9 {
				// skip this basin since it is just a leaky diagonal basin
			} else {
				basins[i][j] = basins[i-1][j+1]
				return false
			}
		}
	}

	basins[i][j] = nextBasinVal
	return true
}
