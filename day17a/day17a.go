package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"jblee.net/adventofcode2018/utils"
)

type groundCoord struct {
	x, y int
}

type coordSet map[groundCoord]bool

type groundMap struct {
	theGrid [][]byte
	minX    int
}

const springX = 500
const springY = 0

const sand = ' '
const clay = '#'
const spring = '+'
const flowingWater = '|'
const pooledWater = '~'

type gridLine struct {
	start, end groundCoord
}

func (theMap *groundMap) validXY(x, y int) bool {
	return x >= 0 && x < len(theMap.theGrid) &&
		y >= 0 && y < len(theMap.theGrid[0])
}

func (theMap *groundMap) validCoord(coord groundCoord) bool {
	return theMap.validXY(coord.x, coord.y)
}

func (theMap *groundMap) atXY(x, y int) byte {
	if theMap.validXY(x, y) {
		return theMap.theGrid[x][y]
	} else {
		return '.'
	}
}

func (theMap *groundMap) atCoord(coord groundCoord) byte {
	return theMap.atXY(coord.x, coord.y)
}

func (theMap *groundMap) putAtXY(x, y int, val byte) {
	if theMap.validXY(x, y) {
		theMap.theGrid[x][y] = val
	}
}

func (theMap *groundMap) putAtCoord(coord groundCoord, val byte) {
	theMap.putAtXY(coord.x, coord.y, val)
}

func printGroundMap(theMap *groundMap) {
	numCols := len(theMap.theGrid)
	// maxX := numCols + theMap.minX - 1
	maxY := len(theMap.theGrid[0])
	yCoordLen := int(math.Ceil(math.Log10(float64(maxY))))

	// +1 for space between coords and grid
	lineLen := numCols + yCoordLen + 1
	lineBuf := make([]byte, lineLen)

	// prep for first three lines
	for i := 0; i <= yCoordLen; i++ {
		lineBuf[i] = ' '
	}

	// line 1
	for i := 0; i < numCols; i++ {
		xVal := i + theMap.minX
		var c byte
		if xVal < 100 {
			c = ' '
		} else {
			c = '0' + byte((xVal/100)%10)
		}
		lineBuf[yCoordLen+1+i] = c
	}
	fmt.Printf("%s\n", string(lineBuf))

	// line 2
	for i := 0; i < numCols; i++ {
		xVal := i + theMap.minX
		var c byte
		if xVal < 10 {
			c = ' '
		} else {
			c = '0' + byte((xVal/10)%10)
		}
		lineBuf[yCoordLen+1+i] = c
	}
	fmt.Printf("%s\n", string(lineBuf))

	// line 3
	for i := 0; i < numCols; i++ {
		xVal := i + theMap.minX
		lineBuf[yCoordLen+1+i] = '0' + byte(xVal%10)
	}
	fmt.Printf("%s\n", string(lineBuf))

	// lines for rows
	for y := 0; y <= maxY-1; y++ {
		bufIdx := 0
		for i := yCoordLen - 1; i >= 0; i-- {
			var c byte
			digit := (y / int(math.Pow10(i)))
			if i > 0 && digit == 0 {
				c = ' '
			} else {
				c = '0' + byte(digit%10)
			}
			lineBuf[bufIdx] = c
			bufIdx++
		}

		lineBuf[bufIdx] = ' '
		bufIdx++

		for x := 0; x < numCols; x++ {
			lineBuf[bufIdx] = theMap.atXY(x, y)
			bufIdx++
		}

		fmt.Printf("%s\n", string(lineBuf))
	}
}

func parseGridLine(line string) gridLine {
	var theGridLine gridLine
	lineParts := strings.Split(line, ", ")

	linePart0Sides := strings.Split(lineParts[0], "=")
	linePart0StartEnd := strings.Split(linePart0Sides[1], "..")

	linePart1Sides := strings.Split(lineParts[1], "=")
	linePart1StartEnd := strings.Split(linePart1Sides[1], "..")

	var xSideStartEnd, ySideStartEnd *([]string)
	if linePart0Sides[0] == "x" {
		xSideStartEnd = &linePart0StartEnd
		ySideStartEnd = &linePart1StartEnd
	} else {
		ySideStartEnd = &linePart0StartEnd
		xSideStartEnd = &linePart1StartEnd
	}

	theGridLine.start.x, _ = strconv.Atoi((*xSideStartEnd)[0])
	if len(*xSideStartEnd) > 1 {
		theGridLine.end.x, _ = strconv.Atoi((*xSideStartEnd)[1])
	} else {
		theGridLine.end.x = theGridLine.start.x
	}

	theGridLine.start.y, _ = strconv.Atoi((*ySideStartEnd)[0])
	if len(*ySideStartEnd) > 1 {
		theGridLine.end.y, _ = strconv.Atoi((*ySideStartEnd)[1])
	} else {
		theGridLine.end.y = theGridLine.start.y
	}

	return theGridLine
}

