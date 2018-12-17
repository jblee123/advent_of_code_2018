package main

import (
	"fmt"
	"math"

	"jblee.net/adventofcode2018/utils"
)

type point struct {
	x, y int
}

type boundingBox struct {
	left, top, right, bottom int
}

type spaceData struct {
	owner int
	dist  int
}

func lines2Points(lines []string) []point {
	points := make([]point, len(lines))
	for i, line := range lines {
		fmt.Sscanf(line, "%d, %d", &points[i].x, &points[i].y)
	}

	return points
}

func getBounds(points []point) boundingBox {
	var bounds boundingBox
	bounds.left = math.MaxInt64
	bounds.top = math.MaxInt64
	bounds.right = math.MinInt64
	bounds.bottom = math.MinInt64

	for _, p := range points {
		if p.x < bounds.left {
			bounds.left = p.x
		}
		if p.x > bounds.right {
			bounds.right = p.x
		}
		if p.y < bounds.top {
			bounds.top = p.y
		}
		if p.y > bounds.bottom {
			bounds.bottom = p.y
		}
	}

	return bounds
}

func trySpread(
	xOffset int, yOffset int, origin point, spaces map[point]spaceData,
	pointQueue *[]point, areas []int, bounds boundingBox) {

	targetCoord := point{x: origin.x + xOffset, y: origin.y + yOffset}
	originSpace := spaces[origin]

	// see if the coord is outside bounds
	if targetCoord.x < bounds.left || targetCoord.x > bounds.right ||
		targetCoord.y < bounds.top || targetCoord.y > bounds.bottom {

		if originSpace.owner >= 0 {
			areas[originSpace.owner] = -1
		}
		return
	}

	newDist := originSpace.dist + 1

	// see if the coord is taken
	targetSpace, exists := spaces[targetCoord]
	if exists {
		if targetSpace.owner != originSpace.owner &&
			targetSpace.dist == newDist {

			spaces[targetCoord] = spaceData{owner: -1, dist: newDist}
		}
		return
	}

	// claim the coord
	spaces[targetCoord] = spaceData{owner: originSpace.owner, dist: newDist}
	*pointQueue = append(*pointQueue, targetCoord)
}

func printBoard(spaces map[point]spaceData, bounds boundingBox) {
	for y := bounds.top; y <= bounds.bottom; y++ {
		for x := bounds.left; x <= bounds.right; x++ {
			space := spaces[point{x: x, y: y}]
			if space.owner >= 0 {
				fmt.Printf("%c", int('A')+space.owner)
			} else {
				fmt.Printf("#")
			}
		}
		fmt.Printf("\n")
	}
}

func main() {
	lines := utils.ReadLinesOrDie("input.txt")
	points := lines2Points(lines)
	bounds := getBounds(points)

	spaces := map[point]spaceData{}
	areas := make([]int, len(points))

	for idx, point := range points {
		spaces[point] = spaceData{owner: idx, dist: 0}
	}

	pointQueue := points[:]
	for len(pointQueue) > 0 {
		trySpread(-1, 0, pointQueue[0], spaces, &pointQueue, areas, bounds)
		trySpread(0, -1, pointQueue[0], spaces, &pointQueue, areas, bounds)
		trySpread(1, 0, pointQueue[0], spaces, &pointQueue, areas, bounds)
		trySpread(0, 1, pointQueue[0], spaces, &pointQueue, areas, bounds)

		pointQueue = pointQueue[1:]
	}

	for _, space := range spaces {
		if space.owner > -1 && areas[space.owner] > -1 {
			areas[space.owner]++
		}
	}

	maxArea := math.MinInt64
	for _, area := range areas {
		if area > maxArea {
			maxArea = area
		}
	}

	// printBoard(spaces, bounds)

	fmt.Printf("max area: %d\n", maxArea)
}
