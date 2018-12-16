package main

import (
	"fmt"

	"jblee.net/adventofcode2018/utils"
)

func main() {
	lines := utils.ReadLinesOrDie("input.txt")

	const NumLetters = 26

	repeated2x := 0
	repeated3x := 0

	for _, line := range lines {
		var letterCnts [NumLetters]int
		for _, letter := range line {
			letterCnts[letter-'a']++
		}

		repeated2xFound := false
		repeated3xFound := false

		for _, cnt := range letterCnts {
			if cnt == 2 {
				repeated2xFound = true
			}
			if cnt == 3 {
				repeated3xFound = true
			}
		}

		if repeated2xFound {
			repeated2x++
		}
		if repeated3xFound {
			repeated3x++
		}
	}

	checksum := repeated2x * repeated3x

	fmt.Printf("checksum: %d\n", checksum)
}
