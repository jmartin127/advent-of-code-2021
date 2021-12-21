package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

const NUM_SIDES_OF_DIE = 100
const STOPPING_SCORE = 1000

func main() {
	list := helpers.ReadFile("input.txt")
	player1StartingSpace, player2StartingSpace := parseInput(list)
	answer := playToNPoints(STOPPING_SCORE, player1StartingSpace, player2StartingSpace)
	fmt.Printf("Answer %d\n", answer)
}

func playToNPoints(nPoints int, player1StartingSpace, player2StartingSpace int) int {
	player1Points := 0
	player2Points := 0

	player1Turn := true
	currentDieValue := 0
	rollIncrement := 3
	player1Space := player1StartingSpace - 1 // 0-based
	player2Space := player2StartingSpace - 1 // 0-based
	numTimesRolled := 0
	for player1Points < STOPPING_SCORE && player2Points < STOPPING_SCORE {
		currentDieValue += rollIncrement
		rolled := toAdd(((currentDieValue - 2) % NUM_SIDES_OF_DIE), ((currentDieValue - 1) % NUM_SIDES_OF_DIE), (currentDieValue % NUM_SIDES_OF_DIE))
		if player1Turn {
			player1Space += rolled
			player1Space = player1Space % 10
			player1Points += player1Space + 1
		} else {
			player2Space += rolled
			player2Space = player2Space % 10
			player2Points += player2Space + 1
		}
		player1Turn = !player1Turn
		numTimesRolled += 3
	}

	losingScore := player1Points
	if player2Points < player1Points {
		losingScore = player2Points
	}
	fmt.Printf("numTimesRolled %d\n", numTimesRolled)
	fmt.Printf("losingScore %d\n", losingScore)
	return numTimesRolled * losingScore
}

func toAdd(val1, val2, val3 int) int {
	if val1 == 0 {
		val1 = 100
	} else if val2 == 0 {
		val2 = 100
	} else if val3 == 0 {
		val3 = 100
	}
	fmt.Printf("Die values %d %d %d...\n", val1, val2, val3)
	return val1 + val2 + val3
}

func parseInput(list []string) (int, int) {
	p1parts := strings.Split(list[0], ": ") // Player 1 starting position: 4
	p2parts := strings.Split(list[1], ": ") // Player 1 starting position: 8

	p1, _ := strconv.Atoi(p1parts[1])
	p2, _ := strconv.Atoi(p2parts[1])

	return p1, p2
}
