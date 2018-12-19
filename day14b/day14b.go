package main

import (
	"bytes"
	"fmt"
)

func main() {
	targetScores := [...]byte{2, 9, 0, 4, 3, 1} // the input
	targetLen := len(targetScores)

	startScores := [...]byte{3, 7}

	const StartingLen = 1024 * 512
	scores := make([]byte, 2, StartingLen)
	copy(scores, startScores[:])

	elf1Pos, elf2Pos := 0, 1

	checkIdx := 0
	for {
		combinedScore := scores[elf1Pos] + scores[elf2Pos]
		if combinedScore >= 10 {
			scores = append(scores, combinedScore/10)
		}
		scores = append(scores, combinedScore%10)

		elf1Pos = (elf1Pos + int(scores[elf1Pos]) + 1) % len(scores)
		elf2Pos = (elf2Pos + int(scores[elf2Pos]) + 1) % len(scores)

		matchFound := false
		for checkIdx+targetLen <= len(scores) {
			scoresToCheck := scores[checkIdx : checkIdx+targetLen]
			if bytes.Equal(targetScores[:], scoresToCheck) {
				matchFound = true
				break
			}
			checkIdx++
		}

		if matchFound {
			break
		}
	}

	fmt.Printf("target appears after %d recipes\n", checkIdx)
}
