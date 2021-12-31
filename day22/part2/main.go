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
	z int
}

type instruction struct {
	isOn   bool
	xStart int
	xEnd   int
	yStart int
	yEnd   int
	zStart int
	zEnd   int
}

func (i *instruction) corners() []*point {
	return []*point{
		{x: i.xStart, y: i.yStart, z: i.zStart},
		{x: i.xStart, y: i.yEnd, z: i.zStart},
		{x: i.xEnd, y: i.yEnd, z: i.zStart},
		{x: i.xEnd, y: i.yStart, z: i.zStart},
		{x: i.xStart, y: i.yStart, z: i.zEnd},
		{x: i.xStart, y: i.yEnd, z: i.zEnd},
		{x: i.xEnd, y: i.yEnd, z: i.zEnd},
		{x: i.xEnd, y: i.yStart, z: i.zEnd},
	}
}

func (i *instruction) containsPoint(p *point) bool {
	return p.x >= i.xStart && p.x <= i.xEnd && p.y >= i.yStart && p.y <= i.yEnd && p.z >= i.zStart && p.z <= i.zEnd
}

func main() {
	instructions := readInstructions()
	fmt.Printf("Number of instructions %+v\n", len(instructions))

	distMatrix := compareCubes(instructions)
	helpers.PrintIntMatrix(distMatrix)
	countNumOverlapsPerCube(distMatrix)

	// OK, after thinking about it... the plan will be to:
	// 1. Compute how much each ON cube uniquely contributes
	// 2. To do this, start with the first ON cube and do the following:
	//     a. Compute the volume of the cube (as a starting point)
	//     b. Cubes BEFORE this one in the list:
	//       i.  If it is an OFF cube, ignore it.
	//       ii. If it is an ON cube, subtract the volume that the 2 cubes overlap
	//     c. Cubes AFTER this one in the list:
	//       i.  If it is an OFF cube, subtact the volume that the 2 cubes overlap
	//       ii. If it is an ON cube, ignore it.
	// 3. Add up how much each cube uniquely contributes to get the final answer!

	// 4. Determine how to proceed, depending on what the overlaps look like
	/*
		1 corner = 1422
		2 corner = 761
		3 corner = 0
		4 corner = 142
		5 corner = 0
		6 corner = 0
		7 corner = 0
		8 corner = 20
	*/

	// Result: Turns out that there are only 4 types of overlaps:
	// 1 corner: 1 corner contained
	// 2 corner: 1 edge contained
	// 4 corner: 1 half contained
	// 5 corner: completely contained
}

/*
Counts by position:
i=0, count=13
i=1, count=13
i=2, count=11
i=3, count=13
i=4, count=13
i=5, count=14
i=6, count=13
i=7, count=13
i=8, count=13
i=9, count=9
i=10, count=4
i=11, count=14
i=12, count=4
i=13, count=12
i=14, count=1
i=15, count=15
i=16, count=8
i=17, count=12
i=18, count=11
i=19, count=14
i=20, count=12
i=21, count=5
*/
func countNumOverlapsPerCube(distMatrix [][]int) {
	result := make([][]string, 0)
	for i := 0; i < len(distMatrix); i++ {
		overlappingCubes := make([]string, 0)
		for j := 0; j < len(distMatrix); j++ {
			if distMatrix[i][j] > 0 {
				overlappingCubes = append(overlappingCubes, fmt.Sprintf("%d:%d", j, distMatrix[i][j]))
			}
		}
		result = append(result, overlappingCubes)
	}
	fmt.Println("Overlapping cubes by position:")
	for i, v := range result {
		fmt.Printf("i=%d, overlapping-cubes=%+v\n", i, v)
	}
}

// 2. Do pair-wise comparison of all cubes
// 3. Notate in a distance matrix which cubes overlap with which other cubes., and what type of overlap they have.
func compareCubes(instructions []*instruction) [][]int {
	result := helpers.NewIntMatrixOfSize(len(instructions), len(instructions), 0)
	for i := 0; i < len(instructions); i++ {
		for j := 0; j < len(instructions); j++ {
			if i == j {
				continue
			}
			fmt.Printf("Setting value for i=%d, j=%d\n", i, j)
			insI := instructions[i]
			insJ := instructions[j]
			numContained := findOverlap(insI, insJ)
			result[i][j] = numContained
			if numContained == 8 {
				fmt.Printf("i=%d,j=%d,val=%d\n", i, j, numContained)
				fmt.Printf("\ti=%+v\n\tj=%+v\n", insI, insJ)
			}
		}
	}
	return result
}

// 1. Write a method which determines if cubes overlap
//   a. For first cube, obtain the coordinates of every corner
//   b. Compare each of the 4 corners with the other cube, and determine how many corners are contained within the otehr cube.
func findOverlap(a, b *instruction) int {
	var numCornersContained int
	for _, p := range a.corners() {
		if b.containsPoint(p) {
			numCornersContained++
		}
	}
	return numCornersContained
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
