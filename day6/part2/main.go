package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

const (
	ITERATIONS     = 256
	CYCLE_DURATION = 7
)

func main() {
	// determine number of initial iterations to run, and how many subsequent cycles
	totalIterations := ITERATIONS
	numInitial := totalIterations % CYCLE_DURATION
	numCycles := (totalIterations - numInitial) / CYCLE_DURATION
	fmt.Printf("Num initial iterations: %d\n", numInitial)
	fmt.Printf("Num cycles after initial: %d\n", numCycles)

	// run initial iterations
	initialFish := parseInput()
	fishAfterInitial := runInitialIterations(initialFish, numInitial)

	// count num of occurences of each fish
	fishByCount := countFish(fishAfterInitial)

	// run remaining iterations
	currentResult := fishByCount
	for i := 0; i < numCycles; i++ {
		result := runFullCycle(currentResult)
		currentResult = result
	}
	fmt.Printf("Answer %d\n", totalInMap(currentResult))
}

// PART 1 //

func parseInput() []int {
	filepath := "/Users/jeff.martin@getweave.com/go/src/github.com/jmartin127/advent-of-code-2021/day6/input.txt"
	list := helpers.ReadFile(filepath)

	initialFishString := strings.Split(list[0], ",")
	initialFish := make([]int, 0)
	for _, valStr := range initialFishString {
		val, _ := strconv.Atoi(valStr)
		initialFish = append(initialFish, val)
	}

	return initialFish
}

func runInitialIterations(fish []int, numIterations int) []int {
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

// PART 2 //

func totalInMap(fish map[int]int) int {
	sum := 0
	for _, v := range fish {
		sum += v
	}
	return sum
}

// Execute a full cycle (7 iterations), by determining how each fish will reproduce in the next iteration
// 0 ==> 0 + 2 (prouces another 0 fish, and also a 2 fish after 1 full cycle)
// 1 ==> 1 + 3 (etc.)
// 2 ==> 2 + 4
// 3 ==> 3 + 5
// 4 ==> 4 + 6
// 5 ==> 5 + 7
// 6 ==> 6 + 8

// 7 ==> 0
// 8 ==> 1
func runFullCycle(fish map[int]int) map[int]int {
	result := make(map[int]int, 0)

	for k, v := range fish {
		if k <= 6 {
			addToMap(result, k, v)
			addToMap(result, k+2, v)
		} else {
			addToMap(result, k-7, v)
		}
	}

	return result
}

func addToMap(fish map[int]int, val int, num int) {
	if _, ok := fish[val]; ok {
		fish[val] = fish[val] + num
	} else {
		fish[val] = num
	}
}

func countFish(fish []int) map[int]int {
	initialFish := make(map[int]int, 0)
	for _, val := range fish {
		addToMap(initialFish, val, 1)
	}
	return initialFish
}
