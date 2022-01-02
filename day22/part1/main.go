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
	filepath := "day22/new-input.txt"
	list := helpers.ReadFile(filepath)

	onCubes := make(map[string]bool, 0)
	for _, line := range list {
		instruction := parseLine(line)
		processInstruction(instruction, onCubes)
		fmt.Printf("Answer %d\n", len(onCubes))
	}
	fmt.Printf("Answer %d\n", len(onCubes))
}

func processInstruction(ins *instruction, onCubes map[string]bool) {
	if !isInRange(ins) {
		return
	}

	for i := ins.xStart; i <= ins.xEnd; i++ {
		for j := ins.yStart; j <= ins.yEnd; j++ {
			for k := ins.zStart; k <= ins.zEnd; k++ {
				if ins.isOn {
					onCubes[makeKey(i, j, k)] = true
				} else {
					delete(onCubes, makeKey(i, j, k))
				}
			}
		}
	}
}

// The initialization procedure only uses cubes that have x, y, and z positions of at least -50 and at most 50. For now, ignore cubes outside this region.
func isInRange(ins *instruction) bool {
	return isPointInRange(ins.xStart, ins.yStart, ins.zStart) && isPointInRange(ins.xEnd, ins.yEnd, ins.zEnd)
}

func isPointInRange(x, y, z int) bool {
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
