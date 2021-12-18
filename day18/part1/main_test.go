package main

import (
	"strconv"
	"testing"
)

func TestExplode(t *testing.T) {
	l := NewList("[[[[[9,8],1],2],3],4]")
	first := l.findFirstNumberNestedInsideFourPairs()
	l.explode(first)
	if l.asString() != "[[[[0,9],2],3],4]" {
		t.Fatalf("(1) result was %s\n", l.asString())
	}

	l = NewList("[7,[6,[5,[4,[3,2]]]]]")
	first = l.findFirstNumberNestedInsideFourPairs()
	l.explode(first)
	if l.asString() != "[7,[6,[5,[7,0]]]]" {
		t.Fatalf("(1) result was %s\n", l.asString())
	}

	l = NewList("[[6,[5,[4,[3,2]]]],1]")
	first = l.findFirstNumberNestedInsideFourPairs()
	l.explode(first)
	if l.asString() != "[[6,[5,[7,0]]],3]" {
		t.Fatalf("(1) result was %s\n", l.asString())
	}

	l = NewList("[[3,[2,[1,[7,3]]]],[6,[5,[4,[3,2]]]]]")
	first = l.findFirstNumberNestedInsideFourPairs()
	l.explode(first)
	if l.asString() != "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]" {
		t.Fatalf("(1) result was %s\n", l.asString())
	}

	l = NewList("[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]")
	first = l.findFirstNumberNestedInsideFourPairs()
	l.explode(first)
	if l.asString() != "[[3,[2,[8,0]]],[9,[5,[7,0]]]]" {
		t.Fatalf("(1) result was %s\n", l.asString())
	}
}

// For example, 10 becomes [5,5], 11 becomes [5,6], 12 becomes [6,6], and so on.
func TestSplit(t *testing.T) {
	l := newListFromNum(10)
	l.split(l.head.next)
	if l.asString() != "[[5,5]]" {
		t.Fatalf("(1) result was %s\n", l.asString())
	}

	l = newListFromNum(11)
	l.split(l.head.next)
	if l.asString() != "[[5,6]]" {
		t.Fatalf("(1) result was %s\n", l.asString())
	}

	l = newListFromNum(12)
	l.split(l.head.next)
	if l.asString() != "[[6,6]]" {
		t.Fatalf("(1) result was %s\n", l.asString())
	}
}

func TestReduceRepeatedly(t *testing.T) {
	l := NewList("[[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]")
	l.reduceRepeatedly()
	if l.asString() != "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]" {
		t.Fatalf("(1) result was %s\n", l.asString())
	}
}

// For example, [1,2] + [[3,4],5] becomes [[1,2],[[3,4],5]].
func TestAddToList(t *testing.T) {
	a := NewList("[1,2]")
	b := NewList("[[3,4],5]")

	a.addToList(b)
	if a.asString() != "[[1,2],[[3,4],5]]" {
		t.Fatalf("(1) result was %s\n", a.asString())
	}
}

func TestMagnitude(t *testing.T) {
	l := NewList("[9,1]")
	l.calcMagnitudeUntilFinished()
	if l.asString() != "[29]" {
		t.Fatalf("(1) result was %s\n", l.asString())
	}

	l = NewList("[1,9]")
	l.calcMagnitudeUntilFinished()
	if l.asString() != "[21]" {
		t.Fatalf("(1) result was %s\n", l.asString())
	}

	l = NewList("[[9,1],[1,9]]")
	l.calcMagnitudeUntilFinished()
	if l.asString() != "[129]" {
		t.Fatalf("(1) result was %s\n", l.asString())
	}

	l = NewList("[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]")
	l.calcMagnitudeUntilFinished()
	if l.asString() != "[3488]" {
		t.Fatalf("(1) result was %s\n", l.asString())
	}
}

func newListFromNum(num int) *list {
	head := NewNode("[")
	val := NewNode(strconv.Itoa(num))
	tail := NewNode("]")

	insertInBetween(head, val, tail)

	return &list{head: head}
}
