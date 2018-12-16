package main

import (
	"fmt"

	"jblee.net/adventofcode2018/utils"
)

func findSimilarString(lines []string) string {
	for i, line := range lines {
		for j := i + 1; j < len(lines); j++ {
			line2 := lines[j]

			diffPos := -1
			for pos := 0; pos < len(line); pos++ {
				letter := line[pos]
				if letter != line2[pos] {
					if diffPos >= 0 {
						diffPos = -1
						break
					} else {
						diffPos = pos
					}
				}
			}

			if diffPos >= 0 {
				return line[:diffPos] + line[diffPos+1:]
			}
		}
	}

	return ""
}

func main() {
	lines := utils.ReadLinesOrDie("input.txt")

	answer := findSimilarString(lines)

	fmt.Printf("answer: %s\n", answer)
}
