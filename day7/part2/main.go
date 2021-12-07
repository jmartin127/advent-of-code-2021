package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

func main() {
	filepath := "/Users/jeff.martin@getweave.com/go/src/github.com/jmartin127/advent-of-code-2021/day7/input.txt"
	list := helpers.ReadFile(filepath)

	// convert to int
	initial := strings.Split(list[0], ",")
	listInt := make([]int, 0)
	for _, val := range initial {
		valInt, _ := strconv.Atoi(val)
		listInt = append(listInt, valInt)
	}

	// get the max value so we know how long to iterate
	var max int
	for _, val := range listInt {
		if val > max {
			max = val
		}
	}
	fmt.Printf("Max %d\n", max)

	// determine the cheapest
	minCost := -1
	for i := 0; i < max; i++ {
		sum := 0
		for _, v := range listInt {
			diff := int(math.Abs(float64(v - i)))
			sum = sum + findCost(diff)
		}
		if minCost == -1 || sum < minCost {
			minCost = sum
		}
	}

	fmt.Printf("Answer %d\n", minCost)
}

func findCost(num int) int {
	result := 0
	for i := 1; i <= num; i++ {
		result = result + i
	}
	return result
}
