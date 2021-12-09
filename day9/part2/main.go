package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

type pos struct {
	val   int
	basin int
}
type basin struct {
	matrix [][]*pos
}

func NewBasin() *basin {
	return &basin{
		matrix: make([][]*pos, 0),
	}
}

func main() {
	filepath := "input.txt"
	list := helpers.ReadFile(filepath)

	// parse input
	b := NewBasin()
	for _, line := range list {
		row := make([]*pos, 0)
		for _, stringVal := range strings.Split(line, "") {
			v, _ := strconv.Atoi(stringVal)
			row = append(row, &pos{val: v, basin: -1})
		}
		b.matrix = append(b.matrix, row)
	}

	// greedily fill each position
	nextBasinVal := 0
	for i, row := range b.matrix {
		for j := range row {
			if b.matrix[i][j].basin == -1 {
				nextBasinVal++
			}
			greedilyFill(b, i, j, nextBasinVal)
		}
	}

	// count how often each occurs
	result := make(map[int]int, 0)
	for _, row := range b.matrix {
		for _, val := range row {
			if val.basin != -1 {
				result[val.basin]++
			}
		}
	}

	// get the top values
	values := make([]int, 0)
	for _, v := range result {
		values = append(values, v)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(values)))

	// top 3
	top1 := values[0]
	top2 := values[1]
	top3 := values[2]
	fmt.Printf("Top 3: %d %d %d\n", top1, top2, top3)
	fmt.Printf("Final: %d\n", top1*top2*top3)
}

func greedilyFill(b *basin, i, j int, nextBasinVal int) {
	// base case
	if i == -1 || j == -1 || i > len(b.matrix)-1 || j > len(b.matrix[0])-1 || valIsNotFillable(b.matrix[i][j].val, b.matrix[i][j].basin) {
		return
	}

	// update the current position
	b.matrix[i][j].basin = nextBasinVal

	// go right
	greedilyFill(b, i, j+1, nextBasinVal)

	// go left
	greedilyFill(b, i, j-1, nextBasinVal)

	// go up
	greedilyFill(b, i-1, j, nextBasinVal)

	// go down
	greedilyFill(b, i+1, j, nextBasinVal)
}

func valIsNotFillable(val int, basin int) bool {
	return val == 9 || basin != -1
}
