package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

type node struct {
	prev   *node
	next   *node
	val    string
	valInt int
}

func (n *node) isOpeningBracket() bool {
	return n.val == "["
}

func (n *node) isClosingBracket() bool {
	return n.val == "]"
}

func (n *node) isComma() bool {
	return n.val == ","
}

func (n *node) isNumber() bool {
	return !n.isOpeningBracket() && !n.isClosingBracket() && !n.isComma()
}

func (n *node) isEven() bool {
	return n.valInt%2 == 0
}

func (n *node) asString() string {
	if n.isNumber() {
		return strconv.Itoa(n.valInt)
	}
	return n.val
}

type list struct {
	head *node
	tail *node
}

func (l *list) asString() string {
	current := l.head
	result := ""
	for current.next != nil {
		result += current.asString()
		current = current.next
	}
	result += current.asString()
	return result
}

func (l *list) printList() {
	fmt.Println(l.asString())
}

// [[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]
func (l *list) findFirstNumberNestedInsideFourPairs() *node {
	var numOpeningBrackets int
	current := l.head
	for current.next != nil {
		if current.isOpeningBracket() {
			numOpeningBrackets++
		} else if current.isClosingBracket() {
			numOpeningBrackets--
		}

		if numOpeningBrackets >= 5 && current.isNumber() && current.next.isComma() && current.next.next.isNumber() {
			return current
		}
		current = current.next
	}

	return nil
}

func (l *list) findFirstNumberGtEqual10() *node {
	current := l.head
	for current.next != nil {
		if current.isNumber() && current.valInt >= 10 {
			return current
		}
		current = current.next
	}

	return nil
}

func main() {
	filepath := "input.txt"
	list := helpers.ReadFile(filepath)

	max := -1
	for i := 0; i < len(list); i++ {
		for j := 0; j < len(list); j++ {
			if i != j {
				a := NewList(list[i])
				b := NewList(list[j])
				a.addToList(b)
				a.reduceRepeatedly()
				a.calcMagnitudeUntilFinished()
				m := a.head.next.valInt
				if m > max {
					max = m
				}
			}
		}
	}
	fmt.Printf("Answer %d\n", max)
}

func NewNode(val string) *node {
	n := &node{val: val}
	if n.isNumber() {
		n.valInt, _ = strconv.Atoi(n.val)
		n.val = ""
	}
	return n
}

// [[[[4,3],4],4],[7,[[8,4],9]]]
func NewList(line string) *list {

	parts := strings.Split(line, "")
	head := NewNode(parts[0])
	prev := head
	for i := 1; i < len(parts); i++ {
		// create the new node
		n := NewNode(parts[i])
		n.prev = prev

		// hook it up
		prev.next = n
		n.prev = prev
		prev = n
	}

	return &list{head: head, tail: prev}
}

// This snailfish homework is about addition. To add two snailfish numbers, form a pair from the left and right parameters of the addition operator.
// For example, [1,2] + [[3,4],5] becomes [[1,2],[[3,4],5]].
func (l *list) addToList(otherList *list) {
	newLeftBracket := NewNode("[")
	newComma := NewNode(",")
	newRightBracket := NewNode("]")

	insertInBetween(l.tail, newComma, otherList.head)
	l.tail = otherList.tail
	l.prependToList(newLeftBracket)
	l.appendToList(newRightBracket)
}

// To reduce a snailfish number, you must repeatedly do the first action in this list that applies to the snailfish number
func (l *list) reduceRepeatedly() {
	for true {
		if !l.reduce() {
			break
		}
	}
}

/*
If any pair is nested inside four pairs, the leftmost such pair explodes.
If any regular number is 10 or greater, the leftmost such regular number splits.
Once no action in the above list applies, the snailfish number is reduced.
*/
func (l *list) reduce() bool {
	firstNested := l.findFirstNumberNestedInsideFourPairs()
	if firstNested != nil {
		l.explode(firstNested)
		return true
	}

	firstGtEq10 := l.findFirstNumberGtEqual10()
	if firstGtEq10 != nil {
		l.split(firstGtEq10)
		return true
	}

	return false
}

