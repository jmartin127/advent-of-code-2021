package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

func main() {
	filepath := "input.txt"
	list := helpers.ReadFile(filepath)

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
	// add to the counts
	if iterationsRemaining == 0 {
		addToCounts(input, countByChar)
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

func addToCounts(input []rune, countByChar map[rune]int) {
	numToCount := len(input) - 1
	for i := 0; i < numToCount; i++ {
		char := input[i]
		countByChar[char] = countByChar[char] + 1
	}
}

func readInsertionLine(line string) (string, rune) {
	parts := strings.Split(line, " -> ") // CH -> B
	key := parts[0]
	val := parts[1]
	return key, []rune(val)[0]
}
