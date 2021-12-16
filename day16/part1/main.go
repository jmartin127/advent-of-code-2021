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

	versionSum, _ := run(binaryStr, 0)
	fmt.Printf("version sum %d\n", versionSum)
}

func run(binaryStr string, startIndex int) (int, int) {
	versionSum := 0

	currentIndex := startIndex
	for true {
		if currentIndex >= len(binaryStr)-1 || isRemainingZero(binaryStr, currentIndex) {
			return versionSum, currentIndex
		}
		newIndex, version, theType := readHeader(binaryStr, currentIndex)
		currentIndex = newIndex
		fmt.Println("******************************")
		fmt.Printf("Version %d\n", version)
		fmt.Printf("Type %d\n", theType)

		versionSum += version

		switch theType {
		case 4:
			newIndex, _ := applyTypeFour(binaryStr, currentIndex)
			currentIndex = newIndex
		default:
			newIndex, additionalVersion := applyOtherTypes(binaryStr, currentIndex)
			currentIndex = newIndex
			versionSum += additionalVersion
			fmt.Printf("NEW INDEX %d\n", newIndex)
		}
	}

	fmt.Printf("NOOOOOOOOOOO!!!!!!!")
	return -1, -1 // shouldn't happen
}

func isRemainingZero(input string, index int) bool {
	fmt.Printf("Len %d\n", len(input))
	remaining := input[index : len(input)-1]
	fmt.Printf("Remaining %s\n", remaining)
	for _, c := range strings.Split(remaining, "") {
		if c != "0" {
			return false
		}
	}
	return true
}

func applyOtherTypes(input string, index int) (int, int) {
	lengthTypeId := input[index : index+1]
	index++

	fmt.Printf("Length Type ID %s\n", lengthTypeId)

	totalVersion := 0
	if lengthTypeId == "0" {
		// If the length type ID is 0, then the next 15 bits are a number that represents the total length in bits of the sub-packets contained by this packet.
		lengthInBits := input[index : index+15]
		index += 15
		lengthOfSubpackets := binaryToDecimal(lengthInBits)
		fmt.Printf("lengthOfSubpackets %d\n", lengthOfSubpackets)
		fmt.Printf("sub-packets %s\n", input[index:index+lengthOfSubpackets])

		// read the subpackets
		newInput := input[0 : index+lengthOfSubpackets]
		additionalVersion, newIndex := run(newInput, index)
		index = newIndex
		totalVersion += additionalVersion
	} else {
		// If the length type ID is 1, then the next 11 bits are a number that represents the number of sub-packets immediately contained by this packet.
		numOfSubPacketsInBits := input[index : index+11]
		index += 11

		numOfSubPackets := binaryToDecimal(numOfSubPacketsInBits)
		fmt.Printf("numOfSubPackets %d\n", numOfSubPackets)
		for i := 0; i < numOfSubPackets; i++ {
			additionalVersion, newIndex := run(input, index)
			index = newIndex
			totalVersion += additionalVersion
		}
	}

	return index, totalVersion
}

func applyTypeFour(input string, index int) (int, string) {
	var binaryStr string

	for true {
		isLastGroupIndicator := input[index : index+1]
		index++ // move the pointer

		// read the next part of the number
		binaryStr += input[index : index+4]
		index += 4

		fmt.Printf("found %s\n", binaryStr)

		if isLastGroupIndicator == "1" {
			fmt.Printf("not last\n")
			// not the last group, continue
		} else { // is the last group
			fmt.Printf("hit last\n")
			return index, binaryStr
		}
	}

	return 0, ""
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
