package main

import (
	"fmt"
	"strings"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

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
	moveUntilFinished(c)
}

// The amphipods would like a method to organize every amphipod into side rooms so that each side room contains
// one type of amphipod and the types are sorted A-D going left to right, like this:
func moveUntilFinished(inputCave *cave) { // all reached destination, all stuck
	c := inputCave.copy()

	c.print()

	// base cases:
	// a) all amphipods have found thier final spot
	if c.allFoundRoom() {
		fmt.Printf("ALL HAVE FOUND A ROOM!")
		c.print()
		return
	}

	for _, a := range c.amphipods {
		// try all of the hall positions
		for hallAssignment := 0; hallAssignment < 7; hallAssignment++ {
			moved := moveIntoHall(c, a, hallAssignment)
			if moved { // If we weren't able to move the amphipod, don't recurse
				moveUntilFinished(c)
			}
		}

		// try going to the room
		moved := moveIntoRoom(c, a)
		if moved { // If we weren't able to move the amphipod, don't recurse
			moveUntilFinished(c)
		}
	}
}

// TODO don't allow a crab to move out if it is the only one in the room, and it is already in the correc position
func moveIntoHall(c *cave, a *amphipod, hallAssignment int) bool {
	// Once an amphipod stops moving in the hallway, it will stay in that spot until it can move into a room.
	if a.isInHall {
		return false
	}

	// Once an amphipod has reached it's final destination, it shouldn't come back out
	if a.hasFoundRoom {
		return false
	}

	// Determine where we are going
	destinationX := HALL_LOCATIONS[hallAssignment]

	// Check if the hallway is clear
	if !isHallwayClear(c, a, destinationX) {
		return false
	}

	// move to that position in the hallway
	c.cells[a.yPos][a.xPos].amph = nil
	c.cells[1][destinationX].amph = a
	a.xPos = destinationX
	a.yPos = 1
	a.isInHall = true
	return true
}

func moveIntoRoom(c *cave, a *amphipod) bool {
	if !a.isInHall { // already in a room
		return false
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
			return false // can't move into the room, there is a amphipod of the wrong type already in the room
		} else {
			// there is already one in the furtherst spot, check if the last spot is open
			lastSpot := ROOM_ASSIGNMENTS[a.aType][1]
			lastCell := c.cells[lastSpot.y][lastSpot.x]
			if lastCell.amph != nil {
				return false // can't move here, there is already an amphipod in the last spot
			} else {
				// the last spot is open, let's try to go there
				xTargetPos = lastSpot.x
				yTargetPos = lastSpot.y
			}
		}
	}

	// move to the target position (if nothing is in the way)
	if !isHallwayClear(c, a, xTargetPos) {
		return false
	}

	// since the hallway is clear, and we've already checked the room, let's move IN!
	c.cells[a.yPos][a.xPos].amph = nil
	c.cells[yTargetPos][xTargetPos].amph = a
	a.xPos = xTargetPos
	a.yPos = yTargetPos
	a.isInHall = false
	a.hasFoundRoom = true
	return true
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
