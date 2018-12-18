package main

import "fmt"

const gridSize = 300

type powerGrid [gridSize][gridSize]int

func getPower(x int, y int, serialNum int) int {
	rackID := x + 10
	power := rackID * y
	power += serialNum
	power *= rackID
	power = (power / 100) % 10
	power -= 5

	return power
}

func getAreaPower(grid *powerGrid, x int, y int, dist int) int {
	power := 0
	for dx := -dist; dx <= dist; dx++ {
		for dy := -dist; dy <= dist; dy++ {
			power += grid[x+dx][y+dy]
		}
	}

	return power
}

func main() {
	const SerialNum = 5791

	var grid powerGrid
	for y := range grid {
		for x := range grid[y] {
			grid[x][y] = getPower(x+1, y+1, SerialNum)
		}
	}

	var maxXCoord, maxYCoord int
	maxAreaPower := 0
	for y := 1; y < len(grid)-1; y++ {
		for x := 1; x < len(grid[y])-1; x++ {
			power := getAreaPower(&grid, x, y, 1)
			if power > maxAreaPower {
				maxAreaPower = power
				maxXCoord, maxYCoord = x, y
			}
		}
	}

	fmt.Printf("coord: %d,%d with power of %d\n",
		maxXCoord, maxYCoord, maxAreaPower)
}
