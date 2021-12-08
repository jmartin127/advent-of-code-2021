package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

/*
  0:      1:      2:      3:      4:
 aaaa    ....    aaaa    aaaa    ....
b    c  .    c  .    c  .    c  b    c
b    c  .    c  .    c  .    c  b    c
 ....    ....    dddd    dddd    dddd
e    f  .    f  e    .  .    f  .    f
e    f  .    f  e    .  .    f  .    f
 gggg    ....    gggg    gggg    ....

  5:      6:      7:      8:      9:
 aaaa    aaaa    aaaa    aaaa    aaaa
b    .  b    .  .    c  b    c  b    c
b    .  b    .  .    c  b    c  b    c
 dddd    dddd    ....    dddd    dddd
.    f  e    f  .    f  e    f  .    f
.    f  e    f  .    f  e    f  .    f
 gggg    gggg    ....    gggg    gggg

*/
var NUMBER_DEFS = map[int][]int{
	0: {TOP, TOP_LEFT, TOP_RIGHT, BOTTOM_LEFT, BOTTOM_RIGHT, BOTTOM},
	1: {TOP_RIGHT, BOTTOM_RIGHT},
	2: {TOP, TOP_RIGHT, MIDDLE, BOTTOM_LEFT, BOTTOM},
	3: {TOP, TOP_RIGHT, MIDDLE, BOTTOM_RIGHT, BOTTOM},
	4: {TOP_LEFT, TOP_RIGHT, MIDDLE, BOTTOM_RIGHT},
	5: {TOP, TOP_LEFT, MIDDLE, BOTTOM_RIGHT, BOTTOM},
	6: {TOP, TOP_LEFT, MIDDLE, BOTTOM_LEFT, BOTTOM_RIGHT, BOTTOM},
	7: {TOP, TOP_RIGHT, BOTTOM_RIGHT},
	8: {TOP, TOP_LEFT, TOP_RIGHT, MIDDLE, BOTTOM_LEFT, BOTTOM_RIGHT, BOTTOM},
	9: {TOP, TOP_LEFT, TOP_RIGHT, MIDDLE, BOTTOM_RIGHT, BOTTOM},
}

const (
	TOP          = 0
	TOP_LEFT     = 1
	TOP_RIGHT    = 2
	MIDDLE       = 3
	BOTTOM_LEFT  = 4
	BOTTOM_RIGHT = 5
	BOTTOM       = 6
)

// Possible numbers by string length
var POSS_NUMS_BY_LEN = map[int][]int{
	2: {1},
	4: {4},
	3: {7},
	7: {8},
	5: {2, 3, 5},
	6: {0, 6, 9},
}

type possibleValues struct {
	vals []string
}

type board struct {
	possibilities map[int]*possibleValues
}

type partialNumber struct {
	mustBePresent map[int]bool
	notFilledIn   map[int]bool
}

func (b *board) print() {
	for k, v := range b.possibilities {
		fmt.Printf("key=%d, val=%v\n", k, v)
	}
}

func (b *board) invert() map[string]int {
	result := make(map[string]int, 0)
	for k, v := range b.possibilities {
		result[v.vals[0]] = k
	}
	return result
}

func main() {
	list := helpers.ReadFile("input.txt")

	result := 0
	for _, line := range list {
		result += deduceSum(line)
	}
	fmt.Printf("Answer %d\n", result)
}

// fgaebd cg bdaec gdafb agbcfd gdcbef bgcad gfac gcb cdgabef | cg cg fdcagb cbg
func deduceSum(line string) int {
	parts := strings.Split(line, " | ")
	inputParts := strings.Split(parts[0], " ")
	outputParts := strings.Split(parts[1], " ")

	b := determineOriginalInputMapping(inputParts)

	// determine what each input value maps to, to narrow down the deduction
	useInputToFinalizeMapping(b, inputParts)
	panicIfNotDeduced(b)

	// now that we know where each letter goes, map each one to a number
	posByCharacter := b.invert()
	resultString := ""
	for _, outputChunk := range outputParts {
		num := determineNumber(posByCharacter, outputChunk)
		resultString = resultString + strconv.Itoa(num)
	}

	result, _ := strconv.Atoi(resultString)

	return result
}

