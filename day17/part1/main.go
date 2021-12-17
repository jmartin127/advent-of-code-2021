package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

var MAX_STEPS = 600

type targetArea struct {
	x1 int
	x2 int
	y1 int
	y2 int
}

type point struct {
	x int
	y int
}

func (ta *targetArea) isInside(x, y int) bool {
	if x >= min(ta.x1, ta.x2) && x <= max(ta.x1, ta.x2) && y >= min(ta.y1, ta.y2) && y <= max(ta.y1, ta.y2) {
		return true
	}
	return false
}

func main() {
	filepath := "input.txt"
	list := helpers.ReadFile(filepath)

	ta := parseTargetArea(list[0])
	fmt.Printf("Target %+v\n", ta)

	answer := runMultipleSimultations(ta)
	fmt.Printf("Answer %d\n", answer)
}

func printIt(ta *targetArea) {
	fmt.Printf("Path:\n")
	matrixSize := 400
	startX := 18
	startY := 6
	matrixOffset := 150
	m := helpers.NewStringMatrixOfSize(matrixSize, matrixSize, ".")
	_, _, path := simulateTrajectory(startX, startY, ta, true)
	m[matrixOffset][matrixOffset] = "S"
	for _, p := range path {
		fmt.Printf("Adding %+v\n", p)
		m[p.y+matrixOffset][p.x+matrixOffset] = "#"
	}
	for i := ta.x1; i <= ta.x2; i++ {
		for j := ta.y1; j <= ta.y2; j++ {
			if m[j+matrixOffset][i+matrixOffset] == "." {
				m[j+matrixOffset][i+matrixOffset] = "T"
			}
		}
	}
	helpers.PrintStringMatrix(m)
}

func runMultipleSimultations(ta *targetArea) int {
	overallHighestY := -100000000000
	for i := -1000; i < 1000; i++ {
		for j := -1000; j < 1000; j++ {
			highestY, hit, _ := simulateTrajectory(i, j, ta, false)
			if hit {
				overallHighestY = max(highestY, overallHighestY)
			}
		}
	}

	return overallHighestY
}

/* The probe's x position increases by its x velocity.
The probe's y position increases by its y velocity.
Due to drag, the probe's x velocity changes by 1 toward the value 0; that is, it decreases by 1 if it is greater than 0, increases by 1 if it is less than 0, or does not change if it is already 0.
Due to gravity, the probe's y velocity decreases by 1.
*/
func simulateTrajectory(initX, initY int, ta *targetArea, print bool) (int, bool, []*point) {
	currentX := 0
	currentY := 0

	xVelocity := initX
	yVelocity := initY

	highestY := -100000000000
	var steps int
	path := make([]*point, 0)
	for true {
		// move one position
		currentX += xVelocity
		currentY += yVelocity

		xVelocity = moveX(xVelocity)
		yVelocity = yVelocity - 1

		highestY = max(highestY, currentY)

		if print {
			fmt.Printf("Checking x=%d, y=%d\n", currentX, currentY)
			path = append(path, &point{x: currentX, y: currentY})
		}

		// check if we hit the target area
		if ta.isInside(currentX, currentY) {
			return highestY, true, path
		}

		// don't iterate for too long
		steps++
		if steps >= MAX_STEPS {
			break
		}
	}

	return highestY, false, path
}

func moveX(currentX int) int {
	if currentX > 0 {
		return currentX - 1
	} else if currentX < 0 {
		return currentX + 1
	} else {
		return currentX
	}
}

// target area: x=20..30, y=-10..-5
func parseTargetArea(line string) *targetArea {
	parts := strings.Split(line, ": ")
	moreParts := strings.Split(parts[1], ", ")

	x1, x2 := getCoords(moreParts[0])
	y1, y2 := getCoords(moreParts[1])

	return &targetArea{
		x1: x1,
		x2: x2,
		y1: y1,
		y2: y2,
	}
}

// x=20..30
func getCoords(line string) (int, int) {
	parts := strings.Split(line, "=")
	numbers := strings.Split(parts[1], "..")

	numOne, _ := strconv.Atoi(numbers[0])
	numTwo, _ := strconv.Atoi(numbers[1])

	return numOne, numTwo
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
