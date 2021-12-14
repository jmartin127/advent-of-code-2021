package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

const NUM_ITERATIONS = 40

var inputsAlreadyCalculated = make(map[string]map[rune]int, 0)

func main() {
	// parse the input
	list := helpers.ReadFile("input.txt")
	insertions := make(map[string]rune, 0)
	for i := 2; i < len(list); i++ {
		k, v := readInsertionLine(list[i])
		insertions[k] = v
	}

	// iterate
	polymer := list[0]
	result := findCountsAfterNIterations([]rune(polymer), insertions, NUM_ITERATIONS)

	// add a count for the last character
	chars := []rune(polymer)
	lastChar := chars[len(chars)-1]
	result[lastChar] = result[lastChar] + 1

	// compute the answer
	top, bottom := countMostLeast(result)
	fmt.Printf("Top bottom %d %d\n", top, bottom)
	fmt.Printf("Answer %d\n", top-bottom)
}

func countMostLeast(counts map[rune]int) (int, int) {
	values := make([]int, 0)
	for _, v := range counts {
		values = append(values, v)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(values)))
	return values[0], values[len(values)-1]
}

// main premise is that we can recurively determine the character counts for each subset of the problem
// also, since the numbers are so huge, don't recompute something we've already computed
func findCountsAfterNIterations(input []rune, insertions map[string]rune, iterationsRemaining int) map[rune]int {
	// see if we already know the answer
	mapKey := fmt.Sprintf("%s%d", string(input), iterationsRemaining)
	if precomupted, ok := inputsAlreadyCalculated[mapKey]; ok {
		return precomupted
	}

	// add to the counts
	if iterationsRemaining == 0 {
		return computeCounts(input)
	}

	// since we don't know the answer, compute it recursively
	total := make(map[rune]int, 0)
	for i := 0; i < len(input)-1; i++ {
		first := input[i]
		second := input[i+1]
		insertion := insertions[string(first)+string(second)]

		result := findCountsAfterNIterations([]rune{first, insertion, second}, insertions, iterationsRemaining-1)
		addToTotal(total, result)
	}

	inputsAlreadyCalculated[mapKey] = total
	return total
}

func addToTotal(total map[rune]int, toAdd map[rune]int) {
	for k, v := range toAdd {
		total[k] = total[k] + v
	}
}

func computeCounts(input []rune) map[rune]int {
	precomupted := make(map[rune]int, 0)
	for i := 0; i < len(input)-1; i++ {
		precomupted[input[i]] = precomupted[input[i]] + 1
	}
	mapKey := fmt.Sprintf("%s%d", string(input), 0) // only on zeroth iteration
	inputsAlreadyCalculated[mapKey] = precomupted
	return precomupted
}

func readInsertionLine(line string) (string, rune) {
	parts := strings.Split(line, " -> ") // CH -> B
	return parts[0], []rune(parts[1])[0]
}
