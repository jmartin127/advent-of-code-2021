package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

func main() {
	filepath := "/Users/jeff.martin@getweave.com/go/src/github.com/jmartin127/advent-of-code-2021/day6/input.txt"
	list := helpers.ReadFile(filepath)

	initialFishString := strings.Split(list[0], ",")
	initialFish := make([]int, 0)
	for _, valStr := range initialFishString {
		val, _ := strconv.Atoi(valStr)
		initialFish = append(initialFish, val)
	}

	result := run(initialFish, 80)
	fmt.Printf("Answer: %d\n", len(result))
}

func run(fish []int, numIterations int) []int {
	current := fish
	for i := 0; i < numIterations; i++ {
		result := runOneIteration(current)
		current = result
	}
	return current
}

func runOneIteration(fish []int) []int {
	result := make([]int, 0)
	for _, f := range fish {
		if f == 0 {
			result = append(result, 6)
			result = append(result, 8)
		} else {
			result = append(result, f-1)
		}
	}
	return result
}
