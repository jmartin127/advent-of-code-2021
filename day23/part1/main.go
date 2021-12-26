package main

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

var totalPrint = 0
var minEnergy = 1000000 // TODO

const (
	WALL    = "#"
	HALL    = "."
	AMBER   = "A"
	BRONZE  = "B"
	COPPER  = "C"
	DESERT  = "D"
	NOTHING = " "
	ROOM    = "R"
)

// ACTIONS
const (
	TO_HALL = "th"
	TO_ROOM = "tr"
)

// HALL LOCATIONS (key = hall ID, value = xCoord)
/*
#############
#01.2.3.4.56#
###B#C#B#D###
  #A#D#C#A#
  #########
*/
var HALL_LOCATIONS = map[int]int{
	0: 1,
	1: 2,
	2: 4,
	3: 6,
	4: 8,
	5: 10,
	6: 11,
}

var perms = []*permutation{}

type permutation struct {
	amphipodIndex int
	goToHall      bool // true for hall, false for room
	hallIndex     int
}

type move struct {
	*permutation
	distanceMoved int
	energy        int
	aStartX       int
	aStartY       int
}

type coord struct {
	x int
	y int
}

// HALL LOCATIONS (key = type of amphipod, value = yCoord,xCoord)
/*
#############
#...........#
###A#B#C#D###
  #A#B#C#D#
  #########
*/
var ROOM_ASSIGNMENTS = map[string][]*coord{
	"A": []*coord{{x: 3, y: 3}, {x: 3, y: 2}},
	"B": []*coord{{x: 5, y: 3}, {x: 5, y: 2}},
	"C": []*coord{{x: 7, y: 3}, {x: 7, y: 2}},
	"D": []*coord{{x: 9, y: 3}, {x: 9, y: 2}},
}

// Amber amphipods require 1 energy per step, Bronze amphipods require 10 energy, Copper
// amphipods require 100, and Desert ones require 1000.
var ENGERY_PER_MOVE = map[string]int{
	"A": 1,
	"B": 10,
	"C": 100,
	"D": 1000,
}

type amphipod struct {
	aType        string
	xPos         int
	yPos         int
	hasFoundRoom bool
	isInHall     bool
}

func (a *amphipod) copy() *amphipod {
	return &amphipod{
		aType:        a.aType,
		xPos:         a.xPos,
		yPos:         a.yPos,
		hasFoundRoom: a.hasFoundRoom,
		isInHall:     a.isInHall,
	}
}

type cell struct {
	cType string // WALL, HALL, ROOM, NOTHING
	amph  *amphipod
}

func (c *cell) copy() *cell {
	var a *amphipod
	if c.amph != nil {
		a = c.amph.copy()
	}
	return &cell{
		cType: c.cType,
		amph:  a,
	}
}

type cave struct {
	cells     [][]*cell
	amphipods []*amphipod
}

func (c *cave) copy() *cave {
	cells := make([][]*cell, 0)
	for _, row := range c.cells {
		newRow := make([]*cell, 0)
		for _, v := range row {
			newRow = append(newRow, v.copy())
		}
		cells = append(cells, newRow)
	}

	amphipods := make([]*amphipod, 0)
	for _, a := range c.amphipods {
		amphipods = append(amphipods, a.copy())
	}

	return &cave{
		cells:     cells,
		amphipods: amphipods,
	}
}

func (c *cave) allFoundRoom() bool {
	for _, a := range c.amphipods {
		if !a.hasFoundRoom {
			return false
		}
	}
	return true
}

