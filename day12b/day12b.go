package main

import (
	"fmt"
	"strings"

	"jblee.net/adventofcode2018/utils"
)

const patternLen = 5
const paddingLen = patternLen
const padding = "....."

func initPatterns(lines []string) map[string]byte {
	patterns := map[string]byte{}
	for _, line := range lines {
		patterns[line[0:5]] = line[9]
	}
	return patterns
}

func normalizeState(state string, leftPotNum int) (string, int) {
	firstPlant := strings.IndexByte(state, '#')
	if firstPlant < paddingLen {
		state = padding[:paddingLen-firstPlant] + state
	} else if firstPlant > paddingLen {
		state = state[firstPlant-paddingLen:]
	}

	firstPlantNum := leftPotNum + firstPlant
	leftPotNum = firstPlantNum - paddingLen

	lastPlant := strings.LastIndexByte(state, '#')
	state = (state + padding)[:lastPlant+paddingLen+1]

	return state, leftPotNum
}

func main() {
	lines := utils.ReadLinesOrDie("input.txt")

	const NumGenerations = 50000000000

	patterns := initPatterns(lines[2:])

	state := padding + lines[0][len("initial state: "):] + padding
	leftPotNum := -paddingLen
	fmt.Println(state)

	state, leftPotNum = normalizeState(state, leftPotNum)
	fmt.Printf("%5d: %s\n", leftPotNum, state)

	workspace := make([]byte, len(state))

	lastState := state
	genNum := 0
	lastLeftPotNum := leftPotNum
	for {
		workspace = workspace[0:0]
		workspace = append(workspace, '.')
		workspace = append(workspace, '.')

		for i := 0; i <= len(state)-patternLen; i++ {
			plantSet := state[i : i+patternLen]
			marker, _ := patterns[plantSet]
			workspace = append(workspace, marker)
		}

		workspace = append(workspace, '.')
		workspace = append(workspace, '.')

		state, leftPotNum = normalizeState(string(workspace), leftPotNum)
		genNum++

		fmt.Printf("%5d: %s\n", leftPotNum, state)

		if lastState == state {
			break
		}

		lastState = state
		lastLeftPotNum = leftPotNum
	}

	deltaLeftPotNum := leftPotNum - lastLeftPotNum
	generationsToGo := NumGenerations - genNum
	leftPotNum += deltaLeftPotNum * generationsToGo

	sum := 0
	for i, c := range state {
		if c == '#' {
			sum += i + leftPotNum
		}
	}

	fmt.Printf("plant sum: %d\n", sum)
	fmt.Printf("stabalized after %d generations\n", genNum)
}
