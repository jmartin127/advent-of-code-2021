package main

import (
	"fmt"
	"log"
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

func (i *instruction) isEqual(o *instruction) bool {
	return i.isOn == o.isOn &&
		i.xStart == o.xStart && i.xEnd == o.xEnd &&
		i.yStart == o.yStart && i.yEnd == o.yEnd &&
		i.zStart == o.zStart && i.zEnd == o.zEnd
}

// on x=-41..9,y=-7..43,z=-33..15
func (i *instruction) asString() string {
	var result string
	if i.isOn {
		result += "on "
	} else {
		result += "off "
	}
	result += fmt.Sprintf("x=%d..%d,y=%d..%d,z=%d..%d", i.xStart, i.xEnd, i.yStart, i.yEnd, i.zStart, i.zEnd)
	return result
}

// NOTE: Volume can be zero after dividing cubes
func (i *instruction) volume() int {
	if i.xStart > i.xEnd {
		return 0
	}
	if i.yStart > i.yEnd {
		return 0
	}
	if i.zStart > i.zEnd {
		return 0
	}
	return (i.xEnd - i.xStart + 1) * (i.yEnd - i.yStart + 1) * (i.zEnd - i.zStart + 1)
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

	currentInstructions := instructions
	for true {
		foundOffCube, newInstructions := findNextOffInstructionAndApplyToPriorOns(currentInstructions)
		currentInstructions = newInstructions
		if !foundOffCube {
			break
		}
	}
	var total int
	for _, ins := range currentInstructions {
		fmt.Printf("%s\n", ins.asString())
		total += ins.volume()
	}
	fmt.Printf("After only %d ON cubes remain.  Note: still has overlapps we need to resolve: %d\n", len(currentInstructions), total)

	/*
		2. Now we are left with ONLY ON cubes, and can use 1 of 2 strategies:
		  a. Implement an algorithm to find the union volume of all ON cubes
		  b. Use the same strategy to split overlapping ON cubes, and then just add up their volume.
	*/
	// NOTE: Going to try my 2a approach, since 2b only worked on all examples, but no perfomant on my input
	var numIterations int
	for true {
		numIterations++
		if len(currentInstructions) > 50000 {
			fmt.Println("NOT GOING TO COMPLETE")
			break
		}
		if numIterations%100 == 0 {
			log.Printf("Iteration %d. Number of instructions %d\n", numIterations, len(currentInstructions))
		}
		foundOverlap, newInstructions := findAndResolveNextPairOfOverlappingOnCubes(currentInstructions)
		currentInstructions = newInstructions
		if !foundOverlap {
			break
		}
	}
	total = 0
	for _, ins := range currentInstructions {
		//fmt.Printf("%s,v=%d\n", ins.asString(), ins.volume())
		total += ins.volume()
	}
	fmt.Printf("Final answer (hopefully!): %d\n", total)
}

func findAndResolveNextPairOfOverlappingOnCubes(instructions []*instruction) (bool, []*instruction) {
	//fmt.Printf("START LEN %d\n", len(instructions))
	additionalInstructions, toDeleteIndexOne, toDeleteIndexTwo := findNextOverlappingOnInstructionAndSplitAsNeeded(instructions)
	//fmt.Printf("Num additional %d\n", len(additionalInstructions))
	if len(additionalInstructions) == 0 { // no more overlaps
		return false, instructions
	}
	for i, ins := range instructions {
		if i == toDeleteIndexOne || i == toDeleteIndexTwo {
			continue
		}
		additionalInstructions = append(additionalInstructions, ins)
	}

	//fmt.Printf("FINISH len %d\n", len(additionalInstructions))
	return true, additionalInstructions
}

// NOTE: Order should not matter at this point
func findNextOverlappingOnInstructionAndSplitAsNeeded(instructions []*instruction) ([]*instruction, int, int) {
	for i, insI := range instructions {
		for j, insJ := range instructions {
			if i == j {
				continue
			}
			if cubesHaveOverlap(insI, insJ) {
				return resolveOverlappingOnCubes(insI, insJ), i, j
			}
		}
	}
	return []*instruction{}, -1, -1
}

func resolveOverlappingOnCubes(a, b *instruction) []*instruction {
	if a.isEqual(b) {
		return []*instruction{a}
	}

	//fmt.Printf("Comparing A %+v\n", a)
	//fmt.Printf("Comparing B %+v\n", b)

	newInstructionsA := divideOnCubeUsingOverlappingOffCube(a, b, true) // only keep the overlap for 1 of them
	//fmt.Printf("num instrucations A %d\n", len(newInstructionsA))
	newInstructionsB := divideOnCubeUsingOverlappingOffCube(b, a, false)
	//fmt.Printf("num instrucations B %d\n", len(newInstructionsB))
	result := make([]*instruction, 0)
	result = append(result, newInstructionsA...)
	result = append(result, newInstructionsB...)

	// for _, ins := range result {
	// 	fmt.Printf("Result %+v\n", ins)
	// }
	return result
}

/*
	1. Loop through all OFF instructions and convert them to ON instructions by:
	  a. Determine which of the ON cubes prior to this OFF cube, it overlaps with.
	  b. Split the ON cube it overlaps with into 27 cubes. These become 26 ON cubes, and 1 gets deleted (the one where it overlaps)
*/
func findNextOffInstructionAndApplyToPriorOns(instructions []*instruction) (bool, []*instruction) {
	found, firstOffCubeIndex := findFirstOffCubeIndex(instructions)
	if !found {
		return false, instructions
	}
	firstOffCube := instructions[firstOffCubeIndex]

	result := make([]*instruction, 0)
	for i := 0; i < firstOffCubeIndex; i++ {
		otherOnCube := instructions[i]
		// we already know that otherCube is an ON cube, due to how we are iterating
		// if it overlaps with the OFF cube, then break it up (to handle the OFF condition)
		if findSharedVolumeBetweenTwoCuboids(firstOffCube, otherOnCube) > 0 {
			newInstructions := divideOnCubeUsingOverlappingOffCube(otherOnCube, firstOffCube, false)
			result = append(result, newInstructions...)
		} else {
			result = append(result, otherOnCube) // still need to keep other ON cubes that don't overlap
		}
	}

	// append everything AFTER the OFF cube
	for i := firstOffCubeIndex + 1; i < len(instructions); i++ {
		result = append(result, instructions[i])
	}

	return true, result
}

// NOTE: this would naturally result in 27 cubes, but we don't care about:
// a) cubes which have zero volume
// b) the overlapping cube
func divideOnCubeUsingOverlappingOffCube(b, offCube *instruction) []*instruction {
	// first find the overlapping cube (guaranteed to be non-zero, due to caller function logic)
	o := findSharedCubeBetweenTwoCuboids(b, offCube)

	// NOTE: This part really reminds me of the cubes in a Rubick's cube: top layer, middle layer, bottom layer
	newInstructions := []*instruction{
		// top layer
		{isOn: true, xStart: b.xStart, xEnd: o.xStart - 1, yStart: o.yEnd + 1, yEnd: b.yEnd, zStart: o.zEnd + 1, zEnd: b.zEnd},     // #1
		{isOn: true, xStart: o.xStart, xEnd: o.xEnd, yStart: o.yEnd + 1, yEnd: b.yEnd, zStart: o.zEnd + 1, zEnd: b.zEnd},           // #2
		{isOn: true, xStart: o.xEnd + 1, xEnd: b.xEnd, yStart: o.yEnd + 1, yEnd: b.yEnd, zStart: o.zEnd + 1, zEnd: b.zEnd},         // #3
		{isOn: true, xStart: b.xStart, xEnd: o.xStart - 1, yStart: o.yStart, yEnd: o.yEnd, zStart: o.zEnd + 1, zEnd: b.zEnd},       // #4
		{isOn: true, xStart: o.xStart, xEnd: o.xEnd, yStart: o.yStart, yEnd: o.yEnd, zStart: o.zEnd + 1, zEnd: b.zEnd},             // #5
		{isOn: true, xStart: o.xEnd + 1, xEnd: b.xEnd, yStart: o.yStart, yEnd: o.yEnd, zStart: o.zEnd + 1, zEnd: b.zEnd},           // #6
		{isOn: true, xStart: b.xStart, xEnd: o.xStart - 1, yStart: b.yStart, yEnd: o.yStart - 1, zStart: o.zEnd + 1, zEnd: b.zEnd}, // #7
		{isOn: true, xStart: o.xStart, xEnd: o.xEnd, yStart: b.yStart, yEnd: o.yStart - 1, zStart: o.zEnd + 1, zEnd: b.zEnd},       // #8
		{isOn: true, xStart: o.xEnd + 1, xEnd: b.xEnd, yStart: b.yStart, yEnd: o.yStart - 1, zStart: o.zEnd + 1, zEnd: b.zEnd},     // #9

		// middle layer
		{isOn: true, xStart: b.xStart, xEnd: o.xStart - 1, yStart: o.yEnd + 1, yEnd: b.yEnd, zStart: o.zStart, zEnd: o.zEnd},     // #1
		{isOn: true, xStart: o.xStart, xEnd: o.xEnd, yStart: o.yEnd + 1, yEnd: b.yEnd, zStart: o.zStart, zEnd: o.zEnd},           // #2
		{isOn: true, xStart: o.xEnd + 1, xEnd: b.xEnd, yStart: o.yEnd + 1, yEnd: b.yEnd, zStart: o.zStart, zEnd: o.zEnd},         // #3
		{isOn: true, xStart: b.xStart, xEnd: o.xStart - 1, yStart: o.yStart, yEnd: o.yEnd, zStart: o.zStart, zEnd: o.zEnd},       // #4
		{isOn: true, xStart: o.xEnd + 1, xEnd: b.xEnd, yStart: o.yStart, yEnd: o.yEnd, zStart: o.zStart, zEnd: o.zEnd},           // #6
		{isOn: true, xStart: b.xStart, xEnd: o.xStart - 1, yStart: b.yStart, yEnd: o.yStart - 1, zStart: o.zStart, zEnd: o.zEnd}, // #7
		{isOn: true, xStart: o.xStart, xEnd: o.xEnd, yStart: b.yStart, yEnd: o.yStart - 1, zStart: o.zStart, zEnd: o.zEnd},       // #8
		{isOn: true, xStart: o.xEnd + 1, xEnd: b.xEnd, yStart: b.yStart, yEnd: o.yStart - 1, zStart: o.zStart, zEnd: o.zEnd},     // #9

		// bottom layer
		{isOn: true, xStart: b.xStart, xEnd: o.xStart - 1, yStart: o.yEnd + 1, yEnd: b.yEnd, zStart: b.zStart, zEnd: o.zStart - 1},     // #1
		{isOn: true, xStart: o.xStart, xEnd: o.xEnd, yStart: o.yEnd + 1, yEnd: b.yEnd, zStart: b.zStart, zEnd: o.zStart - 1},           // #2
		{isOn: true, xStart: o.xEnd + 1, xEnd: b.xEnd, yStart: o.yEnd + 1, yEnd: b.yEnd, zStart: b.zStart, zEnd: o.zStart - 1},         // #3
		{isOn: true, xStart: b.xStart, xEnd: o.xStart - 1, yStart: o.yStart, yEnd: o.yEnd, zStart: b.zStart, zEnd: o.zStart - 1},       // #4
		{isOn: true, xStart: o.xStart, xEnd: o.xEnd, yStart: o.yStart, yEnd: o.yEnd, zStart: b.zStart, zEnd: o.zStart - 1},             // #5
		{isOn: true, xStart: o.xEnd + 1, xEnd: b.xEnd, yStart: o.yStart, yEnd: o.yEnd, zStart: b.zStart, zEnd: o.zStart - 1},           // #6
		{isOn: true, xStart: b.xStart, xEnd: o.xStart - 1, yStart: b.yStart, yEnd: o.yStart - 1, zStart: b.zStart, zEnd: o.zStart - 1}, // #7
		{isOn: true, xStart: o.xStart, xEnd: o.xEnd, yStart: b.yStart, yEnd: o.yStart - 1, zStart: b.zStart, zEnd: o.zStart - 1},       // #8
		{isOn: true, xStart: o.xEnd + 1, xEnd: b.xEnd, yStart: b.yStart, yEnd: o.yStart - 1, zStart: b.zStart, zEnd: o.zStart - 1},     // #9
	}

	// if keepOverlap {
	// 	newInstructions = append(newInstructions, &instruction{isOn: true, xStart: o.xStart, xEnd: o.xEnd, yStart: o.yStart, yEnd: o.yEnd, zStart: o.zStart, zEnd: o.zEnd})
	// }

	// filter out instructions with volume 0
	result := make([]*instruction, 0)
	for _, ins := range newInstructions {
		if ins.volume() > 0 {
			result = append(result, ins)
		}
	}

	return result
}

func findFirstOffCubeIndex(instructions []*instruction) (bool, int) {
	for i, ins := range instructions {
		if !ins.isOn {
			return true, i
		}
	}
	return false, -1
}

// This approach is no good
func determineNumOnForSingleCuboid(i int, instructions []*instruction) int {
	thisCube := instructions[i]
	// 2. To do this, start with the first ON cube and do the following:
	if !thisCube.isOn {
		return 0
	}

	//     a. Compute the volume of the cube (as a starting point)
	numOn := thisCube.volume()
	for j := 0; j < len(instructions); j++ {
		if i == j {
			continue // don't compare against itself
		}

		otherCube := instructions[j]
		if j < i { // before this cuboid
			//     b. Cubes BEFORE this one in the list:
			if !otherCube.isOn {
				//       i.  If it is an OFF cube, ignore it.
				continue
			} else {
				//       ii. If it is an ON cube, subtract the volume where the 2 cubes overlap
				numOn -= findSharedVolumeBetweenTwoCuboids(thisCube, otherCube)
				//       Add back in volume where the other cube was negated by an OFF in between the 2 ON cubes
				numOn += adjustForNegatedCubeBetween2OnCubes(j, i, instructions)
			}
		} else { // after this cuboid
			//     c. Cubes AFTER this one in the list:
			if !otherCube.isOn {
				//       i.  If it is an OFF cube, subtract the volume that the 2 cubes overlap
				// NOTE: still need to account for the fact that could have already been subtracted by another OFF cube that overlaps the same region
				numOn -= findSharedVolumeBetweenTwoCuboids(thisCube, otherCube)
			} else {
				//       ii. If it is an ON cube, ignore it.
				continue
			}
		}
	}

	return numOn
}

// The scenario here is.... We have subtracted from the contribution of an ON cube, because a prior ON cube overlaps it.
// However, After the first one was turned ON, there could have been 1 (or more) cubes negate part of that cube.  So need
// to add that back in.
func adjustForNegatedCubeBetween2OnCubes(first, second int, instructions []*instruction) int {
	var total int
	for i := first + 1; i < second-1; i++ {
		current := instructions[i]
		if current.isOn {
			continue
		}
		sharedVolume := findSharedVolumeBetweenThreeCuboids(instructions[first], instructions[second], current)
		total += sharedVolume
	}
	return total
}

// Reference: https://stackoverflow.com/questions/5556170/finding-shared-volume-of-two-overlapping-cuboids
/*
max(min(a',x')-max(a,x),0)
* max(min(b',y')-max(b,y),0)
* max(min(c',z')-max(c,z),0)
NOTE: x' > x
NOTE: x = a.xStart //
      y = a.yStart //
	  z = a.zStart //
	  x' = a.xEnd //
	  y' = a.yEnd //
	  z' = a.zEnd //
NOTE: a = b.xStart //
      b = b.yStart //
	  c = b.zStart //
	  a' = b.xEnd //
	  b' = b.yEnd //
	  c' = b.zEnd //
*/
func findSharedVolumeBetweenTwoCuboids(a, b *instruction) int {
	shared := findSharedCubeBetweenTwoCuboids(a, b)
	return max(shared.xEnd-shared.xStart+1, 0) *
		max(shared.yEnd-shared.yStart+1, 0) *
		max(shared.zEnd-shared.zStart+1, 0)
}

func cubesHaveOverlap(a, b *instruction) bool {
	shared := findSharedCubeBetweenTwoCuboids(a, b)
	return max(shared.xEnd-shared.xStart+1, 0) > 0 &&
		max(shared.yEnd-shared.yStart+1, 0) > 0 &&
		max(shared.zEnd-shared.zStart+1, 0) > 0
}

// For more than 2, we should be able to just keep finding the intersection
func findSharedVolumeBetweenThreeCuboids(a, b, c *instruction) int {
	shared := findSharedCubeBetweenTwoCuboids(a, b)
	return findSharedVolumeBetweenTwoCuboids(c, shared)
}

func findSharedCubeBetweenTwoCuboids(a, b *instruction) *instruction {
	xStart := max(b.xStart, a.xStart)
	xEnd := min(b.xEnd, a.xEnd)
	yStart := max(b.yStart, a.yStart)
	yEnd := min(b.yEnd, a.yEnd)
	zStart := max(b.zStart, a.zStart)
	zEnd := min(b.zEnd, a.zEnd)
	return &instruction{xStart: xStart, xEnd: xEnd,
		yStart: yStart, yEnd: yEnd,
		zStart: zStart, zEnd: zEnd}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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
