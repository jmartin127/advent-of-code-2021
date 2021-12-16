package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

var hexToBinary = map[string]string{
	"0": "0000",
	"1": "0001",
	"2": "0010",
	"3": "0011",
	"4": "0100",
	"5": "0101",
	"6": "0110",
	"7": "0111",
	"8": "1000",
	"9": "1001",
	"A": "1010",
	"B": "1011",
	"C": "1100",
	"D": "1101",
	"E": "1110",
	"F": "1111",
}

func main() {
	filepath := "input.txt"
	list := helpers.ReadFile(filepath)

	input := list[0]
	fmt.Printf("Input %s\n", input)

	binaryStr := convertToBinary(input)
	fmt.Printf("Binary %s\n", binaryStr)

	result, _ := run(binaryStr, 0)
	var total int
	for _, v := range result {
		total = total + v
	}
	fmt.Printf("Answer %d\n", total)
}

func run(binaryStr string, startIndex int) ([]int, int) {
	subPacketValues := make([]int, 0)

	currentIndex := startIndex
	for true {
		if currentIndex >= len(binaryStr)-1 || isRemainingZero(binaryStr, currentIndex) {
			return subPacketValues, currentIndex
		}
		newIndex, version, theType := readHeader(binaryStr, currentIndex)
		currentIndex = newIndex
		fmt.Println("******************************")
		fmt.Printf("Version %d\n", version)
		fmt.Printf("Type %d\n", theType)

		switch theType {
		case 4:
			newIndex, val := applyTypeFour(binaryStr, currentIndex)
			currentIndex = newIndex
			fmt.Printf("Appending for case 4, value %d\n", val)
			subPacketValues = append(subPacketValues, val)
		default:
			newIndex, val := applyOtherTypes(theType, binaryStr, currentIndex)
			currentIndex = newIndex
			fmt.Printf("Appending for case %d, value %d\n", theType, val)
			subPacketValues = append(subPacketValues, val)
			fmt.Printf("NEW INDEX %d\n", newIndex)
		}
	}

	fmt.Printf("NOOOOOOOOOOO!!!!!!!")
	return []int{}, -1 // shouldn't happen
}

func isRemainingZero(input string, index int) bool {
	remaining := input[index : len(input)-1]
	for _, c := range strings.Split(remaining, "") {
		if c != "0" {
			return false
		}
	}
	return true
}

func applyOtherTypes(theType int, input string, index int) (int, int) {
	lengthTypeId := input[index : index+1]
	index++

	fmt.Printf("Length Type ID %s\n", lengthTypeId)

	var subpacketVals []int
	if lengthTypeId == "0" {
		// If the length type ID is 0, then the next 15 bits are a number that represents the total length in bits of the sub-packets contained by this packet.
		lengthInBits := input[index : index+15]
		index += 15
		lengthOfSubpackets := binaryToDecimal(lengthInBits)
		fmt.Printf("lengthOfSubpackets %d\n", lengthOfSubpackets)
		fmt.Printf("sub-packets %s\n", input[index:index+lengthOfSubpackets])

		// read the subpackets
		newInput := input[0 : index+lengthOfSubpackets]
		newSubs, newIndex := run(newInput, index)
		index = newIndex
		subpacketVals = newSubs
	} else {
		// If the length type ID is 1, then the next 11 bits are a number that represents the number of sub-packets immediately contained by this packet.
		numOfSubPacketsInBits := input[index : index+11]
		index += 11

		numOfSubPackets := binaryToDecimal(numOfSubPacketsInBits)
		fmt.Printf("numOfSubPackets %d\n", numOfSubPackets)
		for i := 0; i < numOfSubPackets; i++ {
			newSubs, newIndex := run(input, index)
			index = newIndex
			fmt.Printf("Appending!  %+v\n", newSubs)
			subpacketVals = append(subpacketVals, newSubs...)
		}
	}

	fmt.Printf("Running for type %d and values %+v\n", theType, subpacketVals)
	var result int
	switch theType {
	case 0:
		// Packets with type ID 0 are sum packets - their value is the sum of the values of their sub-packets. If they only have a single sub-packet, their value is the value of the sub-packet.
		result = applyTypeZero(subpacketVals)
	case 1:
		result = applyTypeOne(subpacketVals)
	case 2:
		result = applyTypeTwo(subpacketVals)
	case 3:
		result = applyTypeThree(subpacketVals)
	case 5:
		result = applyTypeFive(subpacketVals)
	case 6:
		result = applyTypeSix(subpacketVals)
	case 7:
		result = applyTypeSeven(subpacketVals)
	}

	return index, result
}

