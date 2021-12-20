package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

const PADDING = 40

func main() {
	filepath := "input.txt"
	list := helpers.ReadFile(filepath)

	algo, m := parseInput(list)
	fmt.Println("Original")
	helpers.PrintIntMatrixAsHashes(m)

	paddedMatrix := addPaddingToMatrix(m, PADDING, 0)
	fmt.Println("Padded")
	helpers.PrintIntMatrixAsHashes(paddedMatrix)

	enhanced := enhanceImage(paddedMatrix, algo, PADDING, len(m))
	fmt.Println("Enhanced")
	helpers.PrintIntMatrixAsHashes(enhanced)

	final := enhanceImage(enhanced, algo, PADDING/2, len(m)+2)
	helpers.PrintIntMatrixAsHashes(final)
}

func enhanceImage(paddedMatrix [][]int, algo []bool, padding, original int) [][]int {
	enhancedSize := original + padding
	border := padding / 2

	enhancedImage := helpers.NewIntMatrixOfSize(enhancedSize, enhancedSize, 0)
	var count int
	for i := 0; i < enhancedSize; i++ {
		for j := 0; j < enhancedSize; j++ {
			if findOutputValue(paddedMatrix, algo, i+border, j+border) {
				count++
				enhancedImage[i][j] = 1
			}
		}
	}
	fmt.Printf("Count %d\n", count)
	return enhancedImage
}

func makeKey(i, j int) string {
	return fmt.Sprintf("%d,%d", i, j)
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
