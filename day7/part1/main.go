package main

import (
	"fmt"
	"math"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

func main() {
	intList := helpers.ReadSingleLineFileAsInts("input.txt")

	// get the max value so we know how long to iterate
	max := helpers.MaxFromList(intList)
	fmt.Printf("Max %d\n", max)

	// determine the cheapest
	minCost := 10000000
	for i := 0; i < max; i++ {
		var sum int
		for _, v := range intList {
			diff := int(math.Abs(float64(v - i)))
			sum += diff
		}
		if sum < minCost {
			minCost = sum
		}
	}

	fmt.Printf("Answer %d\n", minCost)
}
