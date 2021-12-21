package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

const NUM_ENHANCEMENTS = 50
const PADDING = NUM_ENHANCEMENTS + 2

func main() {
	filepath := "input.txt"
	list := helpers.ReadFile(filepath)

	algo, m := parseInput(list)
	paddedMatrix := addPaddingToMatrix(m, PADDING, 0)
	answer := enhanceNtimes(NUM_ENHANCEMENTS, paddedMatrix, algo, len(m), PADDING)
	fmt.Printf("Answer %d\n", answer)
}

func enhanceNtimes(n int, paddedMatrix [][]int, algo []bool, originalSize int, paddingWidth int) int {
	currentMatrix := paddedMatrix
	prevSize := originalSize
	var answer int
	for i := 0; i < n; i++ {
		currentMatrix, prevSize, answer = enhanceImage(currentMatrix, algo, prevSize)
		paddingWidth--
		alternatePaddingValues(currentMatrix, paddingWidth)
		fmt.Printf("Iteration %d Count %d\n", i, answer)
	}

	return answer
}

func isEven(val int) bool {
	return val%2 == 0
}

func enhanceImage(paddedMatrix [][]int, algo []bool, prevSize int) ([][]int, int, int) {
	helpers.PrintIntMatrixAsHashes(paddedMatrix)

	// 1. make a copy of the padded matrix
	newPadded := copyMatrix(paddedMatrix)

	// 2. determine size of enhanced image
	enhancedSize := prevSize + 2

	// 3. fill out the enhanced image.
	iterationStart := (len(paddedMatrix) - enhancedSize) / 2

	var count int
	for i := iterationStart; i < iterationStart+enhancedSize; i++ {
		for j := iterationStart; j < iterationStart+enhancedSize; j++ {
			if findOutputValue(paddedMatrix, algo, i, j) {
				count++
				newPadded[i][j] = 1
			} else {
				newPadded[i][j] = 0
			}
		}
	}
	return newPadded, enhancedSize, count
}

// swaps the 0/1 values for the cells within the padding area
func alternatePaddingValues(matrix [][]int, paddingWidth int) {
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix); j++ {
			if isInPaddingArea(i, j, paddingWidth, len(matrix)) {
				if matrix[i][j] == 0 {
					matrix[i][j] = 1
				} else {
					matrix[i][j] = 0
				}
			}
		}
	}
}

func isInPaddingArea(i, j int, paddingWidth int, matrixSize int) bool {
	isInterior := i >= paddingWidth && i < matrixSize-paddingWidth && j >= paddingWidth && j < matrixSize-paddingWidth
	return !isInterior
}

func copyMatrix(matrix [][]int) [][]int {
	copy := helpers.NewIntMatrixOfSize(len(matrix), len(matrix), 0)
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix); j++ {
			copy[i][j] = matrix[i][j]
		}
	}
	return copy
}

func findOutputValue(matrix [][]int, algo []bool, xPos, yPos int) bool {
	vals := make([]bool, 0)

	// check all 9 items
	for i := xPos - 1; i < xPos-1+3; i++ {
		for j := yPos - 1; j < yPos-1+3; j++ {
			vals = append(vals, lookupInMatrix(matrix, i, j))
		}
	}

	//fmt.Printf("Vals %+v\n", vals)
	returnV := valueFromAlgo(vals, algo)
	//fmt.Printf("return %t\n", returnV)
	return returnV
}

func lookupInMatrix(matrix [][]int, i, j int) bool {
	return matrix[i][j] == 1
}

func isInBounds(size int, xPos, yPos int) bool {
	return xPos >= 0 && xPos <= size-1 && yPos >= 0 && yPos <= size-1
}

func valueFromAlgo(pixels []bool, algo []bool) bool {
	val := pixelsToDecimal(pixels)
	//fmt.Printf("Decimal %d\n", val)

	return algo[val]
}

func pixelsToDecimal(pixels []bool) int {
	binaryStr := ""
	for _, p := range pixels {
		if p {
			binaryStr += "1"
		} else {
			binaryStr += "0"
		}
	}
	return binaryToDecimal(binaryStr)
}

func binaryToDecimal(binary string) int {
	output, _ := strconv.ParseInt(binary, 2, 64)
	return int(output)
}

func addPaddingToMatrix(matrix [][]int, numToPad int, defaultVal int) [][]int {
	newSize := len(matrix) + 2*numToPad
	paddedMatrix := helpers.NewIntMatrixOfSize(newSize, newSize, defaultVal)
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix); j++ {
			paddedMatrix[i+numToPad][j+numToPad] = matrix[i][j]
		}
	}

	return paddedMatrix
}

/*
#..#.
#....
##..#
..#..
..###
*/
func parseInput(list []string) ([]bool, [][]int) {
	algoStr := list[0]

	matrix := helpers.NewIntMatrixOfSize(len(list[2]), len(list[2]), 0)
	var rowCount int
	for i := 2; i < len(list); i++ {
		chars := strings.Split(list[i], "")
		for j, c := range chars {
			if c == "#" {
				matrix[rowCount][j] = 1
			} else {
				matrix[rowCount][j] = 0
			}
		}
		rowCount++
	}

	algoVals := strings.Split(algoStr, "")
	algo := make([]bool, 0)
	for _, v := range algoVals {
		if v == "#" {
			algo = append(algo, true)
		} else {
			algo = append(algo, false)
		}
	}

	return algo, matrix
}
