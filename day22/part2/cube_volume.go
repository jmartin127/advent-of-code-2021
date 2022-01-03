package main

import (
	"fmt"
	"sort"
)

type interval struct {
	lower int
	upper int
}

type cube struct {
	x *interval
	y *interval
	z *interval
}

type event struct {
	val   int
	delta int
}

type subkey struct {
	x int
	y int
}

func length(inter *interval) int {
	return inter.upper - inter.lower
}

func length_of_union(intervals []*interval) int {
	events := make([]*event, 0)
	for _, inter := range intervals {
		events = append(events, &event{val: inter.lower, delta: 1})
		events = append(events, &event{val: inter.upper, delta: -1})
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].val < events[j].val
	})
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

func all_boundaries(intervals []*interval) []int {
	boundaries := make([]int, 0)
	for _, inter := range intervals {
		boundaries = append(boundaries, inter.lower)
		boundaries = append(boundaries, inter.upper)
	}
	sort.Ints(boundaries)
	return boundaries
}

func subdivide_at(inter *interval, boundaries []int) []*interval {
	result := make([]*interval, 0)

	lower := inter.lower
	sort.Ints(boundaries)
	for _, x := range boundaries { // Resort is O(n) due to Timsort.
		if x < lower {
			continue // TODO: Is this the same as "pass" in Python?
		} else if x < inter.upper {
			result = append(result, &interval{lower: lower, upper: x})
			lower = x
		} else {
			result = append(result, &interval{lower: lower, upper: inter.upper})
			break
		}
	}
	return result
}

type xyInterval struct {
	x *interval
	y *interval
}

func volume_of_union(cubes []*cube) int {
	x_intervals := make([]*interval, 0)
	y_intervals := make([]*interval, 0)
	for _, c := range cubes {
		x_intervals = append(x_intervals, c.x)
		y_intervals = append(y_intervals, c.y)
	}
	x_boundaries := all_boundaries(x_intervals)
	y_boundaries := all_boundaries(y_intervals)

	sub_problems := make(map[string][]*interval, 0) // should be key as "x,y", value as []*interval
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

// func keyToInts(key string) []int {
// 	parts := strings.Split(key, ",")
// 	first, _ := strconv.Atoi(parts[0])
// 	second, _ := strconv.Atoi(parts[1])
// 	return []int{first, second}
// }

func appendToMap(sub_problems map[string][]*interval, key string, inter *interval) {
	if val, ok := sub_problems[key]; ok {
		sub_problems[key] = append(val, inter)
	} else {
		sub_problems[key] = []*interval{inter}
	}
}
