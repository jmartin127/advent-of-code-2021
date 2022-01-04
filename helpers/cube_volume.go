package helpers

// NOTE: there is still a bug in here that I haven't quite worked out. This is a port from the Python code here:
// https://stackoverflow.com/questions/69137352/computing-the-volume-of-the-union-of-axis-aligned-cubes

import (
	"fmt"
	"sort"
)

type Interval struct {
	lower int
	upper int
}

type Cube struct {
	x *Interval
	y *Interval
	z *Interval
}

func (c *Cube) AsString() string {
	return fmt.Sprintf("x=%d..%d y=%d..%d z=%d..%d", c.x.lower, c.x.upper, c.y.lower, c.y.upper, c.z.lower, c.z.upper)
}

func NewCube(x *Interval, y *Interval, z *Interval) *Cube {
	return &Cube{x: x, y: y, z: z}
}

func NewInterval(lower, upper int) *Interval {
	return &Interval{lower: lower, upper: upper}
}

type event struct {
	val   int
	delta int
}

type subkey struct {
	x int
	y int
}

func length(inter *Interval) int {
	return inter.upper - inter.lower
}

func length_of_union(intervals []*Interval) int {
	events := make([]*event, 0)
	for _, inter := range intervals {
		events = append(events, &event{val: inter.lower, delta: 1})
		events = append(events, &event{val: inter.upper, delta: -1})
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].val < events[j].val
	})
	// fmt.Println("Sorted")
	// for _, e := range events {
	// 	fmt.Printf("Event %+v\n", e)
	// }
	previous := 0
	overlap := 0
	total := 0
	for _, e := range events {
		x := e.val
		delta := e.delta
		if overlap > 0 {
			total += x - previous
		}
		previous = x
		overlap += delta
	}
	return total
}

func all_boundaries(intervals []*Interval) []int {
	boundaries := make([]int, 0)
	for _, inter := range intervals {
		boundaries = append(boundaries, inter.lower)
		boundaries = append(boundaries, inter.upper)
	}
	sort.Ints(boundaries)
	return boundaries
}

func subdivide_at(inter *Interval, boundaries []int) []*Interval {
	result := make([]*Interval, 0)

	lower := inter.lower
	sort.Ints(boundaries)
	for _, x := range boundaries { // Resort is O(n) due to Timsort.
		if x < lower {
			// nothing
		} else if x < inter.upper {
			result = append(result, &Interval{lower: lower, upper: x})
			lower = x
		} else {
			result = append(result, &Interval{lower: lower, upper: inter.upper})
			break
		}
	}
	return result
}

type xyInterval struct {
	x *Interval
	y *Interval
}

func Volume_of_union(cubes []*Cube) int {
	x_intervals := make([]*Interval, 0)
	y_intervals := make([]*Interval, 0)
	for _, c := range cubes {
		x_intervals = append(x_intervals, c.x)
		y_intervals = append(y_intervals, c.y)
	}
	x_boundaries := all_boundaries(x_intervals)
	y_boundaries := all_boundaries(y_intervals)

	sub_problems := make(map[string][]*Interval, 0) // should be key as "x,y", value as []*interval
	sub_problems_lookup := make(map[string]*xyInterval)
	for _, cube := range cubes {
		for x, xInt := range subdivide_at(cube.x, x_boundaries) {
			for y, yInt := range subdivide_at(cube.y, y_boundaries) {
				key := fmt.Sprintf("%d,%d", x, y)
				sub_problems_lookup[key] = &xyInterval{x: xInt, y: yInt}
				appendToMap(sub_problems, key, cube.z)
			}
		}
	}

	var result int
	for k, z_intervals := range sub_problems {
		xyInt := sub_problems_lookup[k]
		result += length(xyInt.x) * length(xyInt.y) * length_of_union(z_intervals)
	}

	return result
}

func appendToMap(sub_problems map[string][]*Interval, key string, inter *Interval) {
	if val, ok := sub_problems[key]; ok {
		sub_problems[key] = append(val, inter)
	} else {
		sub_problems[key] = []*Interval{inter}
	}
}
