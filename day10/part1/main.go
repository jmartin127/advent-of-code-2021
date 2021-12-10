package main

import (
	"fmt"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

func main() {
	filepath := "input.txt"
	list := helpers.ReadFile(filepath)

	answer := 0
	for _, line := range list {
		firstIllegal := findFirstIllegal(line)
		points := determinePoints(firstIllegal)
		if points != 0 {
			continue
		}
		fmt.Printf("poitns %d\n", points)
		answer += points
	}

	fmt.Printf("Answer %d\n", answer)
}

func determinePoints(input string) int {
	if input == ")" {
		return 3
	} else if input == "]" {
		return 57
	} else if input == "}" {
		return 1197
	} else if input == ">" {
		return 25137
	}

	return 0
}

func findFirstIllegal(line string) string {
	stack := make([]string, 0)
	for _, char := range strings.Split(line, "") {
		if isOpenChar(char) {
			stack = append(stack, char)
		} else {
			popped, wasPopped, stackNew := popLast(stack)
			stack = stackNew
			if !wasPopped {
				return char
			} else {
				if !isMatchingCloseChar(popped, char) {
					return char
				}
			}
		}
	}

	return ""
}

func isOpenChar(element string) bool {
	return element == "(" || element == "[" || element == "{" || element == "<"
}

func isMatchingCloseChar(openChar, closeChar string) bool {
	return (openChar == "(" && closeChar == ")") || (openChar == "[" && closeChar == "]") || (openChar == "{" && closeChar == "}") || (openChar == "<" && closeChar == ">")
}

func popLast(input []string) (string, bool, []string) {
	if len(input) == 0 {
		return "", false, input
	}

	result := input[len(input)-1]
	input = input[:len(input)-1]
	return result, true, input
}
