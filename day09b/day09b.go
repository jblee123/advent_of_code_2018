package main

import (
	"container/ring"
	"fmt"

	"jblee.net/adventofcode2018/utils"
)

func main() {
	line := utils.ReadLinesOrDie("input.txt")[0]

	var numPlayers, maxPointValue int
	fmt.Sscanf(line, "%d players; last marble is worth %d points",
		&numPlayers, &maxPointValue)

	maxPointValue *= 100

	currentMarble := ring.New(1)
	currentMarble.Value = 0

	player := 0
	scores := make([]int, numPlayers)

	for pointVal := 1; pointVal <= maxPointValue; pointVal++ {
		if pointVal%23 != 0 {
			newMarble := ring.New(1)
			newMarble.Value = pointVal

			currentMarble.Next().Link(newMarble)
			currentMarble = newMarble
		} else {
			for i := 0; i < 7; i++ {
				currentMarble = currentMarble.Prev()
			}

			scores[player] += pointVal + currentMarble.Value.(int)

			currentMarble = currentMarble.Prev()
			currentMarble.Unlink(1)
			currentMarble = currentMarble.Next()
		}

		player = (player + 1) % numPlayers
	}

	maxScore := 0
	for _, score := range scores {
		if score > maxScore {
			maxScore = score
		}
	}

	fmt.Printf("max score: %d\n", maxScore)
}
