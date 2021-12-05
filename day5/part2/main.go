package main

import (
	"fmt"
	"math"
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
	for x := 0; x < len(matrix); x++ {
		for y := 0; y < len(matrix); y++ {
			if isBetweenSegment(x, y, seg) {
				matrix[x][y] = matrix[x][y] + 1
			}
		}
	}
}

// Finally realized there is a mathematical formula for this!
// reference: https://stackoverflow.com/questions/328107/how-can-you-determine-a-point-is-between-two-other-points-on-a-line-segment
func isBetweenSegment(x int, y int, seg *segment) bool {
	a := point{x: seg.start.x, y: seg.start.y}
	b := point{x: seg.end.x, y: seg.end.y}
	c := point{x: x, y: y}

	crossproduct := (c.y-a.y)*(b.x-a.x) - (c.x-a.x)*(b.y-a.y)

	if int(math.Abs(float64(crossproduct))) != 0 {
		return false
	}

	dotproduct := (c.x-a.x)*(b.x-a.x) + (c.y-a.y)*(b.y-a.y)
	if dotproduct < 0 {
		return false
	}

	squaredlengthba := (b.x-a.x)*(b.x-a.x) + (b.y-a.y)*(b.y-a.y)
	if dotproduct > squaredlengthba {
		return false
	}

	return true
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
