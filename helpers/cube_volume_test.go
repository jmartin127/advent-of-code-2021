package helpers

import "testing"

func TestVolumeOfUnion(t *testing.T) {
	cubes := []*Cube{
		{
			x: &Interval{-2, 2},
			y: &Interval{-2, 2},
			z: &Interval{-2, 2},
		},
		{
			x: &Interval{1, 5},
			y: &Interval{1, 2},
			z: &Interval{1, 2},
		},
	}
	volume := Volume_of_union(cubes)
	if volume != 0 {
		t.Fatalf("Expected %d, but was %d\n", 0, volume)
	}
}

func TestVolumeOfUnionRealWorld(t *testing.T) {
	//on x=10..11,y=12..12,z=12..12
	//on x=12..12,y=12..12,z=12..12
	cubes := []*Cube{
		{
			x: &Interval{0, 2},
			y: &Interval{0, 2},
			z: &Interval{0, 2},
		},
		{
			x: &Interval{10, 12},
			y: &Interval{12, 13},
			z: &Interval{12, 13},
		},
		// {
		// 	x: &Interval{12, 13},
		// 	y: &Interval{12, 13},
		// 	z: &Interval{12, 13},
		// },
	}
	volume := Volume_of_union(cubes)
	if volume != 13 {
		t.Fatalf("Expected %d, but was %d\n", 0, volume)
	}
}
