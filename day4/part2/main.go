package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

type tile struct {
	number        int
	hasBeenCalled bool
}

type board struct {
	rows [][]*tile
}

func main() {
	bingoVals, boards := parseInput()

	lastNumCalled, winningBoard := run(bingoVals, boards)
	sumOfUnmarked := winningBoard.sumOfUnmarked()
	answer := lastNumCalled * sumOfUnmarked
	fmt.Printf("Answer: %d\n", answer)
}

func run(bingoVals []int, boards []*board) (int, *board) {
	numBoards := len(boards)
	boardsWhichWon := make(map[int]bool, 0)
	for _, numCalled := range bingoVals {
		for boardIndex, b := range boards {
			b.updateBoard(numCalled)
			if b.hasBingo() {
				boardsWhichWon[boardIndex] = true
				if len(boardsWhichWon) == numBoards {
					return numCalled, b
				}
			}
		}
	}

	return 0, nil
}

func (b *board) updateBoard(numCalled int) bool {
	for _, row := range b.rows {
		for _, t := range row {
			if t.number == numCalled {
				t.hasBeenCalled = true
				return true
			}
		}
	}

	return false
}

func (b *board) hasBingo() bool {
	// check rows
	for _, row := range b.rows {
		rowHasBingo := true
		for _, t := range row {
			if !t.hasBeenCalled {
				rowHasBingo = false
				break
			}
		}
		if rowHasBingo {
			return true
		}
	}

	// check columns
	for colIndex := 0; colIndex < len(b.rows); colIndex++ {
		rowHasBingo := true
		for _, row := range b.rows {
			if !row[colIndex].hasBeenCalled {
				rowHasBingo = false
				break
			}
		}
		if rowHasBingo {
			return true
		}
	}

	return false
}

func (b *board) sumOfUnmarked() int {
	var sum int
	for _, row := range b.rows {
		for _, t := range row {
			if !t.hasBeenCalled {
				sum += t.number
			}
		}
	}
	return sum
}

// 7,4,9,5,11,17,23,2,0,14,21,24,10,16,13,6,15,25,12,22,18,20,8,19,3,26,1

// 22 13 17 11  0
//  8  2 23  4 24
// 21  9 14 16  7
func parseInput() ([]int, []*board) {
	filepath := "/Users/jeff.martin@getweave.com/go/src/github.com/jmartin127/advent-of-code-2021/day4/input.txt"
	list := helpers.ReadFile(filepath)

	numbersCalled := createNumbersCalled(list[0])
	boards := make([]*board, 0)
	currentBoard := &board{}
	for i := 2; i < len(list); i++ {
		line := list[i]
		if line == "" {
			boards = append(boards, currentBoard)
			currentBoard = &board{}
		} else {
			r := createRow(line)
			currentBoard.rows = append(currentBoard.rows, r)
		}
	}
	boards = append(boards, currentBoard)

	return numbersCalled, boards
}

// 21  9 14 16  7
func createRow(line string) []*tile {
	vals := strings.Fields(line)
	result := make([]*tile, 0)
	for _, v := range vals {
		i, _ := strconv.Atoi(v)
		result = append(result, &tile{number: i})
	}
	return result
}

func createNumbersCalled(line string) []int {
	vals := strings.Split(line, ",")
	result := make([]int, 0)
	for _, v := range vals {
		i, _ := strconv.Atoi(v)
		result = append(result, i)
	}

	return result
}
