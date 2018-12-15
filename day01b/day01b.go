package main

import (
	"fmt"
	"strconv"

	"jblee.net/adventofcode2018/utils"
)

func main() {
	lines := utils.ReadLinesOrDie("input.txt")

	freq := 0
	var firstFreqRepeat int
	seenRepeat := false
	seenFreqs := make(map[int]bool)
	seenFreqs[freq] = true

	for !seenRepeat {
		for _, line := range lines {
			delta, _ := strconv.Atoi(line)
			freq += delta
			if seenFreqs[freq] {
				firstFreqRepeat = freq
				seenRepeat = true
				break
			}

			seenFreqs[freq] = true
		}
	}

	fmt.Printf("first repeated freq: %d\n", firstFreqRepeat)
}
