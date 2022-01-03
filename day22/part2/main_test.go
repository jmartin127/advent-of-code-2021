package main

// import (
// 	"testing"
// )

// func TestDivideOnCubeUsingOverlappingOffCube_Simple(t *testing.T) {
// 	onCube := &instruction{isOn: true, xStart: 0, xEnd: 2, yStart: 0, yEnd: 2, zStart: 0, zEnd: 2}
// 	offCube := &instruction{isOn: false, xStart: 1, xEnd: 1, yStart: 1, yEnd: 1, zStart: 1, zEnd: 1}
// 	result := divideOnCubeUsingOverlappingOffCube(onCube, offCube)
// 	if len(result) != 26 {
// 		t.Fatalf("Expected 26 found %d\n", len(result))
// 	}
// 	for _, r := range result {
// 		if r.volume() != 1 {
// 			t.Fatalf("Expected volume to be 1, was %d\n", r.volume())
// 		}
// 	}
// }

// func TestDivideOnCubeUsingOverlappingOffCube_RealWorld(t *testing.T) {
// 	onCube := &instruction{isOn: true, xStart: 10, xEnd: 12, yStart: 10, yEnd: 12, zStart: 10, zEnd: 12}
// 	offCube := &instruction{isOn: false, xStart: 9, xEnd: 11, yStart: 9, yEnd: 11, zStart: 9, zEnd: 11}
// 	result := divideOnCubeUsingOverlappingOffCube(onCube, offCube)

// 	// validate that the volume of the leftover after Off = Total ON - overlap with off
// 	overlap := findSharedCubeBetweenTwoCuboids(onCube, offCube)
// 	leftoverAfterOffVolume := totalVolume(result)
// 	expectedVolume := onCube.volume() - overlap.volume()
// 	if leftoverAfterOffVolume != expectedVolume {
// 		t.Fatalf("Expected volume %d, but was %d", expectedVolume, leftoverAfterOffVolume)
// 	}

// 	if len(result) != 7 {
// 		t.Fatalf("Expected 7 found %d\n", len(result))
// 	}
// }

// func TestDivideOnCubeUsingOverlappingOffCube_RealWorld2(t *testing.T) {
// 	onCube := &instruction{isOn: true, xStart: 11, xEnd: 13, yStart: 11, yEnd: 13, zStart: 11, zEnd: 13}
// 	offCube := &instruction{isOn: false, xStart: 9, xEnd: 11, yStart: 9, yEnd: 11, zStart: 9, zEnd: 11}
// 	result := divideOnCubeUsingOverlappingOffCube(onCube, offCube)

// 	// validate that the volume of the leftover after Off = Total ON - overlap with off
// 	overlap := findSharedCubeBetweenTwoCuboids(onCube, offCube)
// 	leftoverAfterOffVolume := totalVolume(result)
// 	expectedVolume := onCube.volume() - overlap.volume()
// 	if leftoverAfterOffVolume != expectedVolume {
// 		t.Fatalf("Expected volume %d, but was %d", expectedVolume, leftoverAfterOffVolume)
// 	}

// 	if len(result) != 7 {
// 		t.Fatalf("Expected 7 found %d\n", len(result))
// 	}
// }

// func totalVolume(instructions []*instruction) int {
// 	var total int
// 	for _, i := range instructions {
// 		total += i.volume()
// 	}
// 	return total
// }
