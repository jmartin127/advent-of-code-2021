package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

func main() {
	// determine number of initial iterations to run
	totalIterations := 256
	numInitial := totalIterations % 7
	fmt.Printf("Num initial iterations %d\n", numInitial)

	// run initial iterations
	fishAfterInitial := runInitialIterations(numInitial)

	fishByCount := countFish(fishAfterInitial)

	currentResult := fishByCount
	for i := 0; i < 36; i++ {
		result := run7Iterations(currentResult)
		currentResult = result
	}
	fmt.Printf("Total %d\n", totalInMap(currentResult))
}

func runInitialIterations(numInitial int) []int {
	filepath := "/Users/jeff.martin@getweave.com/go/src/github.com/jmartin127/advent-of-code-2021/day6/input.txt"
	list := helpers.ReadFile(filepath)

	initialFishString := strings.Split(list[0], ",")
	initialFish := make([]int, 0)
	for _, valStr := range initialFishString {
		val, _ := strconv.Atoi(valStr)
		initialFish = append(initialFish, val)
	}

	return run(initialFish, numInitial)
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

func totalInMap(fish map[int]int) int {
	sum := 0
	for _, v := range fish {
		sum += v
	}
	return sum
}

// 0 ==> 0 + 2
// 1 ==> 1 + 3
// 2 ==> 2 + 4
// 3 ==> 3 + 5
// 4 ==> 4 + 6
// 5 ==> 5 + 7
// 6 ==> 6 + 8

// 7 ==> 0
// 8 ==> 1
func run7Iterations(fish map[int]int) map[int]int {
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
		if _, ok := initialFish[val]; ok {
			initialFish[val] = initialFish[val] + 1
		} else {
			initialFish[val] = 1
		}
	}
	return initialFish
}
