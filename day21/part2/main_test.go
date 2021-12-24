package main

import "testing"

func TestDetermineSpace(t *testing.T) {
	if determineSpace(3, 3) != 6 {
		t.Fatal("expected 6")
	}
	if determineSpace(9, 3) != 2 {
		t.Fatal("expected 2")
	}
	if determineSpace(10, 1) != 1 {
		t.Fatal("expected 1")
	}
	if determineSpace(9, 2) != 1 {
		t.Fatal("expected 1")
	}
	if determineSpace(9, 3) != 2 {
		t.Fatal("expected 2")
	}
}
