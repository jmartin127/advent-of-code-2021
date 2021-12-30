package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

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

// The instructions naturally fall into 14 groups, 1 for each digit
type instructionGroup struct {
	group []*instruction
}

const MODEL_NUM_LEN = 14

var maxModelNum = 99691891957938 - 1
var minZ = 1000000

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
	instructionGroups := parseInput(list)

	isValid := isValidModelNumber(99691891957938, instructionGroups)
	fmt.Printf("Is valid %t\n", isValid)

	// var numTried int
	// for true {
	// 	numTried++
	// 	modelNum, success := randomlyGenerateNewModelNum()
	// 	if numTried%100000 == 0 {
	// 		log.Printf("Num tried %d\n", numTried)
	// 	}
	// 	if !success {
	// 		continue
	// 	}
	// 	//fmt.Printf("Model num %d\n", modelNum)
	// 	if modelNum > maxModelNum {
	// 		vals = []int{0, 0, 0, 0} // reset the output vars
	// 		isValid := isValidModelNumber(modelNum, instructions)
	// 		if isValid {
	// 			fmt.Printf("%d\n", modelNum)
	// 			maxModelNum = modelNum
	// 		}
	// 		if vals[posByLetter["z"]] < minZ {
	// 			minZ = vals[posByLetter["z"]]
	// 			fmt.Printf("New minZ %d\n", minZ)
	// 		}
	// 	}
	// }

	// for i := 1; i < 10; i++ {
	// 	vals = []int{0, 0, 0, 0} // reset the output vars
	// 	isValid := isValidModelNumber(i, instructions)
	// 	fmt.Printf("MODEL_NUM %d. Is Valid %t, Vals=%+v\n", i, isValid, vals)
	// }
}

// Following this pattern NN69NN91NN799N
// Best number so far: 99691891957938
var hardcoded = map[int]int{
	0: 9,
	1: 9,
	2: 6,
	3: 9,
	//4:  1,
	6: 9,
	//7:  1,
	//9:  5,
	10: 7,
	//11: 9,
	//12: 9,
}

// positions where we should only try 6-9
var positionsToOptimize = map[int]bool{
	//: true,
}

// these changes tend to lead to the answer, based on analysis of the input
// the key is the position in the model number, the value is what should
// change on the subsequent position
var preferences = map[int]int{
	//3: 3,
	//7: -8,
	//8:  -7,
	//11: 2,
	//12: -4,
	//13: 7,
}

func randomlyChooseFromSlice(i []int, num int) []int {
	input := make([]int, 0)
	for _, v := range i {
		input = append(input, v)
	}

	result := make([]int, 0)
	for num > len(result) {
		r := generateRandomNumber(0, len(input)-1)
		result = append(result, input[r])
		remove(input, r)
	}

	return result
}

func remove(input []int, s int) []int {
	return append(input[:s], input[s+1:]...)
}

