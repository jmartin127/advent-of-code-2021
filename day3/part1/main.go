package main

import (
	"fmt"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

func main() {
	filepath := "/Users/jeff.martin@getweave.com/go/src/github.com/jmartin127/advent-of-code-2021/day3/input.txt"
	list := helpers.ReadFile(filepath)

	counts := make(map[int]int, 0)
	total := 0
	for _, v := range list {
		total++
		chars := strings.Split(v, "") // 001011100010
		for i, c := range chars {
			if c == "0" {
				counts[i] = counts[i] + 1
			}
		}
	}

	for i := 0; i < len(counts); i++ {
		v := counts[i]
		if v > (total - v) {
			fmt.Printf("0")
		} else {
			fmt.Printf("1")
		}
	}
	fmt.Printf("\n")

}
