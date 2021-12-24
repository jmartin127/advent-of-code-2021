package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

const STOPPING_SCORE = 21

var player1Wins = 0
var player2Wins = 0

func main() {
	list := helpers.ReadFile("day21/input.txt")
	player1StartingSpace, player2StartingSpace := parseInput(list)
	playToStoppingScore(player1StartingSpace, player2StartingSpace, 1, 0, 0, true, []int{})
	playToStoppingScore(player1StartingSpace, player2StartingSpace, 2, 0, 0, true, []int{})
	playToStoppingScore(player1StartingSpace, player2StartingSpace, 3, 0, 0, true, []int{})
	fmt.Printf("player1Wins=%d\n", player1Wins)
	fmt.Printf("player2Wins=%d\n", player2Wins)
}

func playToStoppingScore(player1Space, player2Space, currentDiceRoll, player1Score, player2Score int, player1Turn bool, allRolls []int) []int {
	// base case
	if player1Score >= STOPPING_SCORE || player2Score >= STOPPING_SCORE {
		if player1Score > player2Score {
			player1Wins++
		} else {
			player2Wins++
		}
		fmt.Printf("Rolls %+v. player1Score=%d. player2Score=%d. \n", allRolls, player1Score, player2Score)
		return allRolls
	}

	allRolls = append(allRolls, currentDiceRoll)
	if player1Turn {
		player1Space = determineSpace(player1Space, currentDiceRoll)
		player1Score += player1Space
	} else {
		player2Space = determineSpace(player2Space, currentDiceRoll)
		player2Score += player2Space
	}
	player1Turn = !player1Turn

	playToStoppingScore(player1Space, player2Space, 1, player1Score, player2Score, player1Turn, allRolls)
	playToStoppingScore(player1Space, player2Space, 2, player1Score, player2Score, player1Turn, allRolls)
	playToStoppingScore(player1Space, player2Space, 3, player1Score, player2Score, player1Turn, allRolls)

	return allRolls
}

func determineSpace(currentSpace int, roll int) int {
	currentSpace += roll
	if currentSpace > 10 {
		currentSpace = currentSpace - 10
	}
	return currentSpace
}

func parseInput(list []string) (int, int) {
	p1parts := strings.Split(list[0], ": ") // Player 1 starting position: 4
	p2parts := strings.Split(list[1], ": ") // Player 1 starting position: 8

	p1, _ := strconv.Atoi(p1parts[1])
	p2, _ := strconv.Atoi(p2parts[1])

	return p1, p2
}
