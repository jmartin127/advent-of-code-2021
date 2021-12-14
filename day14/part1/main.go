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

	insertions := make(map[string]string, 0)
	for i := 2; i < len(list); i++ {
		k, v := readInsertionLine(list[i])
		insertions[k] = v
	}

	polymer := list[0]
	for i := 0; i < 10; i++ {
		polymer = runOneIteration(polymer, insertions)
	}
	fmt.Printf("Result %s\n", polymer)

	top, bottom := countMostLeast(polymer)
	fmt.Printf("Top bottom %d %d\n", top, bottom)
	fmt.Printf("Asnwer %d\n", top-bottom)
}

func countMostLeast(input string) (int, int) {
	counts := make(map[string]int, 0)
	for _, c := range strings.Split(input, "") {
		counts[c] = counts[c] + 1
	}

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

func runOneIteration(input string, insertions map[string]string) string {
	polymer := input
	polymerRunes := []rune(polymer)
	result := ""
	for i := 0; i < len(polymerRunes)-1; i++ {
		first := polymerRunes[i]
		second := polymerRunes[i+1]
		together := fmt.Sprintf("%s%s", string(first), string(second))
		//fmt.Printf("look up %s\n", together)
		insertion := insertions[together]
		result += string(first)
		result += string(insertion)
	}
	result += string(polymerRunes[len(polymerRunes)-1])
	return result
}

func readInsertionLine(line string) (string, string) {
	parts := strings.Split(line, " -> ") // CH -> B
	key := parts[0]
	val := parts[1]
	return key, val
}
