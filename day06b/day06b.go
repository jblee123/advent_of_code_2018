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

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func getDist(x1, y1, x2, y2 int) int {
	return absInt(x1-x2) + absInt(y1-y2)
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

func main() {
	lines := utils.ReadLinesOrDie("input.txt")
	points := lines2Points(lines)
	bounds := getBounds(points)

	const DistBound = 10000

	regionSize := 0
	for y := bounds.top; y <= bounds.bottom; y++ {
		for x := bounds.left; x <= bounds.right; x++ {
			totalDist := 0
			for _, point := range points {
				dist := getDist(x, y, point.x, point.y)
				totalDist += dist
			}
			if totalDist <= DistBound {
				regionSize++
			}
		}
	}

	fmt.Printf("region size: %d\n", regionSize)
}