func generateGroundMap(lines []string) *groundMap {
	minX := math.MaxInt64
	maxX, maxY := math.MinInt64, math.MinInt64
	gridLines := make([]gridLine, len(lines))
	for i, line := range lines {
		gridLines[i] = parseGridLine(line)
		if gridLines[i].start.x < minX {
			minX = gridLines[i].start.x
		}
		if gridLines[i].end.x > maxX {
			maxX = gridLines[i].end.x
		}
		if gridLines[i].end.y > maxY {
			maxY = gridLines[i].end.y
		}
	}

	minX--
	maxX++

	numCols := maxX - minX + 1
	numRows := maxY + 1
	groundMapGrid := make([][]byte, numCols)
	for i := 0; i < numCols; i++ {
		groundMapGrid[i] = make([]byte, numRows)
	}

	for x := 0; x < numCols; x++ {
		for y := 0; y < numRows; y++ {
			groundMapGrid[x][y] = sand
		}
	}

	for _, aGridLine := range gridLines {
		x1, y1 := aGridLine.start.x-minX, aGridLine.start.y
		x2, y2 := aGridLine.end.x-minX, aGridLine.end.y
		for {
			groundMapGrid[x1][y1] = clay
			if x1 == x2 && y1 == y2 {
				break
			}

			if x1 != x2 {
				x1++
			}
			if y1 != y2 {
				y1++
			}
		}
	}

	groundMapGrid[springX-minX][springY] = spring

	return &groundMap{groundMapGrid, minX}
}

func canPoolWater(theMap *groundMap, pos groundCoord) bool {
	for {
		onLeft := theMap.atXY(pos.x-1, pos.y)
		under := theMap.atXY(pos.x, pos.y+1)
		underSupported := under == clay || under == pooledWater
		if !underSupported || onLeft == sand {
			return false
		}
		if onLeft == clay {
			break
		}
		pos.x--
	}

	for {
		onRight := theMap.atXY(pos.x+1, pos.y)
		under := theMap.atXY(pos.x, pos.y+1)
		underSupported := under == clay || under == pooledWater
		if !underSupported || onRight == sand {
			return false
		}
		if onRight == clay {
			return true
		}
		pos.x++
	}
}

func poolWater(theMap *groundMap, pos groundCoord, newHotspots coordSet) {
	for theMap.atXY(pos.x-1, pos.y) != clay {
		pos.x--
	}
	for theMap.atCoord(pos) != clay {
		theMap.putAtCoord(pos, pooledWater)

		coordAbove := groundCoord{pos.x, pos.y - 1}
		if theMap.atCoord(coordAbove) == flowingWater {
			newHotspots[coordAbove] = true
		}

		pos.x++
	}
}

func doCycle(theMap *groundMap, hotspots coordSet) coordSet {

	newHotspots := coordSet{}

	setFlowingWater := func(coord groundCoord) {
		theMap.putAtCoord(coord, flowingWater)
		newHotspots[coord] = true
	}

	// starting or ending...
	if len(hotspots) == 0 {
		underSpringCoord := groundCoord{springX - theMap.minX, 1}
		if theMap.atCoord(underSpringCoord) == sand {
			setFlowingWater(underSpringCoord)
		}

		return newHotspots
	}

	for hotspot := range hotspots {
		if !theMap.validCoord(hotspot) {
			continue
		}

		coordBelow := groundCoord{hotspot.x, hotspot.y + 1}
		coordLeft := groundCoord{hotspot.x - 1, hotspot.y}
		coordRight := groundCoord{hotspot.x + 1, hotspot.y}

		spaceAt := theMap.atCoord(hotspot)
		spaceBelow := theMap.atCoord(coordBelow)
		spaceLeft := theMap.atCoord(coordLeft)
		spaceRight := theMap.atCoord(coordRight)

		if spaceAt == flowingWater {
			if spaceBelow == sand {
				setFlowingWater(coordBelow)
			} else if spaceBelow == clay || spaceBelow == pooledWater {
				if spaceLeft == sand {
					setFlowingWater(coordLeft)
				}
				if spaceRight == sand {
					setFlowingWater(coordRight)
				}

				checkForPool := false ||
					(spaceLeft == flowingWater && spaceRight == flowingWater) ||
					(spaceLeft == clay || spaceRight == clay)

				if checkForPool {
					if canPoolWater(theMap, hotspot) {
						poolWater(theMap, hotspot, newHotspots)
					}
				}
			}
		}
	}

	return newHotspots
}

func countWetSpots(theMap *groundMap) int {
	count := 0
	xLen := len(theMap.theGrid)
	yLen := len(theMap.theGrid[0])

	topClay := -1
	for y := 0; y < yLen && topClay == -1; y++ {
		for x := 0; x < xLen && topClay == -1; x++ {
			if theMap.atXY(x, y) == clay {
				topClay = y
			}
		}
	}

	for x := 0; x < xLen; x++ {
		for y := topClay; y < yLen; y++ {
			c := theMap.atXY(x, y)
			if c == flowingWater || c == pooledWater {
				count++
			}
		}
	}

	return count
}

func main() {
	lines := utils.ReadLinesOrDie("input.txt")
	// lines := utils.ReadLinesOrDie("sample_input1.txt")

	theMap := generateGroundMap(lines)
	hotspots := coordSet{}

	stepsRun := 0
	for {
		hotspots = doCycle(theMap, hotspots)

		if len(hotspots) == 0 {
			break
		}

		stepsRun++
	}

	wetSpots := countWetSpots(theMap)

	fmt.Printf("steps run: %d\n", stepsRun)
	fmt.Printf("wet spots: %d\n\n", wetSpots)

	printGroundMap(theMap)
}
