package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

const STOPPING_SCORE = 21

type cacheResult struct {
	p1Wins uint64
	p2Wins uint64
}

var cache = map[string]*cacheResult{}

func cacheKey(player1Space, player2Space, currentDiceRoll, player1Score, player2Score int, player1Turn bool, numPlayerInRow int, totalRoll int) string {
	return fmt.Sprintf("%d,%d,%d,%d,%d,%t,%d,%d", player1Space, player2Space, currentDiceRoll, player1Score, player2Score, player1Turn, numPlayerInRow, totalRoll)
}

func main() {
	list := helpers.ReadFile("day21/input.txt")
	player1StartingSpace, player2StartingSpace := parseInput(list)
	p1wins1, p2wins1 := rollDice(player1StartingSpace, player2StartingSpace, 0, 0, true, 0, 0)
	fmt.Printf("player1Wins=%d\n", p1wins1)
	fmt.Printf("player2Wins=%d\n", p2wins1)
}

func rollDice(player1Space, player2Space, player1Score, player2Score int, player1Turn bool, numPlayerInRow int, totalRoll int) (uint64, uint64) {
	p1wins1, p2wins1 := playToStoppingScore(player1Space, player2Space, 1, player1Score, player2Score, player1Turn, numPlayerInRow, totalRoll)
	p1wins2, p2wins2 := playToStoppingScore(player1Space, player2Space, 2, player1Score, player2Score, player1Turn, numPlayerInRow, totalRoll)
	p1wins3, p2wins3 := playToStoppingScore(player1Space, player2Space, 3, player1Score, player2Score, player1Turn, numPlayerInRow, totalRoll)

	return p1wins1 + p1wins2 + p1wins3, p2wins1 + p2wins2 + p2wins3
}

func playToStoppingScore(player1Space, player2Space, currentDiceRoll, player1Score, player2Score int, player1Turn bool, numPlayerInRow int, totalRoll int) (uint64, uint64) {
	// check cache
	cacheKey := cacheKey(player1Space, player2Space, currentDiceRoll, player1Score, player2Score, player1Turn, numPlayerInRow, totalRoll)
	if v, ok := cache[cacheKey]; ok {
		return v.p1Wins, v.p2Wins
	}

	totalRoll += currentDiceRoll
	numPlayerInRow++
	if numPlayerInRow >= 3 {
		if player1Turn {
			player1Space = determineSpace(player1Space, totalRoll)
			player1Score += player1Space
			if player1Score >= STOPPING_SCORE {
				cache[cacheKey] = &cacheResult{p1Wins: 1, p2Wins: 0}
				return uint64(1), uint64(0)
			}
		} else {
			player2Space = determineSpace(player2Space, totalRoll)
			player2Score += player2Space
			if player2Score >= STOPPING_SCORE {
				cache[cacheKey] = &cacheResult{p1Wins: 0, p2Wins: 1}
				return uint64(0), uint64(1)
			}
		}
		player1Turn = !player1Turn
		numPlayerInRow = 0
		totalRoll = 0
	}

	// add result to cache
	p1Wins, p2Wins := rollDice(player1Space, player2Space, player1Score, player2Score, player1Turn, numPlayerInRow, totalRoll)
	cache[cacheKey] = &cacheResult{p1Wins: p1Wins, p2Wins: p2Wins}

	return p1Wins, p2Wins
}

func determineSpace(currentSpace int, roll int) int {
	result := (currentSpace + roll) % 10
	if result == 0 {
		result = 10
	}
	return result
}

func parseInput(list []string) (int, int) {
	p1parts := strings.Split(list[0], ": ") // Player 1 starting position: 4
	p2parts := strings.Split(list[1], ": ") // Player 1 starting position: 8

	p1, _ := strconv.Atoi(p1parts[1])
	p2, _ := strconv.Atoi(p2parts[1])

	return p1, p2
}
