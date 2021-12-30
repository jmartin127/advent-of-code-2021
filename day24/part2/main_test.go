package main

import (
	"testing"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

func TestIsValidModelNumber(t *testing.T) {
	filepath := "../input.txt"
	list := helpers.ReadFile(filepath)
	instructionGroups := parseInput(list)

	toCheck := []int{
		27691191213911,
		99691891979938,
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

	for _, tc := range toCheck {
		if !isValidModelNumber(tc, instructionGroups, false) {
			t.Fatalf("Model number should have passed the check %d", tc)
		}
	}
	for _, tc := range toCheck {
		if !isValidModelNumber(tc, instructionGroups, true) {
			t.Fatalf("Model number should have passed the check %d", tc)
		}
	}

	// should not pass
	toCheck = []int{
		13579246899999,
	}
	for _, tc := range toCheck {
		if isValidModelNumber(tc, instructionGroups, false) {
			t.Fatalf("Model number should NOT have passed the check %d", tc)
		}
	}

	// should pass (last 2 digits auto-filled)
	toCheck = []int{
		99691891957911,
	}
	for _, tc := range toCheck {
		if !isValidModelNumber(tc, instructionGroups, true) {
			t.Fatalf("Model number should have been auto filled %d", tc)
		}
	}
}
