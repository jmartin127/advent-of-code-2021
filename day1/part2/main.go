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

	// create the initial list
	initialList := make([]int, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		current := readLine(line)
		initialList = append(initialList, current)
	}

	// create the windows
	windows := make([]int, 0)
	for i := 1; i < len(initialList)-1; i++ {
		sum := initialList[i-1] + initialList[i] + initialList[i+1]
		windows = append(windows, sum)
	}

	// loop through the windows
	var prev int
	var result int
	for _, current := range windows {
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
