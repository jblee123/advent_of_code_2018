package main

import (
	"fmt"

	"jblee.net/adventofcode2018/utils"
)

const patternLen = 5

func initPatterns(lines []string) map[string]bool {
	patterns := map[string]bool{}
	for _, line := range lines {
		patterns[line[0:5]] = line[9] == '#'
	}
	return patterns
}

func main() {
	lines := utils.ReadLinesOrDie("input.txt")

	const NumGenerations = 20

	patterns := initPatterns(lines[2:])

	currentState := lines[0][len("initial state: "):]
	neededPadding := NumGenerations + 2
	leftNum := -neededPadding

	workspace := make([]byte, len(currentState)+(neededPadding*2))
	for i := range workspace {
		workspace[i] = '.'
	}
	for i := range currentState {
		workspace[i+neededPadding] = currentState[i]
	}

	currentState = string(workspace)

	for gen := 0; gen < NumGenerations; gen++ {
		for i := 0; i <= len(currentState)-patternLen; i++ {
			idxToWrite := i + (patternLen / 2)
			plantSet := currentState[i : i+patternLen]
			sprouts, _ := patterns[plantSet]
			toWrite := byte('.')
			if sprouts {
				toWrite = byte('#')
			}
			workspace[idxToWrite] = toWrite
		}
		currentState = string(workspace)
		fmt.Println(currentState)
	}

	sum := 0
	for i, c := range currentState {
		if c == '#' {
			sum += i + leftNum
		}
	}

	fmt.Printf("plant sum: %d\n", sum)
}
