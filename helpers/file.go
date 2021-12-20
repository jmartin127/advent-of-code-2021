package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadFile(filepath string) []string {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	result := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)
	}

	return result
}

func ReadFileAsInts(filepath string) []int {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	result := make([]int, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, readLine(line))
	}

	return result
}

func ReadFileAsMatrixOfInts(filepath string) [][]int {
	list := ReadFile(filepath)

	matrix := make([][]int, 0)
	for _, line := range list {
		row := make([]int, 0)
		for _, v := range strings.Split(line, "") {
			i, _ := strconv.Atoi(v)
			row = append(row, i)
		}
		matrix = append(matrix, row)
	}

	return matrix
}

func ReadSingleLineFileAsInts(filepath string) []int {
	// convert to ints
	firstLineParts := ReadSingleLineFileAsStrings(filepath)
	result := make([]int, 0)
	for _, val := range firstLineParts {
		intVal, _ := strconv.Atoi(val)
		result = append(result, intVal)
	}

	return result
}

func ReadSingleLineFileAsStrings(filepath string) []string {
	list := ReadFile(filepath)
	return strings.Split(list[0], ",")
}

func MaxFromList(list []int) int {
	var max int
	for _, val := range list {
		if val > max {
			max = val
		}
	}
	return max
}

func NewIntMatrixOfSize(numRows int, numCols int, initVal int) [][]int {
	result := make([][]int, 0)
	for i := 0; i < numRows; i++ {
		row := make([]int, 0)
		for j := 0; j < numCols; j++ {
			row = append(row, initVal)
		}
		result = append(result, row)
	}

	return result
}

func NewStringMatrixOfSize(numRows int, numCols int, initVal string) [][]string {
	result := make([][]string, 0)
	for i := 0; i < numRows; i++ {
		row := make([]string, 0)
		for j := 0; j < numCols; j++ {
			row = append(row, initVal)
		}
		result = append(result, row)
	}

	return result
}

func PrintIntMatrix(matrix [][]int) {
	for i := len(matrix) - 1; i >= 0; i-- {
		row := matrix[i]
		for _, v := range row {
			fmt.Printf("%d ", v)
		}
		fmt.Println()
	}
}

func PrintIntMatrixAsHashes(matrix [][]int) {
	for i := 0; i < len(matrix); i++ {
		row := matrix[i]
		for _, v := range row {
			if v == 0 {
				fmt.Printf("%s", ".")
			} else {
				fmt.Printf("%s", "#")
			}
		}
		fmt.Println()
	}
}

const (
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

func PrintStringMatrix(matrix [][]string) {
	for i := len(matrix) - 1; i >= 0; i-- {
		row := matrix[i]
		for _, v := range row {
			if v == "#" {
				fmt.Printf(WarningColor, v)
			} else if v == "T" {
				fmt.Printf(ErrorColor, v)
			} else if v == "S" {
				fmt.Printf(DebugColor, v)
			} else {
				fmt.Printf("%s", v)
			}
		}
		fmt.Println()
	}
}

func readLine(input string) int {
	i, err := strconv.Atoi(input)
	if err != nil {
		panic(err)
	}
	return i
}