/*
To explode a pair, the pair's left value is added to the first regular number to the left of the exploding pair (if any), and
the pair's right value is added to the first regular number to the right of the exploding pair (if any). Exploding pairs will
always consist of two regular numbers. Then, the entire exploding pair is replaced with the regular number 0.
*/
// [[[[[9,8],1],2],3],4] becomes [[[[0,9],2],3],4] (the 9 has no regular number to its left, so it is not added to any regular number).
func (l *list) explode(n *node) {
	first := n
	second := n.next.next
	if !first.isNumber() {
		log.Fatalf("expected first to be a number, but was %+v", first)
	}
	if !second.isNumber() {
		log.Fatalf("expected first to be a number, but was %+v", first)
	}

	// the pair's left value is added to the first regular number to the left of the exploding pair (if any)
	current := first.prev
	for current != nil {
		if current.isNumber() {
			current.valInt = current.valInt + first.valInt
			break
		}
		current = current.prev
	}

	// the pair's right value is added to the first regular number to the right of the exploding pair (if any)
	current = second.next
	for current != nil {
		if current.isNumber() {
			current.valInt = current.valInt + second.valInt
			break
		}
		current = current.next
	}

	// Then, the entire exploding pair is replaced with the regular number 0.
	prevNode := first.prev.prev
	nextNode := second.next.next

	// create the new node
	newZero := NewNode("0")
	newZero.prev = prevNode
	newZero.next = nextNode

	// hook it up
	prevNode.next = newZero
	nextNode.prev = newZero
}

/*
To split a regular number, replace it with a pair; the left element of the pair should be the regular number divided by two and rounded down,
while the right element of the pair should be the regular number divided by two and rounded up.
*/
func (l *list) split(n *node) {
	if !n.isNumber() {
		log.Fatalf("expected a number for split, but was %+v", n)
	}

	// the left element of the pair should be the regular number divided by two and rounded down
	leftVal := n.valInt / 2

	// the right element of the pair should be the regular number divided by two and rounded up
	rightVal := n.valInt / 2
	if !n.isEven() {
		rightVal++
	}

	leftBracket := NewNode("[")
	leftNum := NewNode(strconv.Itoa(leftVal))
	comma := NewNode(",")
	rightNum := NewNode(strconv.Itoa(rightVal))
	rightBracket := NewNode("]")

	insertInBetween(n.prev, leftBracket, leftNum)
	insertInBetween(leftBracket, leftNum, comma)
	insertInBetween(leftNum, comma, rightNum)
	insertInBetween(comma, rightNum, rightBracket)
	insertInBetween(rightNum, rightBracket, n.next)
}

func (l *list) calcMagnitudeUntilFinished() {
	l.prependToList(NewNode("["))
	l.appendToList(NewNode("]"))
	for l.calcMagnitude() {
		// will stop once fully resolved
	}
}

// To check whether it's the right answer, the snailfish teacher only checks the magnitude of the final sum.
// The magnitude of a pair is 3 times the magnitude of its left element plus 2 times the magnitude of its right element. The magnitude of a regular number is just that number.
// For example, the magnitude of [9,1] is 3*9 + 2*1 = 29; the magnitude of [1,9] is 3*1 + 2*9 = 21. Magnitude calculations are recursive: the magnitude of [[9,1],[1,9]] is 3*29 + 2*21 = 129.
func (l *list) calcMagnitude() bool {
	current := l.head
	for current.next != nil {
		if current.isNumber() && current.next.isComma() && current.next.next.isNumber() {
			// 3 times the magnitude of its left element plus 2 times the magnitude of its right element
			m := current.valInt*3 + 2*current.next.next.valInt
			newNode := NewNode(strconv.Itoa(m))
			insertInBetween(current.prev.prev, newNode, current.next.next.next.next)
			return true
		}
		current = current.next
	}

	return false
}

func insertInBetween(a, newNode, b *node) {
	a.next = newNode
	newNode.prev = a

	newNode.next = b
	b.prev = newNode
}

func (l *list) prependToList(n *node) {
	l.head.prev = n
	n.next = l.head
	l.head = n
}

func (l *list) appendToList(n *node) {
	l.tail.next = n
	n.prev = l.tail
	l.tail = n
}
