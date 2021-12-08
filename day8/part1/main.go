package main

import (
	"fmt"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

func main() {
	list := helpers.ReadFile("input.txt")

	result := 0
	for _, line := range list {
		parts := strings.Split(line, " | ")
		chunks := strings.Split(parts[1], " ")
		result += numFromChunks(chunks)
	}
	fmt.Printf("Answer %d\n", result)
}

func numFromChunks(chunks []string) int {
	result := 0
	for _, c := range chunks {
		if len(c) == 2 || len(c) == 4 || len(c) == 3 || len(c) == 7 {
			result++
		}
	}
	return result
}
