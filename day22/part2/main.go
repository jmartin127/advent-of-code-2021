package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

type instruction struct {
	isOn   bool
	xStart int
	xEnd   int
	yStart int
	yEnd   int
	zStart int
	zEnd   int
}

func main() {
	instructions := readInstructions()
	fmt.Printf("Number of instructions %+v\n", len(instructions))

	// Plan of attack
	// 1. Write a method which determines if cubes overlap
	//   a. For first cube, obtain the coordinates of every corner
	//   b. Compare each of the 4 corners with the other cube, and determine how many corners are contained within the otehr cube.
	// 2. Do pair-wise comparison of all cubes
	// 3. Notate in a distance matrix which cubes overlap with which other cubes., and what type of overlap they have.
	// 4. Determine how to proceed, depending on what the overlaps look like
}

func readInstructions() []*instruction {
	filepath := "day22/input.txt"
	list := helpers.ReadFile(filepath)

	instructions := make([]*instruction, 0)
	for _, line := range list {
		instruction := parseLine(line)
		instructions = append(instructions, instruction)
	}

	return instructions
}

func processInstruction(ins *instruction, onCubes map[string]bool) {
	if !isInRange(ins) {
		return
	}

	// obtain all "off" instructions that come after this one
}

// The initialization procedure only uses cubes that have x, y, and z positions of at least -50 and at most 50. For now, ignore cubes outside this region.
func isInRange(ins *instruction) bool {
	return true
	//return isPointInRange(ins.xStart, ins.yStart, ins.zStart) && isPointInRange(ins.xEnd, ins.yEnd, ins.zEnd)
}

func isPointInRange(x, y, z int) bool {
	return true
	return x >= -50 && x <= 50 && y >= -50 && y <= 50 && z >= -50 && z <= 50
}

func makeKey(x, y, z int) string {
	return fmt.Sprintf("%d,%d,%d", x, y, z)
}

// on x=10..12,y=10..12,z=10..12
func parseLine(input string) *instruction {
	parts := strings.Split(input, " ")
	onOff := parts[0]

	var isOn bool
	if onOff == "on" {
		isOn = true
	}

	coords := strings.Split(parts[1], ",")
	_, xStart, xEnd := parseCoord(coords[0])
	_, yStart, yEnd := parseCoord(coords[1])
	_, zStart, zEnd := parseCoord(coords[2])

	return &instruction{isOn: isOn, xStart: xStart, xEnd: xEnd, yStart: yStart, yEnd: yEnd, zStart: zStart, zEnd: zEnd}
}

// z=10..12
func parseCoord(input string) (string, int, int) {
	parts := strings.Split(input, "=")
	direction := parts[0]
	theRange := strings.Split(parts[1], "..")

	start, _ := strconv.Atoi(theRange[0])
	end, _ := strconv.Atoi(theRange[1])

	return direction, start, end
}
