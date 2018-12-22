package main

import (
	"fmt"
	"math"
	"sort"
	"time"

	"jblee.net/adventofcode2018/utils"
)

type mapDir byte

const (
	left  mapDir = iota
	up    mapDir = iota
	right mapDir = iota
	down  mapDir = iota
)

type mapCoord struct {
	x, y int
}

type unit struct {
	unitType byte
	pos      mapCoord
	hp       int
	attack   int
}

type mineMap [][]byte

func readingOrderLess(x1, y1, x2, y2 int) bool {
	return (y1 < y2) || ((y1 == y2) && (x1 < x2))
}

////////////////////////////////////////////////////////////////////////////////

type byUnitPos [](*unit)

func (a byUnitPos) Len() int {
	return len(a)
}

func (a byUnitPos) Less(i, j int) bool {
	return readingOrderLess(a[i].pos.x, a[i].pos.y, a[j].pos.x, a[j].pos.y)
}

func (a byUnitPos) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

////////////////////////////////////////////////////////////////////////////////

type byCoordReadingOrder []mapCoord

func (a byCoordReadingOrder) Len() int {
	return len(a)
}

func (a byCoordReadingOrder) Less(i, j int) bool {
	return readingOrderLess(a[i].x, a[i].y, a[j].x, a[j].y)
}

func (a byCoordReadingOrder) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

////////////////////////////////////////////////////////////////////////////////

func createtUnit(unitType byte, x, y int) *unit {
	return &unit{
		unitType: unitType,
		pos:      mapCoord{x, y},
		hp:       200,
		attack:   3,
	}
}

func parseMapInput(lines []string, elfAttack int) (
	mineMap, [](*unit), [](*unit)) {

	xLen := len(lines[0])
	yLen := len(lines)

	theMap := make(mineMap, xLen)
	elves := make([](*unit), 0)
	goblins := make([](*unit), 0)
	for x := 0; x < xLen; x++ {
		theMap[x] = make([]byte, yLen)
		for y := 0; y < yLen; y++ {
			c := lines[y][x]
			theMap[x][y] = c
			if c == 'E' {
				elves = append(elves, createtUnit(c, x, y))
				elves[len(elves)-1].attack = elfAttack
			} else if c == 'G' {
				goblins = append(goblins, createtUnit(c, x, y))
			}
		}
	}

	return theMap, elves, goblins
}

func printMap(theMap mineMap) {
	xLen := len(theMap)
	yLen := len(theMap[0])
	strLen := (xLen + 1) * yLen
	byteBuf := make([]byte, strLen)
	bufIdx := 0
	for y := 0; y < yLen; y++ {
		for x := 0; x < xLen; x++ {
			byteBuf[bufIdx] = theMap[x][y]
			bufIdx++
		}
		byteBuf[bufIdx] = '\n'
		bufIdx++
	}

	fmt.Print(string(byteBuf))
}

func printUnits(elves []*unit, goblins []*unit) {
	for i, elf := range elves {
		fmt.Printf("Elf %2d: %d\n", i+1, elf.hp)
	}

	fmt.Println()
	for i, goblin := range goblins {
		fmt.Printf("Goblin %2d: %d\n", i+1, goblin.hp)
	}
}

func combatRemaining(elves []*unit, goblins []*unit) bool {
	elvesAlive := false
	for _, elf := range elves {
		elvesAlive = elvesAlive || (elf.hp > 0)
	}
	goblinsAlive := false
	for _, goblin := range goblins {
		goblinsAlive = goblinsAlive || (goblin.hp > 0)
	}

	return elvesAlive && goblinsAlive
}

func getUnitMoveOrder(elves []*unit, goblins []*unit) []*unit {
	units := make([]*unit, 0, len(elves)+len(goblins))
	units = append(units, elves...)
	units = append(units, goblins...)
	sort.Sort(byUnitPos(units))
	return units
}

func oppositeUnitType(unitType byte) byte {
	if unitType == 'E' {
		return 'G'
	} else {
		return 'E'
	}
}

func mapSpaceIsValid(theMap mineMap, coord mapCoord) bool {
	xLen := len(theMap)
	yLen := len(theMap[0])
	if coord.x < 0 || coord.x >= xLen || coord.y < 0 || coord.y >= yLen {
		return false
	}
	return true
}

func mapSpaceIs(spaceType byte, theMap mineMap, coord mapCoord) bool {
	if !mapSpaceIsValid(theMap, coord) {
		return false
	}

	return theMap[coord.x][coord.y] == spaceType
}

func coordTo(dir mapDir, coord mapCoord) mapCoord {
	switch dir {
	case left:
		return mapCoord{coord.x - 1, coord.y}
	case up:
		return mapCoord{coord.x, coord.y - 1}
	case right:
		return mapCoord{coord.x + 1, coord.y}
	case down:
		return mapCoord{coord.x, coord.y + 1}
	}

	panic(fmt.Sprintf("bad dir: %v", dir))
}

