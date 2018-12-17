package main

import (
	"fmt"
	"math"
	"strings"

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

func react(s string) string {
	for i := 0; i < len(s)-1; {
		if cancelOut(s[i], s[i+1]) {
			s = s[:i] + s[i+2:]
			if i > 0 {
				i--
			}
		} else {
			i++
		}
	}

	return s
}

func removeAll(s string, c byte) string {
	c = toLower(c)
	var builder strings.Builder
	for i := 0; i < len(s); i++ {
		if toLower(s[i]) != c {
			builder.WriteByte(s[i])
		}
	}

	return builder.String()
}

func main() {
	line := utils.ReadLinesOrDie("input.txt")[0]

	minLen := math.MaxInt64
	for c := byte('a'); c <= byte('z'); c++ {
		line2 := removeAll(line, c)
		line2 = react(line2)
		if len(line2) < minLen {
			minLen = len(line2)
		}
	}

	fmt.Printf("minLen: %d\n", minLen)
}
