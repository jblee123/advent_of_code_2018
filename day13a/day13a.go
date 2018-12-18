package main

import (
	"fmt"
	"sort"

	"jblee.net/adventofcode2018/utils"
)

type turnDir byte

const (
	straightTurn turnDir = iota
	leftTurn     turnDir = iota
	rightTurn    turnDir = iota
)

var leftTurns = map[byte]byte{
	'^': '<',
	'>': '^',
	'v': '>',
	'<': 'v',
}

var rightTurns = map[byte]byte{
	'^': '>',
	'>': 'v',
	'v': '<',
	'<': '^',
}

var bendTurns = map[string]byte{
	`^\`: '<',
	`^/`: '>',
	`>\`: 'v',
	`>/`: '^',
	`v\`: '>',
	`v/`: '<',
	`<\`: '^',
	`</`: 'v',
}

var nextTurns = map[turnDir]turnDir{
	leftTurn:     straightTurn,
	straightTurn: rightTurn,
	rightTurn:    leftTurn,
}

type mapCoord struct {
	x, y int
}

type movementOffset struct {
	dx, dy int
}

type cart struct {
	pos       mapCoord
	direction byte
	nextTurn  turnDir
}

type mapSpace struct {
	track byte
}

type mineMap [][]mapSpace

var movementOffsets = map[byte]movementOffset{
	'^': movementOffset{0, -1},
	'>': movementOffset{1, 0},
	'v': movementOffset{0, 1},
	'<': movementOffset{-1, 0},
}

type byCartPos [](*cart)

func (a byCartPos) Len() int {
	return len(a)
}
func (a byCartPos) Less(i, j int) bool {
	x1, y1 := a[i].pos.x, a[i].pos.y
	x2, y2 := a[j].pos.x, a[j].pos.y
	return (x1 < x2) || ((x1 == x2) && (y1 < y2))
}
func (a byCartPos) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func inputCharToTrackChar(c byte) byte {
	if c == '<' || c == '>' {
		return '-'
	} else if c == '^' || c == 'v' {
		return '|'
	} else {
		return c
	}
}

func inputCharToCart(c byte) *cart {
	var newCart *cart
	if c == '<' || c == '>' || c == '^' || c == 'v' {
		newCart = new(cart)
		newCart.direction = c
		newCart.nextTurn = leftTurn
	}

	return newCart
}

func parseMapInput(lines []string) (mineMap, [](*cart)) {
	xLen := len(lines[0])
	yLen := len(lines)

	theMap := make(mineMap, xLen)
	carts := make([](*cart), 0)
	for x := 0; x < xLen; x++ {
		theMap[x] = make([]mapSpace, yLen)
		for y := 0; y < yLen; y++ {
			theMap[x][y].track = inputCharToTrackChar(lines[y][x])
			newCart := inputCharToCart(lines[y][x])
			if newCart != nil {
				newCart.pos = mapCoord{x, y}
				carts = append(carts, newCart)
			}
		}
	}

	sort.Sort(byCartPos(carts))

	return theMap, carts
}

func printMap(theMap mineMap) {
	xLen := len(theMap)
	yLen := len(theMap[0])
	strLen := (xLen + 1) * yLen
	byteBuf := make([]byte, strLen)
	bufIdx := 0
	for y := 0; y < yLen; y++ {
		for x := 0; x < xLen; x++ {
			byteBuf[bufIdx] = theMap[x][y].track
			bufIdx++
		}
		byteBuf[bufIdx] = '\n'
		bufIdx++
	}

	fmt.Print(string(byteBuf))
}

func applyTurn(theCart *cart) {
	if theCart.nextTurn == leftTurn {
		theCart.direction, _ = leftTurns[theCart.direction]
	} else if theCart.nextTurn == rightTurn {
		theCart.direction, _ = rightTurns[theCart.direction]
	} else if theCart.nextTurn == straightTurn {
		// keep same direction
	} else {
		panic(fmt.Sprintf("bad cart nextTurn: %d", theCart.nextTurn))
	}

	theCart.nextTurn, _ = nextTurns[theCart.nextTurn]
}

func moveCart(theMap mineMap, theCart *cart) {

	mapSpot := &theMap[theCart.pos.x][theCart.pos.y]

	// do any turns
	if mapSpot.track == '+' {
		applyTurn(theCart)
	} else if mapSpot.track == '/' || mapSpot.track == '\\' {
		bendKey := string([]byte{theCart.direction, mapSpot.track})
		theCart.direction, _ = bendTurns[bendKey]
	}

	offset := movementOffsets[theCart.direction]
	theCart.pos.x += offset.dx
	theCart.pos.y += offset.dy
}

func checkForCrash(carts [](*cart)) (bool, mapCoord) {
	for i := 0; i < len(carts)-1; i++ {
		if carts[i].pos == carts[i+1].pos {
			return true, carts[i].pos
		}
	}

	return false, mapCoord{}
}

func executeStep(theMap mineMap, carts [](*cart)) (
	crashed bool, crashSite mapCoord) {

	crashed = false

	origOrderedCarts := make([](*cart), len(carts))
	copy(origOrderedCarts, carts)

	for _, theCart := range origOrderedCarts {
		moveCart(theMap, theCart)
		sort.Sort(byCartPos(carts))
		crashed, crashSite = checkForCrash(carts)
		if crashed {
			return
		}
	}

	return
}

func main() {
	lines := utils.ReadLinesOrDie("input.txt")

	theMap, carts := parseMapInput(lines)

	stepNum := 0
	var crashSite mapCoord
	for {
		stepNum++

		var crashed bool
		crashed, crashSite = executeStep(theMap, carts)
		if crashed {
			break
		}
	}

	fmt.Printf("crashed at: (%d,%d) on step %d\n",
		crashSite.x, crashSite.y, stepNum)
	printMap(theMap)
}
