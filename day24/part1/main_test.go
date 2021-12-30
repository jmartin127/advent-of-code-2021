package main

import (
	"testing"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

func TestIsValidModelNumber(t *testing.T) {
	toCheck := []int{
		99691891957938,
		97691291357918,
		98691591657928,
		77691891957916,
		47691591657913,
		88691391457927,
		38691891957922,
		38691891957922,
		47691691757913,
		38691791857922,
		58691891957924,
	}

	filepath := "../input.txt"
	list := helpers.ReadFile(filepath)
	instructionGroups := parseInput(list)

	for _, tc := range toCheck {
		if !isValidModelNumber(tc, instructionGroups) {
			t.Fatalf("Model number should have passed the check %d", tc)
		}
	}
}