func applyTypeZero(subpacketVals []int) int {
	result := 0
	for _, v := range subpacketVals {
		result = result + v
	}
	return result
}

// Packets with type ID 1 are product packets - their value is the result of multiplying together the values of their sub-packets. If they only have a single sub-packet, their value is the value of the sub-packet.
func applyTypeOne(subpacketVals []int) int {
	result := subpacketVals[0]
	for i := 1; i < len(subpacketVals); i++ {
		result = result * subpacketVals[i]
	}
	return result
}

// Packets with type ID 2 are minimum packets - their value is the minimum of the values of their sub-packets.
func applyTypeTwo(subpacketVals []int) int {
	result := 4294967295
	for _, v := range subpacketVals {
		if v < result {
			result = v
		}
	}
	return result
}

// Packets with type ID 3 are maximum packets - their value is the maximum of the values of their sub-packets.
func applyTypeThree(subpacketVals []int) int {
	result := -1
	for _, v := range subpacketVals {
		if v > result {
			result = v
		}
	}
	return result
}

// Packets with type ID 5 are greater than packets - their value is 1 if the value of the first sub-packet is greater than the value of the second sub-packet; otherwise, their value is 0. These packets always have exactly two sub-packets.
func applyTypeFive(subpacketVals []int) int {
	if subpacketVals[0] > subpacketVals[1] {
		return 1
	}
	return 0
}

// Packets with type ID 6 are less than packets - their value is 1 if the value of the first sub-packet is less than the value of the second sub-packet; otherwise, their value is 0. These packets always have exactly two sub-packets.
func applyTypeSix(subpacketVals []int) int {
	if subpacketVals[0] < subpacketVals[1] {
		return 1
	}
	return 0
}

// Packets with type ID 7 are equal to packets - their value is 1 if the value of the first sub-packet is equal to the value of the second sub-packet; otherwise, their value is 0. These packets always have exactly two sub-packets.
func applyTypeSeven(subpacketVals []int) int {
	if subpacketVals[0] == subpacketVals[1] {
		return 1
	}
	return 0
}

func applyTypeFour(input string, index int) (int, int) {
	var binaryStr string

	for true {
		isLastGroupIndicator := input[index : index+1]
		index++ // move the pointer

		// read the next part of the number
		binaryStr += input[index : index+4]
		index += 4

		if isLastGroupIndicator == "1" {
			// not the last group, continue
		} else { // is the last group
			return index, binaryToDecimal(binaryStr)
		}
	}

	return 0, -1
}

// Every packet begins with a standard header: the first three bits encode the packet version, and the next three bits encode the packet type ID.
// These two values are numbers; all numbers encoded in any packet are represented as binary with the most significant bit first.
// For example, a version encoded as the binary sequence 100 represents the number 4.
func readHeader(input string, index int) (int, int, int) {
	versionBinary := input[index : index+3]
	index = index + 3

	if index+3 > len(input) {
		return 0, 0, 0
	}
	typeBinary := input[index : index+3]
	index = index + 3

	return index, binaryToDecimal(versionBinary), binaryToDecimal(typeBinary)
}

func binaryToDecimal(binary string) int {
	output, _ := strconv.ParseInt(binary, 2, 64)
	return int(output)
}

func convertToBinary(line string) string {
	var result string
	for _, char := range strings.Split(line, "") {
		bin := hexToBinary[char]
		result += bin
	}
	return result
}
