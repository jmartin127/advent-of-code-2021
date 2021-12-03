package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

type set struct {
	rows [][]int
}

func main() {
	s := readInput()

	result, _ := applyFilter(s, 0, true)
	fmt.Printf("Result: %+v\n", result)

	result, _ = applyFilter(s, 0, false)
	fmt.Printf("Result: %+v\n", result)

}

func readInput() *set {
	filepath := "/Users/jeff.martin@getweave.com/go/src/github.com/jmartin127/advent-of-code-2021/day3/input.txt"
	list := helpers.ReadFile(filepath)
	s := &set{
		rows: make([][]int, 0),
	}
	for _, v := range list {
		chars := strings.Split(v, "") // 001011100010
		row := make([]int, 0)
		for _, c := range chars {
			intVal, _ := strconv.Atoi(c)
			row = append(row, intVal)
		}
		s.rows = append(s.rows, row)
	}
	return s
}

func countsAtIndex(s *set, i int) (int, int) {
	countZero := 0
	for _, row := range s.rows {
		if row[i] == 0 {
			countZero++
		}
	}

	return countZero, len(s.rows) - countZero
}

func applyFilter(s *set, index int, isOxygenFilter bool) (*set, int) {
	if len(s.rows) == 1 {
		return s, index
	}

	countZero, countOne := countsAtIndex(s, index)

	var keepOnes bool
	if countZero > countOne {
		keepOnes = false
	} else if countOne > countZero {
		keepOnes = true
	} else {
		keepOnes = true
	}

	if !isOxygenFilter {
		keepOnes = !keepOnes
	}

	// iterate and keep those that we should keep
	result := &set{
		rows: make([][]int, 0),
	}
	for _, row := range s.rows {
		if keepOnes {
			if row[index] == 1 {
				result.rows = append(result.rows, row)
			}
		} else {
			if row[index] == 0 {
				result.rows = append(result.rows, row)
			}
		}
	}

	return applyFilter(result, index+1, isOxygenFilter)
}
