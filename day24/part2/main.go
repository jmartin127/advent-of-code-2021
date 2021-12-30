package main

import (
	"fmt"
	"log"
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

var minModelNum = 27691591657911 + 1

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

	var numTried int
	for true {
		numTried++
		modelNum, success := randomlyGenerateNewModelNum()
		if numTried%100000 == 0 {
			log.Printf("Num tried %d\n", numTried)
		}
		if !success {
			continue
		}
		if modelNum < minModelNum {
			isValidModelNumber(modelNum, instructionGroups, true)
		}
	}
}

// Originally NN69NN91NN799N
var hardcoded = map[int]int{
	0: 2,
	1: 7,
	2: 6,
	3: 9,
	//4:  1,
	//5:  -1,
	//6: 9,
	//7:  1,
	//8:  -1,
	//9:  5,
	//10: 7,
	11: 1, // will get replaced by lower value
	12: 1, // will get replaced by lower value
	13: 1, // will get replaced by lower value
}

// positions where we should only try low numbers (1,2,3,4)
var positionsToOptimize = map[int]bool{
	//0: true,
	//1: true,
}

// these changes tend to lead to the answer, based on analysis of the input
// the key is the position in the model number, the value is what should
// change on the subsequent position
// HIGH was: 99691891979938
// Current lowest: 27691591657911
// New lowest:     27691191279911
// New lowest:     27691191224911
// New lowest:     27691191213911
var preferences = map[int]int{
	//3: 3,
	// 7: -8,
	//8:  -7,
	//11: 2,
	//12: -4,
	//13: 7,
}

var specialCases = map[int]int{
	13: -1,
	12: -9,
	11: -5,
	10: 0,
	8:  -8,
	7:  -14,
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
			toAppend = generateRandomNumber(1, 4) // TODO maybe only 1,2,3
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

func isValidModelNumber(modelNum int, instructionGroups []*instructionGroup, fixNumber bool) bool {
	vals = []int{0, 0, 0, 0} // reset the output vars
	return checkModelNumber(modelNum, instructionGroups, fixNumber)
}

// Example model number: 13579246899999
func checkModelNumber(modelNumStr int, instructionGroups []*instructionGroup, fixNumber bool) bool {
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
		currentZ, success, newModelNum := applyInstructionGroup(ig, modelNum, i, priorZ, fixNumber)
		modelNum = newModelNum
		if !success {
			return false
		}
		priorZ = currentZ
	}

	if vals[posByLetter["z"]] == 0 {
		answer := convertIntArrayToInt(modelNum)
		if answer < minModelNum {
			minModelNum = answer
			fmt.Printf("NEW MIN %d\n", answer)
		}
		return true
	}
	return false
}

func convertIntArrayToInt(input []int) int {
	var result string
	for _, i := range input {
		result += strconv.Itoa(i)
	}
	resultingInteger, _ := strconv.Atoi(result)
	return resultingInteger
}

// return vals are: newZ, success, newModelNum
func applyInstructionGroup(ig *instructionGroup, modelNum []int, modelNumPointer int, priorZ int, fixNumber bool) (int, bool, []int) {
	if fixNumber {
		if v, ok := specialCases[modelNumPointer]; ok {
			inputMustBe := (priorZ % 26) + v
			if !checkRange(inputMustBe) {
				return -1, false, modelNum
			}
			modelNum[modelNumPointer] = inputMustBe
		}
	}

	for _, ins := range ig.group {
		success := applyInstruction(ins, modelNum, modelNumPointer, priorZ)
		if !success {
			return -1, false, modelNum
		}
	}

	// after applying the group, return the z value
	return vals[posByLetter["z"]], true, modelNum
}

/*
 */
func applyInstruction(i *instruction, modelNum []int, modelNumPointer int, priorZ int) bool {
	switch i.cmd {
	case "inp": // inp a - Read an input value and write it to variable a.
		vals[posByLetter[i.operand1]] = modelNum[modelNumPointer]
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
