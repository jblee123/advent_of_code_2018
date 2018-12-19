package main

import (
	"fmt"
)

func main() {
	const NumPracticeRecipes = 290431 // the input
	const NumRecipesToReport = 10
	const NumRecipesToMake = NumPracticeRecipes + NumRecipesToReport
	startScores := [...]byte{3, 7}

	scores := make([]byte, 2, NumRecipesToMake+2)
	copy(scores, startScores[:])

	elf1Pos, elf2Pos := 0, 1

	for len(scores) < NumRecipesToMake {
		combinedScore := scores[elf1Pos] + scores[elf2Pos]
		if combinedScore >= 10 {
			scores = append(scores, combinedScore/10)
		}
		scores = append(scores, combinedScore%10)

		elf1Pos = (elf1Pos + int(scores[elf1Pos]) + 1) % len(scores)
		elf2Pos = (elf2Pos + int(scores[elf2Pos]) + 1) % len(scores)
	}

	lastScoresTxt := make([]byte, NumRecipesToReport)
	for i, score := range scores[NumPracticeRecipes:NumRecipesToMake] {
		lastScoresTxt[i] = score + byte('0')
	}

	fmt.Println(string(lastScoresTxt))
}
