package main

import (
	"fmt"
	"sort"
	"strings"

	"jblee.net/adventofcode2018/utils"
)

type dependency struct {
	first, second int
}

type task struct {
	step    int
	endTime int
}

type byEndTime []task

func (a byEndTime) Len() int           { return len(a) }
func (a byEndTime) Less(i, j int) bool { return a[i].endTime < a[j].endTime }
func (a byEndTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

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

func getNextSteps(deps []dependency, stepsStarted string, maxSteps int) []int {
	var steps []int
	for step := int('A'); step <= int('Z') && len(steps) < maxSteps; step++ {
		if strings.IndexRune(stepsStarted, rune(step)) > -1 {
			continue
		}

		if !stepHasDeps(step, deps) {
			steps = append(steps, step)
		}
	}
	return steps
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

func timeForStep(step int) int {
	return (step - 'A') + 60 + 1
}

func main() {
	lines := utils.ReadLinesOrDie("input.txt")
	deps := lines2Dep(lines)
	stepCount := getStepCount(deps)

	stepsStarted := ""
	idleWorkers := 5
	var finishEvents []task
	currentTime := 0

	for len(stepsStarted) < stepCount || len(finishEvents) > 0 {
		steps := getNextSteps(deps, stepsStarted, idleWorkers)
		fmt.Printf("got steps: %v\n", steps)
		for _, step := range steps {
			finishTime := currentTime + timeForStep(step)
			finishEvents = append(finishEvents,
				task{step: step, endTime: finishTime})
			stepsStarted += string(step)
		}
		sort.Sort(byEndTime(finishEvents))
		idleWorkers -= len(steps)

		currentTime = finishEvents[0].endTime
		deps = filterStep(finishEvents[0].step, deps)
		finishEvents = finishEvents[1:]
		idleWorkers++
	}

	fmt.Printf("time: %d\n", currentTime)
}
