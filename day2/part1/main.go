package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

func main() {
	filepath := "/Users/jeff.martin@getweave.com/go/src/github.com/jmartin127/advent-of-code-2021/day2/input.txt"
	list := helpers.ReadFile(filepath)

	horizontal := 0
	vertical := 0
	for _, v := range list {
		s := strings.Split(v, " ")
		direction := s[0]
		d := s[1]
		distance, _ := strconv.Atoi(d)

		if direction == "forward" {
			horizontal = horizontal + distance
		} else if direction == "down" {
			vertical = vertical + distance
		} else if direction == "up" {
			vertical = vertical - distance
		}
	}

	fmt.Printf("result %d\n", horizontal*vertical)
}
