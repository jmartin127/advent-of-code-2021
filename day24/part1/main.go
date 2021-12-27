package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

type instruction struct {
	cmd              string
	operand1         string
	hasSecondOperand bool
	operand2IsNum    bool
	operand2         string
	operand2Num      int
}

var posByLetter = map[string]int{
	"w": 0,
	"x": 1,
	"y": 2,
	"z": 3,
}

var vals = []int{0, 0, 0, 0}

func main() {
	// read the instructions
	filepath := "day24/input.txt"
	list := helpers.ReadFile(filepath)
	instructions := parseInput(list)

	isValid := isValidModelNumber(6, instructions)
	fmt.Printf("Valid=%t\n", isValid)
	fmt.Printf("Result %+v", vals)
}

func isValidModelNumber(modelNum int, instructions []*instruction) bool {
	checkModelNumber(modelNum, instructions)
	if vals[posByLetter["z"]] == 0 {
		return true
	}
	return false
}

// Example model number: 13579246899999
func checkModelNumber(modelNumStr int, instructions []*instruction) {
	modelNum := strings.Split(strconv.Itoa(modelNumStr), "")
	var modelNumPointer int

	for _, i := range instructions {
		var inputValInt int
		if modelNumPointer <= len(modelNum)-1 {
			inputValInt, _ = strconv.Atoi(modelNum[modelNumPointer])
		}
		usedInputVal := applyInstruction(i, inputValInt)
		if usedInputVal {
			modelNumPointer++
		}
	}
}

/*
 */
func applyInstruction(i *instruction, inputVal int) bool {
	var usedInputVal bool

	switch i.cmd {
	case "inp": // inp a - Read an input value and write it to variable a.
		vals[posByLetter[i.operand1]] = inputVal
		usedInputVal = true
	case "add": // add a b - Add the value of a to the value of b, then store the result in variable a.
		if i.operand2IsNum {
			vals[posByLetter[i.operand1]] += i.operand2Num
		} else {
			vals[posByLetter[i.operand1]] += vals[posByLetter[i.operand2]]
		}
	case "mul": // mul a b - Multiply the value of a by the value of b, then store the result in variable a.
		if i.operand2IsNum {
			vals[posByLetter[i.operand1]] *= i.operand2Num
		} else {
			vals[posByLetter[i.operand1]] *= vals[posByLetter[i.operand2]]
		}
	case "div": // div a b - Divide the value of a by the value of b, truncate the result to an integer, then store the result in variable a. (Here, "truncate" means to round the value toward zero.)
		if i.operand2IsNum {
			vals[posByLetter[i.operand1]] /= i.operand2Num
		} else {
			vals[posByLetter[i.operand1]] /= vals[posByLetter[i.operand2]]
		}
	case "mod": // mod a b - Divide the value of a by the value of b, then store the remainder in variable a. (This is also called the modulo operation.)
		if i.operand2IsNum {
			vals[posByLetter[i.operand1]] %= i.operand2Num
		} else {
			vals[posByLetter[i.operand1]] %= vals[posByLetter[i.operand2]]
		}
	case "eql": // eql a b - If the value of a and b are equal, then store the value 1 in variable a. Otherwise, store the value 0 in variable a.
		var equalVal int
		if i.operand2IsNum {
			if vals[posByLetter[i.operand1]] == i.operand2Num {
				equalVal = 1
			}
		} else {
			if vals[posByLetter[i.operand1]] == vals[posByLetter[i.operand2]] {
				equalVal = 1
			}
		}
		vals[posByLetter[i.operand1]] = equalVal
	}

	return usedInputVal
}

func parseInput(list []string) []*instruction {
	instructions := make([]*instruction, 0)
	for _, line := range list {
		instructions = append(instructions, parseLine(line))
	}
	return instructions
}

/*
inp w
mul x 0
add x z
mod x 26
div z 1
add x 13
eql x w
eql x 0
*/
func parseLine(line string) *instruction {
	parts := strings.Split(line, " ")

	var hasSecondOperand bool
	var operand2 string
	var operand2Num int
	var operand2IsNum bool
	if len(parts) > 2 {
		hasSecondOperand = true
		operand2Str := parts[2]
		if operand2Str == "w" || operand2Str == "x" || operand2Str == "y" || operand2Str == "z" {
			operand2 = operand2Str
		} else {
			operand2IsNum = true
			v, _ := strconv.Atoi(operand2Str)
			operand2Num = v
		}
	}

	return &instruction{
		cmd:              parts[0],
		operand1:         parts[1],
		hasSecondOperand: hasSecondOperand,
		operand2:         operand2,
		operand2Num:      operand2Num,
		operand2IsNum:    operand2IsNum,
	}
}