func getElfOrGoblinAt(pos mapCoord, elves []*unit, goblins []*unit) *unit {
	for _, elf := range elves {
		if elf.pos == pos {
			return elf
		}
	}

	for _, goblin := range goblins {
		if goblin.pos == pos {
			return goblin
		}
	}

	panic("bad pos passed to getElfOrGoblinAt")
}

func isTargetNear(toMove *unit, theMap mineMap, elves []*unit,
	goblins []*unit) (bool, mapCoord) {

	targetType := oppositeUnitType(toMove.unitType)

	coords := [...]mapCoord{
		coordTo(up, toMove.pos),
		coordTo(left, toMove.pos),
		coordTo(right, toMove.pos),
		coordTo(down, toMove.pos),
	}

	var targetCoord mapCoord
	targetHp := math.MaxInt64

	for _, coord := range coords {
		if mapSpaceIs(targetType, theMap, coord) {
			enemy := getElfOrGoblinAt(coord, elves, goblins)
			if enemy.hp < targetHp {
				targetCoord = coord
				targetHp = enemy.hp
			}
		}
	}

	return targetHp != math.MaxInt64, targetCoord
}

func attackUnit(toAttack *unit, attacker *unit, theMap mineMap) {
	toAttack.hp -= attacker.attack
	if toAttack.hp <= 0 {
		toAttack.hp = 0
		theMap[toAttack.pos.x][toAttack.pos.y] = '.'
		toAttack.pos = mapCoord{-1, -1}
	}
}

func checkCoordForSearch(coord mapCoord, theMap mineMap, targetType byte,
	searchGrid [][]int, nextDist int) byte {

	if !mapSpaceIsValid(theMap, coord) {
		return 0
	}

	if searchGrid[coord.x][coord.y] > -1 {
		return 0
	}

	c := theMap[coord.x][coord.y]
	if c == '.' || c == targetType {
		searchGrid[coord.x][coord.y] = nextDist
		return c
	}

	return 0
}

func getNextSearchSeeds(theMap mineMap, targetType byte,
	searchGrid [][]int, searchSeeds []mapCoord) ([]mapCoord, bool) {

	enemyFound := false

	x1 := searchSeeds[0].x
	y1 := searchSeeds[0].y
	nextDist := searchGrid[x1][y1] + 1

	nextSeeds := make([]mapCoord, 0, len(searchSeeds))

	for _, coord := range searchSeeds {
		adjCoords := [...]mapCoord{
			coordTo(up, coord),
			coordTo(left, coord),
			coordTo(right, coord),
			coordTo(down, coord),
		}

		for _, adjCoord := range adjCoords {
			chk := checkCoordForSearch(
				adjCoord, theMap, targetType, searchGrid, nextDist)
			if chk != 0 {
				nextSeeds = append(nextSeeds, adjCoord)
				if chk != '.' {
					enemyFound = true
				}
			}
		}
	}

	return nextSeeds, enemyFound
}

// assuming an enemy exists
func getEnemyToTargetFromSeeds(theMap mineMap, targetType byte,
	searchSeeds []mapCoord) mapCoord {

	enemies := make([]mapCoord, 0)
	for _, searchSeed := range searchSeeds {
		if mapSpaceIs(targetType, theMap, searchSeed) {
			enemies = append(enemies, searchSeed)
		}
	}
	sort.Sort(byCoordReadingOrder(enemies))
	return enemies[0]
}

func getNextTracebackSeeds(tracebackSeeds []mapCoord,
	searchGrid [][]int) []mapCoord {

	x := tracebackSeeds[0].x
	y := tracebackSeeds[0].y
	curDist := searchGrid[x][y] - 1

	newSeeds := make([]mapCoord, 0)

	for _, seed := range tracebackSeeds {
		adjCoords := [...]mapCoord{
			coordTo(up, seed),
			coordTo(left, seed),
			coordTo(right, seed),
			coordTo(down, seed),
		}

		for _, adjCoord := range adjCoords {
			if searchGrid[adjCoord.x][adjCoord.y] == curDist {
				newSeeds = append(newSeeds, adjCoord)
			}
		}
	}

	return newSeeds
}

func printSearchGrid(searchGrid [][]int) {
	xLen := len(searchGrid)
	yLen := len(searchGrid[0])

	for y := 0; y < yLen; y++ {
		for x := 0; x < xLen; x++ {
			fmt.Printf("%2d ", searchGrid[x][y])
		}
		fmt.Printf("\n")
	}
}

