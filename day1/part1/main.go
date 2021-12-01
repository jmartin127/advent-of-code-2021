package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("/Users/jeff.martin@getweave.com/go/src/github.com/jmartin127/advent-of-code-2021/day1/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var prev int
	var result int
	for scanner.Scan() {
		line := scanner.Text()
		current := readLine(line)
		if current > prev {
			result++
		}
		prev = current
	}
	result = result - 1

	fmt.Printf("Result %+v\n", result)
}

func readLine(input string) int {
	i, err := strconv.Atoi(input)
	if err != nil {
		panic(err)
	}
	return i
}
