package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

func main() {
	filepath := "input.txt"
	list := helpers.ReadFile(filepath)

	remainders := make([][]string, 0)
	for _, line := range list {
		firstIllegal, remainingStack := findFirstIllegal(line)
		points := determinePoints(firstIllegal)
		if points != 0 {
			continue
		}

		toClose := make([]string, 0)
		lengthOfStack := len(remainingStack)
		for i := 0; i < lengthOfStack; i++ {
			popped, _, stackNew := popLast(remainingStack)
			remainingStack = stackNew
			close := findCloseForOpen(popped)
			toClose = append(toClose, close)
		}

		remainders = append(remainders, toClose)
		fmt.Printf("adding remainder %s Stack len %+v HERE: %+v\n", line, remainingStack, toClose)
	}

	answers := make([]int, 0)
	for _, rem := range remainders {
		answer := scoreCompletionString(rem)
		fmt.Printf("answer %d\n", answer)
		answers = append(answers, answer)
	}

	sort.Ints(answers)
	middle := len(answers) / 2
	fmt.Printf("answer %d\n", answers[middle])
}

func scoreCompletionString(input []string) int {
	result := 0
	for _, v := range input {
		result *= 5
		result += closingPointsForChar(v)
	}
	return result
}

func closingPointsForChar(input string) int {
	if input == ")" {
		return 1
	} else if input == "]" {
		return 2
	} else if input == "}" {
		return 3
	} else if input == ">" {
		return 4
	}
	return 0
}

func findCloseForOpen(open string) string {
	if open == "(" {
		return ")"
	} else if open == "[" {
		return "]"
	} else if open == "{" {
		return "}"
	} else if open == "<" {
		return ">"
	}

	return ""
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

func findFirstIllegal(line string) (string, []string) {
	stack := make([]string, 0)
	for _, char := range strings.Split(line, "") {
		if isOpenChar(char) {
			stack = append(stack, char)
		} else {
			popped, wasPopped, stackNew := popLast(stack)
			stack = stackNew
			if !wasPopped {
				return char, stack
			} else {
				if !isMatchingCloseChar(popped, char) {
					return char, stack
				}
			}
		}
	}

	return "", stack
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
