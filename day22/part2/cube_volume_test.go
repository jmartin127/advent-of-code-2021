package main

import "testing"

func TestVolumeOfUnion(t *testing.T) {
	cubes := []*cube{
		{
			x: &interval{-2, 2},
			y: &interval{-2, 2},
			z: &interval{-2, 2},
		},
		{
			x: &interval{1, 5},
			y: &interval{1, 2},
			z: &interval{1, 2},
		},
	}
	volume := volume_of_union(cubes)
	if volume != 0 {
		t.Fatalf("Expected %d, but was %d\n", 0, volume)
	}
}
