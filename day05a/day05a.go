package main

import (
	"fmt"

	"jblee.net/adventofcode2018/utils"
)

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return 'a' + (c - 'A')
	}
	return c
}

func cancelOut(a, b byte) bool {
	aLower := toLower(a)
	bLower := toLower(b)
	return aLower == bLower && a != b
}

func main() {
	line := utils.ReadLinesOrDie("input.txt")[0]

	for i := 0; i < len(line)-1; {
		if cancelOut(line[i], line[i+1]) {
			line = line[:i] + line[i+2:]
			if i > 0 {
				i--
			}
		} else {
			i++
		}
	}

	fmt.Printf("line len: %d\n", len(line))
}
