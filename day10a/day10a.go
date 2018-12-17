package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"jblee.net/adventofcode2018/utils"
)

type movingPoint struct {
	x, y       int
	xVel, yVel int
}

type boundingBox struct {
	left, top, right, bottom int
}

func linesToPoints(lines []string) []movingPoint {
	points := make([]movingPoint, len(lines))
	for i, line := range lines {
		points[i].x, _ = strconv.Atoi(strings.TrimSpace(line[10:16]))
		points[i].y, _ = strconv.Atoi(strings.TrimSpace(line[17:24]))
		points[i].xVel, _ = strconv.Atoi(strings.TrimSpace(line[36:38]))
		points[i].yVel, _ = strconv.Atoi(strings.TrimSpace(line[39:42]))
	}

	return points
}

func getBounds(points []movingPoint) boundingBox {
	box := boundingBox{
		left: math.MaxInt64, top: math.MaxInt64,
		right: math.MinInt64, bottom: math.MinInt64}

	for _, p := range points {
		if p.x < box.left {
			box.left = p.x
		}
		if p.x > box.right {
			box.right = p.x
		}
		if p.y < box.top {
			box.top = p.y
		}
		if p.y > box.bottom {
			box.bottom = p.y
		}
	}

	return box
}

func doTimestep(points []movingPoint, moveFwd bool) {
	for i := range points {
		if moveFwd {
			points[i].x += points[i].xVel
			points[i].y += points[i].yVel
		} else {
			points[i].x -= points[i].xVel
			points[i].y -= points[i].yVel
		}
	}
}

func main() {
	lines := utils.ReadLinesOrDie("input.txt")
	points := linesToPoints(lines)

	time := 0
	minRange := math.MaxInt64
	minRangeTime := 0
	var minRangeBounds boundingBox
	for {
		bounds := getBounds(points)

		yRange := bounds.bottom - bounds.top + 1
		if yRange < minRange {
			minRange = yRange
			minRangeTime = time
			minRangeBounds = bounds
		} else {
			break
		}

		doTimestep(points, true)

		time++
	}

	fmt.Printf(
		"min range: %d (at %d sec), min range bounds: %v\n",
		minRange, minRangeTime, minRangeBounds)

	doTimestep(points, false)

	xRange := minRangeBounds.right - minRangeBounds.left + 1
	outLines := make([][]byte, minRange)
	for i := range outLines {
		outLines[i] = make([]byte, xRange)
		for j := range outLines[i] {
			outLines[i][j] = ' '
		}
	}

	for i := range points {
		x := points[i].x - minRangeBounds.left
		y := points[i].y - minRangeBounds.top
		outLines[y][x] = '#'
	}

	fmt.Println()
	for i := range outLines {
		fmt.Println(string(outLines[i]))
	}
}