func determineNumber(posByCharacter map[string]int, outputChunk string) int {
	numberTemplate := make([]int, 0)
	for _, c := range strings.Split(outputChunk, "") {
		pos := posByCharacter[c]
		numberTemplate = append(numberTemplate, pos)
	}

	for num, numDef := range NUMBER_DEFS {
		if slicesMatch(numberTemplate, numDef) {
			return num
		}
	}

	return -1
}

func slicesMatch(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	sort.Ints(a)
	sort.Ints(b)

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func panicIfNotDeduced(b *board) {
	for k, v := range b.possibilities {
		if len(v.vals) > 1 {
			panic(fmt.Sprintf("not reduced for %d\n", k))
		}
	}
}

// Iterate over each chunk and apply the deductive reasoning
// example: fgaebd
func useInputToFinalizeMapping(b *board, inputChunks []string) {
	for _, inputChunk := range inputChunks {
		lenOfChunk := len(inputChunk)
		possibleNums := POSS_NUMS_BY_LEN[lenOfChunk]
		if len(possibleNums) == 1 {
			continue // already solved
		}

		useSingleInputToUpdateMapping(b, inputChunk, possibleNums)
	}

	// check if there are still unresolved parts of the board. resolves these cases:
	/*
		key=1, val=&{[f]}
		key=3, val=&{[a]}
		key=6, val=&{[d]}
		key=4, val=&{[d e]}
		key=2, val=&{[c]}
		key=5, val=&{[g]}
		key=0, val=&{[b]}
	*/
	for pos, vals := range b.possibilities {
		if len(vals.vals) > 1 {
			applyLastReduction(b, pos)
		}
	}
}

func applyLastReduction(b *board, pos int) {
	vals := b.possibilities[pos].vals
	for _, pos2 := range b.possibilities {
		if len(pos2.vals) == 1 {
			if sliceContains(vals, pos2.vals[0]) {
				b.possibilities[pos].vals = findSecondNotInFirst(pos2.vals, vals)
				return
			}
		}
	}
}

// Example:
// inputChunk:   gdafb
// possibleNums: 2,3,4
func useSingleInputToUpdateMapping(b *board, inputChunk string, possibleNums []int) {
	pn := &partialNumber{
		mustBePresent: make(map[int]bool, 0),
	}

	// check each position in the number template to see if the required chars are present
	for i := 0; i <= 6; i++ { // iterate over all positions in the number template
		poss := b.possibilities[i].vals
		if chunkContainsAllValues(inputChunk, poss) {
			pn.mustBePresent[i] = true
		}
	}

	// check if a single numbers satifies the required fields
	possibleMatches := make([]int, 0)
	for _, possibleNum := range possibleNums {
		if possibleNumMatchesTemplate(possibleNum, pn) {
			possibleMatches = append(possibleMatches, possibleNum)
		}
	}

	// check if we were able to narrow down to a single number
	if len(possibleMatches) != 1 {
		return
	}

	// if we were able to narrow it down, then do further deduction of outstanding characters
	matchedNum := possibleMatches[0]

	compareTemplateToMatch(inputChunk, matchedNum, pn, b)
}

func compareTemplateToMatch(inputChunk string, matchedNum int, pn *partialNumber, b *board) {
	// reduce chunk to only those that are not required
	charsToRemove := []string{}
	for mustBePresentPos := range pn.mustBePresent {
		poss := b.possibilities[mustBePresentPos].vals
		for _, p := range poss {
			charsToRemove = append(charsToRemove, p)
		}
	}
	chunkParts := strings.Split(inputChunk, "")
	leftOverFromChunk := findSecondNotInFirst(charsToRemove, chunkParts)

	// update template with fields that must be filled in for the matched num
	numDef := NUMBER_DEFS[matchedNum]
	for _, pos := range numDef {
		if _, ok := pn.mustBePresent[pos]; ok {
			continue
		}

		poss := b.possibilities[pos]

		// find intersection of possibilities + those left over in chunk
		intersection := findIntersection(poss.vals, leftOverFromChunk)

		// reduce the induction space to just these values
		b.possibilities[pos].vals = intersection
	}
}

func findIntersection(a, b []string) []string {
	result := make([]string, 0)
	for _, va := range a {
		if sliceContains(b, va) {
			result = append(result, va)
		}
	}
	return result
}

// Example:
// possibleNum: 0
// pn: map of required pieces of the number (top, top-left, etc.)
func possibleNumMatchesTemplate(possibleNum int, pn *partialNumber) bool {
	numberDef := NUMBER_DEFS[possibleNum]

	// loop through required positions in partial number
	for posMustBePresent := range pn.mustBePresent {
		if !sliceContainsInt(numberDef, posMustBePresent) {
			return false
		}
	}
	return true
}

// Example:
// inputChunk: gdafb
// required: b (for example, top part of the 7)
func chunkContainsAllValues(inputChunk string, required []string) bool {
	inputVals := strings.Split(inputChunk, "")
	for _, req := range required {
		if !sliceContains(inputVals, req) {
			return false
		}
	}
	return true
}

func determineOriginalInputMapping(inputChunks []string) *board {
	// find unique chunks
	mapping := make(map[int]string, 0)
	for _, inputChunk := range inputChunks {
		if len(inputChunk) == 2 {
			mapping[1] = inputChunk
		} else if len(inputChunk) == 4 {
			mapping[4] = inputChunk
		} else if len(inputChunk) == 3 {
			mapping[7] = inputChunk
		} else if len(inputChunk) == 7 {
			mapping[8] = inputChunk
		}
	}

	// create the board
	b := board{
		possibilities: make(map[int]*possibleValues, 0),
	}

	// map unique chunks onto the board
	// map the 1
	oneVals := strings.Split(mapping[1], "")
	b.possibilities[TOP_RIGHT] = &possibleValues{vals: oneVals}
	b.possibilities[BOTTOM_RIGHT] = &possibleValues{vals: oneVals}

	// find the 1 value in the 7 that is unique, and set it
	sevenVals := strings.Split(mapping[7], "")
	uniqueIn7 := findSecondNotInFirst(oneVals, sevenVals)
	if len(uniqueIn7) != 1 {
		panic("go back to the drawing board for 7s")
	}
	b.possibilities[TOP] = &possibleValues{vals: uniqueIn7}

	// find the 2 values in the 4 that are unique, and set them
	fourVals := strings.Split(mapping[4], "")
	uniqueIn4 := findSecondNotInFirst(oneVals, fourVals)
	if len(uniqueIn4) != 2 {
		panic("go back to the drawing board for 4s")
	}
	b.possibilities[TOP_LEFT] = &possibleValues{vals: uniqueIn4}
	b.possibilities[MIDDLE] = &possibleValues{vals: uniqueIn4}

	// set the possible values for the bottom and bottom left, from the 8
	eightVals := strings.Split(mapping[8], "")
	uniqueIn8 := findSecondNotInFirst(fourVals, eightVals)
	uniqueIn8 = findSecondNotInFirst(sevenVals, uniqueIn8)
	uniqueIn8 = findSecondNotInFirst(oneVals, uniqueIn8)
	if len(uniqueIn4) != 2 {
		panic("go back to the drawing board for 8s")
	}
	b.possibilities[BOTTOM] = &possibleValues{vals: uniqueIn8}
	b.possibilities[BOTTOM_LEFT] = &possibleValues{vals: uniqueIn8}

	return &b
}

func findSecondNotInFirst(first []string, second []string) []string {
	result := make([]string, 0)
	for _, sec := range second {
		if !sliceContains(first, sec) {
			result = append(result, sec)
		}
	}
	return result
}

func sliceContains(ref []string, val string) bool {
	for _, r := range ref {
		if r == val {
			return true
		}
	}

	return false
}

func sliceContainsInt(ref []int, val int) bool {
	for _, r := range ref {
		if r == val {
			return true
		}
	}

	return false
}
