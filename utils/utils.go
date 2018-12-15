package utils

import (
	"bufio"
	"log"
	"os"
)

// ReadLinesOrDie reads a whole file into memory
// and returns a slice of its lines.
func ReadLinesOrDie(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}
