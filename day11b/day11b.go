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

func getAreaPower(grid *powerGrid, x int, y int, squareSize int) int {
	power := 0
	for dx := 0; dx < squareSize; dx++ {
		for dy := 0; dy < squareSize; dy++ {
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
	var maxPowerSquareSize int
	maxAreaPower := 0
	for squareSize := 1; squareSize <= gridSize; squareSize++ {
		for y := 0; y <= len(grid)-squareSize; y++ {
			for x := 0; x <= len(grid[y])-squareSize; x++ {
				power := getAreaPower(&grid, x, y, squareSize)
				if power > maxAreaPower {
					maxAreaPower = power
					maxXCoord, maxYCoord = x+1, y+1
					maxPowerSquareSize = squareSize
				}
			}
		}
	}

	fmt.Printf("coord: %d,%d with power of %d and size %d\n",
		maxXCoord, maxYCoord, maxAreaPower, maxPowerSquareSize)
}
