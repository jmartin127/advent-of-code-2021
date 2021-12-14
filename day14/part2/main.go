package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

var inputsAlreadyCalculated map[string]map[rune]int

func main() {
	filepath := "input.txt"
	list := helpers.ReadFile(filepath)

	inputsAlreadyCalculated = make(map[string]map[rune]int, 0)

	insertions := make(map[string]rune, 0)
	for i := 2; i < len(list); i++ {
		k, v := readInsertionLine(list[i])
		insertions[k] = v
	}

	polymer := list[0]
	numIterationsToDo := 10
	countByChar := make(map[rune]int, 0)
	findCountsAfterNIterations([]rune(polymer), insertions, numIterationsToDo, countByChar)

	// add the last char
	chars := []rune(polymer)
	lastChar := chars[len(chars)-1]
	countByChar[lastChar] = countByChar[lastChar] + 1

	top, bottom := countMostLeast(countByChar)
	fmt.Printf("Top bottom %d %d\n", top, bottom)
	fmt.Printf("Answer %d\n", top-bottom)
}

func printMap(counts map[string]int) {
	for k, v := range counts {
		fmt.Printf("k=%s,v=%d\n", k, v)
	}
}

func countMostLeast(counts map[rune]int) (int, int) {
	// get the top values
	values := make([]int, 0)
	for _, v := range counts {
		values = append(values, v)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(values)))

	top := values[0]
	bottom := values[len(values)-1]
	return top, bottom
}

func findCountsAfterNIterations(input []rune, insertions map[string]rune, iterationsRemaining int, countByChar map[rune]int) {
	inputStr := string(input)
	mapKey := fmt.Sprintf("%s%d", inputStr, iterationsRemaining)
	if val, ok := inputsAlreadyCalculated[mapKey]; ok {
		addPrecomputedResultToCounts(val, countByChar)
		return
	}

	// add to the counts
	if iterationsRemaining == 0 {
		addToCounts(input, countByChar, iterationsRemaining)
		return
	}

	for i := 0; i < len(input)-1; i++ {
		first := input[i]
		second := input[i+1]
		insertion := insertions[string(first)+string(second)]

		// count first pair
		findCountsAfterNIterations([]rune{first, insertion, second}, insertions, iterationsRemaining-1, countByChar)
	}
}

func addPrecomputedResultToCounts(precomupted map[rune]int, countByChar map[rune]int) {
	for k, v := range precomupted {
		countByChar[k] = countByChar[k] + v
	}
}

func cachePrecomputed(input []rune, iterationsRemaining int) map[rune]int {
	precomupted := make(map[rune]int, 0)
	for i := 0; i < len(input)-1; i++ {
		precomupted[input[i]] = precomupted[input[i]] + 1
	}

	mapKey := fmt.Sprintf("%s%d", string(input), iterationsRemaining)
	inputsAlreadyCalculated[mapKey] = precomupted

	return precomupted
}

func addToCounts(input []rune, countByChar map[rune]int, iterationsRemaining int) {
	precomupted := cachePrecomputed(input, iterationsRemaining)
	addPrecomputedResultToCounts(precomupted, countByChar)
}

func readInsertionLine(line string) (string, rune) {
	parts := strings.Split(line, " -> ") // CH -> B
	key := parts[0]
	val := parts[1]
	return key, []rune(val)[0]
}
