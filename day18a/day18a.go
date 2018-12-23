package main

import (
	"fmt"

	"jblee.net/adventofcode2018/utils"
)

type groundMap [][]byte

const openLand = '.'
const trees = '|'
const lumberyard = '#'

func initEmptyMap(numCols, numRows int) *groundMap {
	var theMap groundMap
	theMap = make([][]byte, numCols)
	for x := 0; x < numCols; x++ {
		theMap[x] = make([]byte, numRows)
	}

	return &theMap
}

func (theMap *groundMap) size() (int, int) {
	numCols := len(*theMap)
	numRows := len((*theMap)[0])

	return numCols, numRows
}

func (theMap *groundMap) isValidXY(x, y int) bool {
	numCols := len(*theMap)
	numRows := len((*theMap)[0])

	return x >= 0 && x < numCols && y >= 0 && y < numRows
}

func (theMap *groundMap) getAt(x, y int) byte {
	if theMap.isValidXY(x, y) {
		return (*theMap)[x][y]
	}
	return 0
}

func (theMap *groundMap) putAt(x, y int, val byte) {
	if theMap.isValidXY(x, y) {
		(*theMap)[x][y] = val
	}
}

func (theMap *groundMap) print() {
	numCols, numRows := theMap.size()

	strLen := numCols * (numRows + 1)
	buf := make([]byte, strLen)

	idx := 0
	for y := 0; y < numRows; y++ {
		for x := 0; x < numCols; x++ {
			buf[idx] = (*theMap)[x][y]
			idx++
		}
		buf[idx] = '\n'
		idx++
	}

	fmt.Printf(string(buf))
}

func generateGroundMap(lines []string) *groundMap {
	numRows := len(lines)
	numCols := len(lines[0])

	theMap := initEmptyMap(numCols, numRows)
	for x := 0; x < numCols; x++ {
		for y := 0; y < numRows; y++ {
			(*theMap)[x][y] = lines[y][x]
		}
	}

	return theMap
}

func countAdjacencies(theMap *groundMap, x, y int, targetType byte) int {
	numFound := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			onStartSpot := dx == 0 && dy == 0
			if !onStartSpot && theMap.getAt(x+dx, y+dy) == targetType {
				numFound++
			}
		}
	}
	return numFound
}

func getNextLandType(theMap *groundMap, x, y int) byte {
	startType := theMap.getAt(x, y)

	resultType := startType
	if startType == openLand {
		surroundingTrees := countAdjacencies(theMap, x, y, trees)
		if surroundingTrees >= 3 {
			resultType = trees
		}
	} else if startType == trees {
		surroundingLumberyards := countAdjacencies(theMap, x, y, lumberyard)
		if surroundingLumberyards >= 3 {
			resultType = lumberyard
		}
	} else if startType == lumberyard {
		surroundingLumberyards := countAdjacencies(theMap, x, y, lumberyard)
		surroundingTrees := countAdjacencies(theMap, x, y, trees)
		if surroundingLumberyards < 1 || surroundingTrees < 1 {
			resultType = openLand
		}
	}

	return resultType
}

func doCycle(src, dst *groundMap) {
	numCols, numRows := src.size()

	for x := 0; x < numCols; x++ {
		for y := 0; y < numRows; y++ {
			nextType := getNextLandType(src, x, y)
			dst.putAt(x, y, nextType)
		}
	}
}

func countTypes(
	theMap *groundMap) (openSpots int, treeSpots int, lumberyardSpots int) {

	numCols, numRows := theMap.size()

	for x := 0; x < numCols; x++ {
		for y := 0; y < numRows; y++ {
			switch theMap.getAt(x, y) {
			case openLand:
				openSpots++
				break
			case trees:
				treeSpots++
				break
			case lumberyard:
				lumberyardSpots++
				break
			}
		}
	}

	return
}

func main() {
	lines := utils.ReadLinesOrDie("input.txt")
	// lines := utils.ReadLinesOrDie("sample_input1.txt")

	srcMap := generateGroundMap(lines)
	dstMap := initEmptyMap(srcMap.size())

	const numCycles = 10
	for cycle := 0; cycle < numCycles; cycle++ {
		doCycle(srcMap, dstMap)
		srcMap, dstMap = dstMap, srcMap
	}

	openSpots, treeSpots, lumberyardSpots := countTypes(srcMap)

	srcMap.print()
	fmt.Printf("\nopenings: %d, trees: %d, lumberyards: %d\n",
		openSpots, treeSpots, lumberyardSpots)
	fmt.Printf("trees * lumberyards: %d\n", treeSpots*lumberyardSpots)
}
