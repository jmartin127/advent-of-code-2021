package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

type scanner struct {
	id              int
	beaconPositions [][]int // x,y,z
}

type orientedScanner struct {
	s        *scanner
	position []int // position relative to scanner 0
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

var permutations = []*permutation{}

// d = just take difference
// m = multiply by -1, then take difference
var combinationsOfMatchTypes = [][]string{
	{"d", "d", "d"},
	{"d", "d", "m"},
	{"d", "m", "d"},
	{"d", "m", "m"},
	{"m", "d", "d"},
	{"m", "d", "m"},
	{"m", "m", "d"},
	{"m", "m", "m"},
}

func (s *scanner) print() {
	fmt.Printf("--- scanner %d ---\n", s.id)
	for _, bp := range s.beaconPositions {
		fmt.Printf("%v\n", bp)
	}
	fmt.Println()
}

func main() {
	setPermutations()
	list := helpers.ReadFile("input.txt")
	scanners := parseScanners(list)

	// remove the first scanner from the map
	firstScanner := scanners[0]
	delete(scanners, 0)

	lockedScanners, _ := lockScannerPositions(map[int]*orientedScanner{0: &orientedScanner{s: firstScanner, position: []int{0, 0, 0}}}, scanners)
	fmt.Printf("Locked scanners %d\n", len(lockedScanners))

	// Next steps:
	// 3. Once all scanners are locked in (relative to 0), loop through and obtain increment common beacon map (relative to 0)
}

// 1. Once a match is found with a scanner relative to a previously locked in scanner, LOCK that position in place (map of scanner position --> *scanner (rotated relative to zero))
func lockScannerPositions(lockedScanners map[int]*orientedScanner, remainingScanners map[int]*scanner) (map[int]*orientedScanner, map[int]*scanner) {
	if len(lockedScanners) == 0 {
		return lockedScanners, remainingScanners
	}

	// 2. Then compare to locked in positions and continue to lock in until all scanners are locked in.
	resultLocked := copyOrientedMap(lockedScanners)
	resultRemaining := copyMap(remainingScanners)
	for rsPos, rs := range remainingScanners {
		for lsPos, ls := range lockedScanners {
			if hasOverlap, os := ls.s.scannersMatch(rs); hasOverlap {
				positionRelativeToZero := addByPos(ls.position, os.position)
				resultLocked[rsPos] = &orientedScanner{s: os.s, position: positionRelativeToZero}
				fmt.Printf("Adding %+v and %+v\n", ls.position, os.position)
				fmt.Printf("Locking in scanner %d at position %+v, relative to scanner at position %d\n", rsPos, positionRelativeToZero, lsPos)
				delete(resultRemaining, rsPos)
				return lockScannerPositions(resultLocked, resultRemaining)
			}
		}
	}

	return lockedScanners, remainingScanners
}

func addByPos(a []int, b []int) []int {
	result := make([]int, 0)
	for i := 0; i < len(a); i++ {
		result = append(result, a[i]+b[i])
	}
	return result
}

func copyMap(a map[int]*scanner) map[int]*scanner {
	result := make(map[int]*scanner, 0)
	for k, v := range a {
		result[k] = v
	}
	return result
}

func copyOrientedMap(a map[int]*orientedScanner) map[int]*orientedScanner {
	result := make(map[int]*orientedScanner, 0)
	for k, v := range a {
		result[k] = v
	}
	return result
}

func setPermutations() {
	rotations := generate24PossibleRotations([]string{"x", "y", "z"})
	perms := make([]*permutation, 0)
	for _, rotation := range rotations {
		perm := convertRotationToPermutation(rotation)
		perms = append(perms, perm)
	}
	permutations = perms
}

func (s *scanner) scannersMatch(other *scanner) (bool, *orientedScanner) {
	for _, perm := range permutations {
		//fmt.Printf("perm %d\n", i)
		other = other.applyPermutation(perm)
		if hasOverlap, relativePos := s.scannersMatchAtCurrentPosition(other); hasOverlap {
			return hasOverlap, &orientedScanner{s: other, position: relativePos}
		}
	}

	return false, &orientedScanner{s: other, position: []int{}}
}

/*
Options:
1) just take the diff
2) multiple by -1, then take diff
*/
func (s *scanner) scannersMatchAtCurrentPosition(other *scanner) (bool, []int) {
	// try each combination at each position (e.g., d,d,m)
	for _, matchType := range combinationsOfMatchTypes {
		// do a pair-wise comparison of every beacon, just taking the diffs
		diffMap := make(map[string]int, 0)
		for _, bp := range s.beaconPositions {
			for _, otherBp := range other.beaconPositions {
				diff := make([]int, 0)
				for pos := 0; pos < 3; pos++ {
					diffType := matchType[pos]
					if diffType == "d" {
						diff = append(diff, otherBp[pos]-bp[pos])
					} else { // m
						diff = append(diff, (otherBp[pos]*-1)-bp[pos])
					}
				}
				incrementMap(diffMap, convertIntArrayToString(diff))
			}
		}
		// check if any of these occur more than N number of times
		if hasOverlap, relativePos := findOverlap(diffMap); hasOverlap {
			//fmt.Printf("1st approach match %+v\n", relativePos) // TODO
			return hasOverlap, relativePos
		}
	}
	return false, []int{}
}

func incrementMap(diffMap map[string]int, key string) {
	if _, ok := diffMap[key]; ok {
		diffMap[key] = diffMap[key] + 1
	} else {
		diffMap[key] = 1
	}
}

// By finding pairs of scanners that both see at least 12 of the same beacons, you can assemble the entire map.
func findOverlap(input map[string]int) (bool, []int) {
	for k, v := range input {
		if v >= 12 {
			return true, findRelativePosition(convertStringToIntArray(k)) // TODO perhaps in the -1 case, may NOT want to invert the value
		}
	}
	return false, []int{}
}

func findRelativePosition(input []int) []int {
	result := make([]int, 0)
	for _, v := range input {
		result = append(result, v*-1)
	}

	return result
}

func convertIntArrayToString(input []int) string {
	return strconv.Itoa(input[0]) + "," + strconv.Itoa(input[1]) + "," + strconv.Itoa(input[2])
}

func convertStringToIntArray(input string) []int {
	vals := strings.Split(input, ",")
	result := make([]int, 0)
	for _, v := range vals {
		i, _ := strconv.Atoi(v)
		result = append(result, i)
	}
	return result
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
func parseScanners(list []string) map[int]*scanner {
	scanners := make(map[int]*scanner, 0)
	var currScanner *scanner
	var currentCount int
	for _, line := range list {
		if strings.Contains(line, "scanner") {
			if currScanner != nil {
				scanners[currentCount] = *&currScanner
				currentCount++
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
		scanners[currentCount] = *&currScanner
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

// convert -x,-z,-y --> pos0: &operation{sourcePos: 0, negative: true}, pos1: &operation{sourcePos: 2, negative: true}, pos2: &operation{sourcePos: 1, negative: true}
func convertRotationToPermutation(rotation []string) *permutation {
	pos0op := convertLetterToOperation(rotation[0])
	pos1op := convertLetterToOperation(rotation[1])
	pos2op := convertLetterToOperation(rotation[2])

	return &permutation{pos0: pos0op, pos1: pos1op, pos2: pos2op}
}

func convertLetterToOperation(letter string) *operation {
	switch letter {
	case "-x":
		return &operation{sourcePos: 0, negative: true}
	case "-y":
		return &operation{sourcePos: 1, negative: true}
	case "-z":
		return &operation{sourcePos: 2, negative: true}
	case "x":
		return &operation{sourcePos: 0, negative: false}
	case "y":
		return &operation{sourcePos: 1, negative: false}
	case "z":
		return &operation{sourcePos: 2, negative: false}
	default:
		log.Fatal("oops")
		return nil
	}
}

// https://stackoverflow.com/questions/16452383/how-to-get-all-24-rotations-of-a-3-dimensional-array
func generate24PossibleRotations(v []string) [][]string {
	result := make([][]string, 0)
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			v = roll(v)
			result = append(result, v)
			for k := 0; k < 3; k++ {
				v = turn(v)
				result = append(result, v)
			}
		}
		v = roll(turn(roll(v)))
	}
	return result
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
