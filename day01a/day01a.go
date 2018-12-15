package main

import (
	"fmt"
	"strconv"

	"jblee.net/adventofcode2018/utils"
)

func main() {
	lines := utils.ReadLinesOrDie("input.txt")

	freq := 0
	for _, line := range lines {
		delta, _ := strconv.Atoi(line)
		freq += delta
	}

	fmt.Printf("freq: %d\n", freq)
}
