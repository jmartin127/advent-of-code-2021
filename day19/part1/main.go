package main

import (
	"fmt"
	"strconv"
	"strings"
)

type scanner struct {
	id              int
	beaconPositions [][]int // x,y,z
}

type operation struct {
	sourcePos int
	negative  bool
}

type permutation struct {
	pos0 *operation
	pos1 *operation
	pos2 *operation
}

// TODO need 24 of these, but which ones?
var permutations = []*permutation{
	{pos0: &operation{sourcePos: 0, negative: false}, pos1: &operation{sourcePos: 1, negative: false}, pos2: &operation{sourcePos: 2, negative: false}}, // x,y,z
	{pos0: &operation{sourcePos: 0, negative: true}, pos1: &operation{sourcePos: 2, negative: true}, pos2: &operation{sourcePos: 1, negative: true}},    // -x,-z,-y
	{pos0: &operation{sourcePos: 2, negative: true}, pos1: &operation{sourcePos: 1, negative: false}, pos2: &operation{sourcePos: 0, negative: false}},  // -z,y,x
	{pos0: &operation{sourcePos: 2, negative: false}, pos1: &operation{sourcePos: 1, negative: true}, pos2: &operation{sourcePos: 0, negative: false}},  // z,-y,x
	{pos0: &operation{sourcePos: 1, negative: true}, pos1: &operation{sourcePos: 2, negative: false}, pos2: &operation{sourcePos: 0, negative: true}},   // -y,z,-x
}

func (s *scanner) print() {
	fmt.Printf("--- scanner %d ---\n", s.id)
	for _, bp := range s.beaconPositions {
		fmt.Printf("%v\n", bp)
	}
	fmt.Println()
}

func main() {
	test([]string{"x", "y", "z"})

	// list := helpers.ReadFile("input.txt")
	// scanners := parseScanners(list)
	// scanners[0].rotateScanner()
}

/*
Options:
1) just take the diff
2) multiple by -1, then take diff
*/
func (s *scanner) scannersMatch(other *scanner) {
	diffMap := make(map[string]int, 0)

	// do a pair-wise comparison of every beacon, just taking the diffs
	for _, bp := range s.beaconPositions {
		for _, otherBp := range other.beaconPositions {
			diff := make([]int, 0)
			for pos := 0; pos < 3; pos++ {
				diff = append(diff, otherBp[pos]-bp[pos])
			}
			key := convertIntArrayToString(diff)
			if _, ok := diffMap[key]; ok {
				diffMap[key] = diffMap[key] + 1
			} else {
				diffMap[key] = 1
			}
		}
	}

	// check if any of these occur more than

}

// By finding pairs of scanners that both see at least 12 of the same beacons, you can assemble the entire map.
func findOverlap(input map[string]int) {
	for k, v := range input {
		if v >= 12 {
			fmt.Printf("Found %s\n", k)
		}
	}
}

func convertIntArrayToString(input []int) string {
	return strconv.Itoa(input[0]) + "," + strconv.Itoa(input[1]) + "," + strconv.Itoa(input[2])
}

func (s *scanner) rotateScanner() {
	for _, p := range permutations {
		newScanner := s.applyPermutation(p)
		newScanner.print()
	}
}

func (s *scanner) applyPermutation(p *permutation) *scanner {
	new := s.copy()

	for i := range s.beaconPositions {
		oldBp := s.beaconPositions[i]
		newBp := new.beaconPositions[i]

		newBp[0] = oldBp[p.pos0.sourcePos]
		if p.pos0.negative {
			newBp[0] = newBp[0] * -1
		}

		newBp[1] = oldBp[p.pos1.sourcePos]
		if p.pos1.negative {
			newBp[1] = newBp[1] * -1
		}

		newBp[2] = oldBp[p.pos2.sourcePos]
		if p.pos2.negative {
			newBp[2] = newBp[2] * -1
		}
	}

	return new
}

func (s *scanner) copy() *scanner {
	copy := &scanner{
		id:              s.id,
		beaconPositions: make([][]int, 0),
	}
	for _, bp := range s.beaconPositions {
		copy.beaconPositions = append(copy.beaconPositions, copyPosition(bp))
	}
	return copy
}

func copyPosition(input []int) []int {
	result := make([]int, 0)
	for _, v := range input {
		result = append(result, v)
	}
	return result
}

/*
--- scanner 0 ---
-1,-1,1
-2,-2,2
-3,-3,3
-2,-3,1
5,6,-4
8,0,7
*/
func parseScanners(list []string) []*scanner {
	scanners := make([]*scanner, 0)
	var currScanner *scanner
	for _, line := range list {
		if strings.Contains(line, "scanner") {
			if currScanner != nil {
				scanners = append(scanners, currScanner)
			}
			parts := strings.Split(line, " ") // --- scanner 0 ---
			id, _ := strconv.Atoi(parts[2])
			currScanner = &scanner{
				id:              id,
				beaconPositions: make([][]int, 0),
			}
		} else if strings.Contains(line, ",") {
			currScanner.beaconPositions = append(currScanner.beaconPositions, parseInts(line))
		} else {
			// skip empty lines
		}
	}
	if currScanner != nil {
		scanners = append(scanners, currScanner)
	}
	return scanners
}

// -1,-1,1
func parseInts(line string) []int {
	vals := strings.Split(line, ",")
	x, _ := strconv.Atoi(vals[0])
	y, _ := strconv.Atoi(vals[1])
	z, _ := strconv.Atoi(vals[2])
	return []int{x, y, z}
}

// https://stackoverflow.com/questions/16452383/how-to-get-all-24-rotations-of-a-3-dimensional-array
func test(v []string) {
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			v = roll(v)
			fmt.Printf("%+v\n", v)
			for k := 0; k < 3; k++ {
				v = turn(v)
				fmt.Printf("%+v\n", v)
			}
		}
		v = roll(turn(roll(v)))
	}
}

var flipMapping = map[string]string{
	"x":  "-x",
	"y":  "-y",
	"z":  "-z",
	"-x": "x",
	"-y": "y",
	"-z": "z",
}

func roll(input []string) []string {
	return []string{input[0], input[2], flipMapping[input[1]]}
}

func turn(input []string) []string {
	return []string{flipMapping[input[1]], input[0], input[2]}
}
