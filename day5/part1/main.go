package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

type point struct {
	x int
	y int
}

type segment struct {
	start *point
	end   *point
}

func main() {
	m := initializeMatrix()

	segments := readInput()
	for _, seg := range segments {
		addToMatrix(m, seg)
	}

	answer := findNumInMatrixGreaterThan(m, 1)
	fmt.Printf("Answer %d\n", answer)
}

func findNumInMatrixGreaterThan(matrix [][]int, greaterThan int) int {
	var result int
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix); j++ {
			if matrix[i][j] > greaterThan {
				result++
			}
		}
	}

	return result
}

func initializeMatrix() [][]int {
	m := make([][]int, 0)
	for i := 0; i < 1000; i++ {
		row := make([]int, 0)
		for j := 0; j < 1000; j++ {
			row = append(row, 0)
		}
		m = append(m, row)
	}
	return m
}

func addToMatrix(matrix [][]int, seg *segment) {
	// For now, only consider horizontal and vertical lines: lines where either x1 = x2 or y1 = y2.
	if seg.start.x != seg.end.x && seg.start.y != seg.end.y {
		return
	}

	if seg.start.x == seg.end.x { // vertical line, iterate over the y's and update
		yStartPos := seg.start.y
		yEndPos := seg.end.y
		if seg.end.y < seg.start.y { // choose the smaller "y"
			yStartPos = seg.end.y
			yEndPos = seg.start.y
		}

		// update all positions along the vertical path
		for y := yStartPos; y <= yEndPos; y++ {
			matrix[seg.start.x][y] = matrix[seg.start.x][y] + 1
		}
	} else {
		xStartPos := seg.start.x
		xEndPos := seg.end.x
		if seg.end.x < seg.start.x { // choose the smaller "x"
			xStartPos = seg.end.x
			xEndPos = seg.start.x
		}

		// update all positions along the horizontal path
		for x := xStartPos; x <= xEndPos; x++ {
			matrix[x][seg.start.y] = matrix[x][seg.start.y] + 1
		}
	}
}

func readInput() []*segment {
	filepath := "/Users/jeff.martin@getweave.com/go/src/github.com/jmartin127/advent-of-code-2021/day5/input.txt"
	list := helpers.ReadFile(filepath)

	segments := make([]*segment, 0)
	for _, line := range list {
		segmentStrings := strings.Split(line, " -> ") // 0,9 -> 5,9
		startSegString := segmentStrings[0]
		endSegString := segmentStrings[1]

		startPoint := parsePoint(startSegString)
		endPoint := parsePoint(endSegString)

		s := segment{
			start: startPoint,
			end:   endPoint,
		}
		segments = append(segments, &s)
	}

	return segments
}

// 0,9
func parsePoint(line string) *point {
	parts := strings.Split(line, ",")
	xCoord, _ := strconv.Atoi(parts[0])
	yCoord, _ := strconv.Atoi(parts[1])

	return &point{
		x: xCoord,
		y: yCoord,
	}
}
