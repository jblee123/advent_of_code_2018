package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"jblee.net/adventofcode2018/utils"
)

func main() {
	lines := utils.ReadLinesOrDie("input.txt")
	sort.Strings(lines)

	const NumMinutes = 60
	schedule := map[int][]int{}

	var id int
	var asleepTime int
	for _, line := range lines {
		minute, _ := strconv.Atoi(line[15:17])
		event := line[19:]

		if strings.HasPrefix(event, "Guard") {
			fmt.Sscanf(event, "Guard #%d begins shift", &id)
			_, ok := schedule[id]
			if !ok {
				schedule[id] = make([]int, NumMinutes)
			}
		} else if event == "falls asleep" {
			asleepTime = minute
		} else if event == "wakes up" {
			mins := schedule[id]
			for i := asleepTime; i < minute; i++ {
				mins[i]++
			}
		} else {
			panic(fmt.Sprintf("unexpected event on line: %s", line))
		}
	}

	maxGuardSleepCount := 0
	maxGuardSleepID := 0
	maxGuardSleptMinute := 0
	for id, minutes := range schedule {
		currentMaxSleepCount := 0
		currentMaxSleptMinute := 0

		for minute, sleepCount := range minutes {
			if sleepCount > currentMaxSleepCount {
				currentMaxSleepCount = sleepCount
				currentMaxSleptMinute = minute
			}
		}

		if currentMaxSleepCount > maxGuardSleepCount {
			maxGuardSleepCount = currentMaxSleepCount
			maxGuardSleepID = id
			maxGuardSleptMinute = currentMaxSleptMinute
		}
	}

	fmt.Printf("maxGuardSleepCount: %d\n", maxGuardSleepCount)
	fmt.Printf("maxGuardSleepID: %d\n", maxGuardSleepID)
	fmt.Printf("maxGuardSleptMinute: %d\n", maxGuardSleptMinute)
	fmt.Printf("ID x Minute: %d\n", maxGuardSleepID*maxGuardSleptMinute)
}