func (c *cave) print() {
	for _, row := range c.cells {
		for _, v := range row {
			val := v.cType
			if v.cType == ROOM || v.cType == HALL {
				if v.amph != nil {
					val = v.amph.aType
				} else {
					val = "."
				}
			}
			fmt.Printf("%s", val)
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	c := parseInput()
	c.print()

	perms = createPermutations(c.amphipods)
	fmt.Printf("Num of permutations %d\n", len(perms))

	moveUntilFinished(c, []*move{}, 0)
}

// Should be 8 * 8 = 64 permutations
// 8 amphipods
// 8 possible moves for amphipods (7 ways to enter hall, 1 way to enter room)
func createPermutations(amphipods []*amphipod) []*permutation {
	result := make([]*permutation, 0)
	for amphipodIndex := range amphipods {
		// hall options
		for hallIndex := range HALL_LOCATIONS {
			p := &permutation{amphipodIndex: amphipodIndex, goToHall: true, hallIndex: hallIndex}
			result = append(result, p)
		}

		// room option
		p := &permutation{amphipodIndex: amphipodIndex, goToHall: false, hallIndex: -1}
		result = append(result, p)
	}
	return result
}

// The amphipods would like a method to organize every amphipod into side rooms so that each side room contains
// one type of amphipod and the types are sorted A-D going left to right, like this:
func moveUntilFinished(c *cave, path []*move, energyToThisPoint int) {
	if energyToThisPoint > minEnergy {
		return
	}

	// base case:
	if c.allFoundRoom() {
		fmt.Printf("FOUND A ROOM!")
		//c.print()
		//fmt.Println("Path... (made up of moves)")
		var totalEnergy int
		for _, m := range path {
			// fmt.Printf("\tmove energy %d, perm %t %d, a=%s, dist=%d, aStartX=%d, aStartY=%d\n", m.energy, m.permutation.goToHall, m.permutation.hallIndex, c.amphipods[m.amphipodIndex].aType, m.distanceMoved, m.aStartX, m.aStartY)
			totalEnergy += m.energy
		}
		if totalEnergy < minEnergy {
			minEnergy = totalEnergy
		}
		fmt.Printf("\tTOTAL ENERGY %d. Min energy %d\n", totalEnergy, minEnergy)
		return
	}

	for _, p := range perms {
		applyPermutation(c.copy(), p, path, energyToThisPoint)
	}
}

func applyPermutation(c *cave, p *permutation, path []*move, energyToThisPoint int) {
	var moved bool
	var m *move
	a := c.amphipods[p.amphipodIndex]
	if p.goToHall {
		moved, m = moveIntoHall(c, a, p.hallIndex)
	} else {
		moved, m = moveIntoRoom(c, a)
	}
	if moved {
		energy := m.distanceMoved * ENGERY_PER_MOVE[a.aType]
		m.energy = energy
		m.permutation = p
		path = append(path, m)
		moveUntilFinished(c, path, energyToThisPoint+energy)
	}
}

func moveIntoHall(c *cave, a *amphipod, hallAssignment int) (bool, *move) {
	// make sure the cave thinks it is here as well
	if c.cells[a.yPos][a.xPos].amph == nil {
		log.Fatal("The cave doesn't think anything is here")
	}

	// Once an amphipod stops moving in the hallway, it will stay in that spot until it can move into a room.
	if a.isInHall {
		return false, nil
	}

	// Once an amphipod has reached it's final destination, it shouldn't come back out
	if a.hasFoundRoom {
		return false, nil
	}

	// Determine where we are going
	destinationX := HALL_LOCATIONS[hallAssignment]

	// If it is on the bottom, check if there is another amphipod blocking
	if a.yPos == 3 {
		if c.cells[a.yPos-1][a.xPos].amph != nil {
			return false, nil
		}
	}

	// Check if the hallway is clear
	if !isHallwayClear(c, a, destinationX) {
		return false, nil
	}

	// cost to go right/left in the hall
	distanceMoved := int(math.Abs(float64(a.xPos-destinationX))) + int(math.Abs(float64(a.yPos-1)))

	// record keeping
	aStartX := a.xPos
	aStartY := a.yPos

	// move to that position in the hallway
	// make sure the cave thinks it is here as well
	c.cells[a.yPos][a.xPos].amph = nil
	c.cells[1][destinationX].amph = a
	a.xPos = destinationX
	a.yPos = 1
	a.isInHall = true
	return true, &move{distanceMoved: distanceMoved, aStartX: aStartX, aStartY: aStartY}
}

func moveIntoRoom(c *cave, a *amphipod) (bool, *move) {
	if !a.isInHall { // already in a room
		return false, nil
	}

	// First check if the furtherst spot in the room is available, if it isn't, take the last spot
	furthestSpot := ROOM_ASSIGNMENTS[a.aType][0]
	furthestCell := c.cells[furthestSpot.y][furthestSpot.x]
	var xTargetPos int
	var yTargetPos int
	if furthestCell.amph == nil {
		xTargetPos = furthestSpot.x
		yTargetPos = furthestSpot.y
	} else {
		// if there is an amphipod in the preferred spot, see if it is the righ type
		if furthestCell.amph.aType != a.aType {
			return false, nil // can't move into the room, there is a amphipod of the wrong type already in the room
		} else {
			// there is already one in the furtherst spot, check if the last spot is open
			lastSpot := ROOM_ASSIGNMENTS[a.aType][1]
			lastCell := c.cells[lastSpot.y][lastSpot.x]
			if lastCell.amph != nil {
				return false, nil // can't move here, there is already an amphipod in the last spot
			} else {
				// the last spot is open, let's try to go there
				xTargetPos = lastSpot.x
				yTargetPos = lastSpot.y
			}
		}
	}

	// move to the target position (if nothing is in the way)
	if !isHallwayClear(c, a, xTargetPos) {
		return false, nil
	}

	// calculate the distance
	distanceMoved := int(math.Abs(float64(a.xPos-xTargetPos))) + int(math.Abs(float64(a.yPos-yTargetPos)))

	// record keeping
	aStartX := a.xPos
	aStartY := a.yPos

	// since the hallway is clear, and we've already checked the room, let's move IN!
	c.cells[a.yPos][a.xPos].amph = nil
	c.cells[yTargetPos][xTargetPos].amph = a
	a.xPos = xTargetPos
	a.yPos = yTargetPos
	a.isInHall = false
	a.hasFoundRoom = true
	return true, &move{distanceMoved: distanceMoved, aStartX: aStartX, aStartY: aStartY}
}

func isHallwayClear(c *cave, a *amphipod, destinationX int) bool {
	startingX := a.xPos

	// Check if the hallway is clear
	if destinationX < startingX { // move left
		for i := startingX - 1; i >= destinationX; i-- {
			if c.cells[1][i].amph != nil {
				return false // can't move here (blocked)
			}
		}
	} else { // move right
		for i := startingX + 1; i <= destinationX; i++ {
			if c.cells[1][i].amph != nil {
				return false // can't move here (blocked)
			}
		}
	}
	return true
}

/*
#############
#...........#
###B#C#B#D###
  #A#D#C#A#
  #########
*/
func parseInput() *cave {
	filepath := "day23/input.txt"
	list := helpers.ReadFile(filepath)

	cells := make([][]*cell, 0)
	amphipods := make([]*amphipod, 0)
	for i, line := range list {
		vals := strings.Split(line, "")
		row := make([]*cell, 0)
		for j, v := range vals {
			var cType string
			var amph *amphipod
			switch v {
			case "#":
				cType = WALL
			case ".":
				cType = HALL
			case " ":
				cType = NOTHING
			default:
				cType = ROOM
				amph = &amphipod{
					aType: v,
					yPos:  i,
					xPos:  j,
				}
				firstCavePos := ROOM_ASSIGNMENTS[amph.aType][0]
				if amph.xPos == firstCavePos.x && amph.yPos == firstCavePos.y {
					amph.hasFoundRoom = true
					fmt.Printf("Already in the right spot %s\n", amph.aType)
				}
				amphipods = append(amphipods, amph)
			}
			c := &cell{
				cType: cType,
				amph:  amph,
			}
			row = append(row, c)
		}
		cells = append(cells, row)
	}

	return &cave{cells: cells, amphipods: amphipods}
}
