package main

import (
	"fmt"
	"strings"

	"jblee.net/adventofcode2018/utils"
)

type dependency struct {
	first, second int
}

func lines2Dep(lines []string) []dependency {
	deps := make([]dependency, len(lines))
	for idx, line := range lines {
		deps[idx].first = int(line[5])
		deps[idx].second = int(line[36])
	}

	return deps
}

func getStepCount(deps []dependency) int {
	steps := map[int]bool{}
	for _, dep := range deps {
		steps[dep.first] = true
		steps[dep.second] = true
	}

	return len(steps)
}

func stepHasDeps(step int, deps []dependency) bool {
	for _, dep := range deps {
		if dep.second == step {
			return true
		}
	}

	return false
}

func getNextStep(deps []dependency, sequenceSoFar string) int {
	for step := int('A'); step <= int('Z'); step++ {
		if strings.IndexRune(sequenceSoFar, rune(step)) > -1 {
			continue
		}

		if !stepHasDeps(step, deps) {
			return step
		}
	}
	return -1
}

func filterStep(step int, deps []dependency) []dependency {
	var newDeps []dependency
	for _, dep := range deps {
		if dep.first != step {
			newDeps = append(newDeps, dep)
		}
	}
	return newDeps
}

func main() {
	lines := utils.ReadLinesOrDie("input.txt")
	deps := lines2Dep(lines)
	stepCount := getStepCount(deps)

	sequence := ""
	for len(sequence) < stepCount {
		step := getNextStep(deps, sequence)
		sequence += string(step)
		deps = filterStep(step, deps)
	}

	fmt.Printf("sequence: %s\n", sequence)
}
