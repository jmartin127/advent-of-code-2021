package helpers

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func ReadFile(filepath string) []string {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	result := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)
	}

	return result
}

func ReadFileAsInts(filepath string) []int {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	result := make([]int, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, readLine(line))
	}

	return result
}

func ReadSingleLineFileAsInts(filepath string) []int {
	list := ReadFile(filepath)

	// convert to ints
	firstLineParts := strings.Split(list[0], ",")
	result := make([]int, 0)
	for _, val := range firstLineParts {
		intVal, _ := strconv.Atoi(val)
		result = append(result, intVal)
	}

	return result
}

func MaxFromList(list []int) int {
	var max int
	for _, val := range list {
		if val > max {
			max = val
		}
	}
	return max
}

func readLine(input string) int {
	i, err := strconv.Atoi(input)
	if err != nil {
		panic(err)
	}
	return i
}
