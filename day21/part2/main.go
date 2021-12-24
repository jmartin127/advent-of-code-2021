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
	playToStoppingScore(player1StartingSpace, player2StartingSpace, 1, 0, 0, true, 0, 0)
	playToStoppingScore(player1StartingSpace, player2StartingSpace, 2, 0, 0, true, 0, 0)
	playToStoppingScore(player1StartingSpace, player2StartingSpace, 3, 0, 0, true, 0, 0)
	fmt.Printf("player1Wins=%d\n", player1Wins)
	fmt.Printf("player2Wins=%d\n", player2Wins)
}

func playToStoppingScore(player1Space, player2Space, currentDiceRoll, player1Score, player2Score int, player1Turn bool, numPlayerInRow int, totalRoll int) {
	// base case
	if player1Score >= STOPPING_SCORE || player2Score >= STOPPING_SCORE {
		if player1Score > player2Score {
			player1Wins++
		} else {
			player2Wins++
		}
		if player1Wins%1000000000 == 0 {
			fmt.Printf("player1Wins=%d. player2Wins=%d. \n", player1Wins, player2Wins)
		}
		//fmt.Printf("Rolls %+v. player1Score=%d. player2Score=%d. \n", allRolls, player1Score, player2Score)
		return
	}

	totalRoll += currentDiceRoll
	numPlayerInRow++
	if numPlayerInRow >= 3 {
		if player1Turn {
			player1Space = (player1Space + totalRoll) % 10
			player1Score += player1Space
		} else {
			player2Space = (player2Space + totalRoll) % 10
			player2Score += player2Space
		}
		player1Turn = !player1Turn
		numPlayerInRow = 0
		totalRoll = 0
	}

	playToStoppingScore(player1Space, player2Space, 1, player1Score, player2Score, player1Turn, numPlayerInRow, totalRoll)
	playToStoppingScore(player1Space, player2Space, 2, player1Score, player2Score, player1Turn, numPlayerInRow, totalRoll)
	playToStoppingScore(player1Space, player2Space, 3, player1Score, player2Score, player1Turn, numPlayerInRow, totalRoll)
}

// func determineSpace(currentSpace int, roll int) int {
// 	return (currentSpace + roll) % 10
// }

func parseInput(list []string) (int, int) {
	p1parts := strings.Split(list[0], ": ") // Player 1 starting position: 4
	p2parts := strings.Split(list[1], ": ") // Player 1 starting position: 8

	p1, _ := strconv.Atoi(p1parts[1])
	p2, _ := strconv.Atoi(p2parts[1])

	return p1, p2
}