func getKeys(m map[int]int) []int {
	keys := make([]int, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func getPrefsToUse() map[int]int {
	keys := getKeys(preferences)
	keysToUse := randomlyChooseFromSlice(keys, 6)
	result := make(map[int]int, 0)
	for _, k := range keysToUse {
		result[k] = preferences[k]
	}
	return result
}

func randomlyGenerateNewModelNum() (int, bool) {
	//preferencesToUse := getPrefsToUse() // TODO
	preferencesToUse := preferences // TODO
	//fmt.Printf("using %+v\n", preferencesToUse)

	newModelNum := make([]int, 0)
	for i := 0; i < MODEL_NUM_LEN; i++ {
		var toAppend int
		if v, ok := hardcoded[i]; ok {
			toAppend = v
		} else if _, ok := positionsToOptimize[i]; ok {
			toAppend = generateRandomNumber(6, 9) // TODO maybe only 7-9
		} else if v, ok := preferencesToUse[i]; ok {
			prevNum := newModelNum[i-1]
			toAppend = prevNum + v
			if !checkRange(toAppend) {
				return -1, false
			}
		} else {
			toAppend = generateRandomNumber(1, 9)
		}
		newModelNum = append(newModelNum, toAppend)
	}
	newModelNumStr := strings.Trim(strings.Replace(fmt.Sprint(newModelNum), " ", "", -1), "[]")

	result, _ := strconv.Atoi(newModelNumStr)
	return result, true
}

func checkRange(v int) bool {
	return v >= 1 && v <= 9
}

func generateRandomNumber(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func isValidModelNumber(modelNum int, instructionGroups []*instructionGroup) bool {
	vals = []int{0, 0, 0, 0} // reset the output vars
	return checkModelNumber(modelNum, instructionGroups)
}

// Example model number: 13579246899999
func checkModelNumber(modelNumStr int, instructionGroups []*instructionGroup) bool {
	// convert model number to slice of int
	modelNumStrParts := strings.Split(strconv.Itoa(modelNumStr), "")
	modelNum := make([]int, 0)
	for _, m := range modelNumStrParts {
		inputValInt, _ := strconv.Atoi(m)
		modelNum = append(modelNum, inputValInt)
	}

	// apply each group
	var priorZ int
	for i, ig := range instructionGroups {
		currentZ, success := applyInstructionGroup(ig, modelNum, i, priorZ)
		if !success {
			fmt.Printf("False for ig index %d\n", i)
			return false
		}
		priorZ = currentZ
	}

	return true
}

func applyInstructionGroup(ig *instructionGroup, modelNum []int, modelNumPointer int, priorZ int) (int, bool) {
	// FOUND: For the last digit to work, priorZ must be <= 10 && >= 2... BECAUSE we are subtracting 1 after doing the mod
	if modelNumPointer == 13 {
		if priorZ < 2 || priorZ > 10 {
			return -1, false
		} else {
			return 0, true
		}
	}

	for _, ins := range ig.group {
		success := applyInstruction(ins, modelNum, modelNumPointer, priorZ)
		if !success {
			return -1, false
		}
	}
	// after applying the group, return the z value
	return vals[posByLetter["z"]], true
}

func applyCustomModelNumAtPos(defaultModelNum int, pos int, priorZ int) (int, bool) {
	// FOUND: For the last digit to work... Need to set w=mod(z from prior, 26)-1.  If w is out of range at this point, the model number won't work out.
	if pos == 13 { // last
		mustBeVal := (priorZ % 26) - 1
		if checkRange(mustBeVal) {
			fmt.Printf("Setting to %d, for priorZ %d\n", mustBeVal, priorZ)
			return mustBeVal, true
		}
		return -1, false
	}
	return defaultModelNum, true
}

/*
 */
func applyInstruction(i *instruction, modelNum []int, modelNumPointer int, priorZ int) bool {
	switch i.cmd {
	case "inp": // inp a - Read an input value and write it to variable a.
		customVal, success := applyCustomModelNumAtPos(modelNum[modelNumPointer], modelNumPointer, priorZ)
		if !success {
			return false
		}
		vals[posByLetter[i.operand1]] = customVal
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
		// handle ASCII
		//fmt.Printf("Converting %d to ascii\n", vals[posByLetter[i.operand1]])
		//character := string(65 + vals[posByLetter[i.operand1]])
		//fmt.Printf("ASCII %s\n", character)
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

	return true
}

func parseInput(list []string) []*instructionGroup {
	result := make([]*instructionGroup, 0)
	ig := &instructionGroup{
		group: make([]*instruction, 0),
	}
	for _, line := range list {
		ins := parseLine(line)
		if ins.cmd == "inp" && len(ig.group) > 0 {
			result = append(result, ig)
			ig = &instructionGroup{
				group: make([]*instruction, 0),
			}
		}
		ig.group = append(ig.group, ins)
	}
	result = append(result, ig)
	return result
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