func doMoveSearch(
	theMap mineMap, targetType byte,
	searchGrid [][]int, searchSeeds []mapCoord) []mapCoord {

	nextSeeds, enemyFound := getNextSearchSeeds(
		theMap, targetType, searchGrid, searchSeeds)

	// printSearchGrid(searchGrid)
	// fmt.Println()

	if len(nextSeeds) == 0 {
		return nil
	}

	var tracebackSeeds []mapCoord

	if enemyFound {
		tracebackSeeds = make([]mapCoord, 1)
		tracebackSeeds[0] = getEnemyToTargetFromSeeds(
			theMap, targetType, nextSeeds)
	} else {
		tracebackSeeds = doMoveSearch(theMap, targetType, searchGrid, nextSeeds)
		if len(tracebackSeeds) > 0 {
			tracebackSeeds = getNextTracebackSeeds(tracebackSeeds, searchGrid)
		}
	}

	return tracebackSeeds
}

func tryMoveTowardEnemy(toMove *unit, theMap mineMap) {
	targetType := oppositeUnitType(toMove.unitType)

	xLen := len(theMap)
	yLen := len(theMap[0])
	searchGrid := make([][]int, xLen)
	for x := 0; x < xLen; x++ {
		searchGrid[x] = make([]int, yLen)
		for y := 0; y < yLen; y++ {
			searchGrid[x][y] = -1
		}
	}

	searchGrid[toMove.pos.x][toMove.pos.y] = 0
	searchSeeds := make([]mapCoord, 1)
	searchSeeds[0] = toMove.pos

	moveOptions := doMoveSearch(theMap, targetType, searchGrid, searchSeeds)

	if len(moveOptions) > 0 {
		sort.Sort(byCoordReadingOrder(moveOptions))
		theMap[toMove.pos.x][toMove.pos.y] = '.'
		toMove.pos = moveOptions[0]
		theMap[toMove.pos.x][toMove.pos.y] = toMove.unitType
	}
}

func doMove(toMove *unit, theMap mineMap, elves []*unit, goblins []*unit) {
	if toMove.hp <= 0 {
		return
	}

	if isNear, pos := isTargetNear(toMove, theMap, elves, goblins); isNear {
		toAttack := getElfOrGoblinAt(pos, elves, goblins)
		attackUnit(toAttack, toMove, theMap)
	} else {
		tryMoveTowardEnemy(toMove, theMap)

		if isNear, pos := isTargetNear(toMove, theMap, elves, goblins); isNear {
			toAttack := getElfOrGoblinAt(pos, elves, goblins)
			attackUnit(toAttack, toMove, theMap)
		}
	}
}

func sumRemainingUnitHp(elves []*unit, goblins []*unit) int {
	hp := 0

	for _, elf := range elves {
		hp += elf.hp
	}
	for _, goblin := range goblins {
		hp += goblin.hp
	}

	return hp
}

func runBattle(theMap mineMap, elves [](*unit), goblins [](*unit)) {
	// printMap(theMap)
	// fmt.Println()
	// printUnits(elves, goblins)

	fullRounds := 0
	finalFullRounds := -1
	for combatRemaining(elves, goblins) {
		unitsInMoveOrder := getUnitMoveOrder(elves, goblins)
		effectivelyFinished := false
		for _, toMove := range unitsInMoveOrder {
			if toMove.hp <= 0 {
				continue
			}

			if effectivelyFinished {
				finalFullRounds = fullRounds
			}
			doMove(toMove, theMap, elves, goblins)
			effectivelyFinished = !combatRemaining(elves, goblins)
		}

		fullRounds++

		// fmt.Println()
		// fmt.Printf("After round %d ======================\n", fullRounds)
		// printMap(theMap)
		// fmt.Println()
		// printUnits(elves, goblins)
	}

	if finalFullRounds == -1 {
		finalFullRounds = fullRounds
	}

	hpSum := sumRemainingUnitHp(elves, goblins)

	// printMap(theMap)
	// fmt.Println()
	printUnits(elves, goblins)

	fmt.Println()
	fmt.Printf("full rounds: %d\n", finalFullRounds)
	fmt.Printf("hp sum: %d\n", hpSum)
	fmt.Printf("rounds * hp: %d\n", finalFullRounds*hpSum)
}

func main() {
	lines := utils.ReadLinesOrDie("input.txt")
	// lines := utils.ReadLinesOrDie("sample_input2.txt")

	elfAttack := 3
	allElvesSurvived := false
	for !allElvesSurvived {
		fmt.Printf("========== testing elf attack of %d ==========\n",
			elfAttack)
		theMap, elves, goblins := parseMapInput(lines, elfAttack)

		start := time.Now()
		runBattle(theMap, elves, goblins)
		elapsed := time.Since(start)
		fmt.Printf("\nbattle took %s\n", elapsed)

		allElvesSurvived = true
		for _, elf := range elves {
			if elf.hp <= 0 {
				allElvesSurvived = false
			}
		}

		elfAttack++

		fmt.Println()
	}
}
