package main

import (
	"fmt"

	"jblee.net/adventofcode2018/utils"
)

type sheetPiece struct {
	id     int
	x      int
	y      int
	width  int
	height int
}

func linesToPieces(lines []string) []sheetPiece {
	pieces := []sheetPiece{}

	for _, line := range lines {
		var a sheetPiece
		// #1 @ 555,891: 18x12
		fmt.Sscanf(line, "#%d @ %d,%d: %dx%d",
			&a.id, &a.x, &a.y, &a.width, &a.height)
		pieces = append(pieces, a)
	}

	return pieces
}

func main() {
	lines := utils.ReadLinesOrDie("input.txt")
	pieces := linesToPieces(lines)

	const SideLen = 1000
	var sheet [SideLen][SideLen]int

	for _, piece := range pieces {
		maxX := piece.x + piece.width
		maxY := piece.y + piece.height
		for x := piece.x; x < maxX; x++ {
			for y := piece.y; y < maxY; y++ {
				sheet[x][y]++
			}
		}
	}

	goodId := 0
	for _, piece := range pieces {
		maxCovered := 0
		maxX := piece.x + piece.width
		maxY := piece.y + piece.height
		for x := piece.x; x < maxX; x++ {
			for y := piece.y; y < maxY; y++ {
				if sheet[x][y] > maxCovered {
					maxCovered = sheet[x][y]
				}
			}
		}

		if maxCovered == 1 {
			goodId = piece.id
			break
		}
	}

	fmt.Printf("goodId: %d\n", goodId)
}
