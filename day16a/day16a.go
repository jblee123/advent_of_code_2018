package main

import (
	"fmt"

	"jblee.net/adventofcode2018/utils"
)

const numRegisters = 4
const numOpCodes = 16

type registerSet [numRegisters]int

type instruction struct {
	opCode                 int
	paramA, paramB, paramC int
}

type opFunc func(paramA, paramB, paramC int, registers registerSet) registerSet

func addr(paramA, paramB, paramC int, registers registerSet) registerSet {
	registers[paramC] = registers[paramA] + registers[paramB]
	return registers
}

func addi(paramA, paramB, paramC int, registers registerSet) registerSet {
	registers[paramC] = registers[paramA] + paramB
	return registers
}

func mulr(paramA, paramB, paramC int, registers registerSet) registerSet {
	registers[paramC] = registers[paramA] * registers[paramB]
	return registers
}

func muli(paramA, paramB, paramC int, registers registerSet) registerSet {
	registers[paramC] = registers[paramA] * paramB
	return registers
}

func banr(paramA, paramB, paramC int, registers registerSet) registerSet {
	registers[paramC] = registers[paramA] & registers[paramB]
	return registers
}

func bani(paramA, paramB, paramC int, registers registerSet) registerSet {
	registers[paramC] = registers[paramA] & paramB
	return registers
}

func borr(paramA, paramB, paramC int, registers registerSet) registerSet {
	registers[paramC] = registers[paramA] | registers[paramB]
	return registers
}

func bori(paramA, paramB, paramC int, registers registerSet) registerSet {
	registers[paramC] = registers[paramA] | paramB
	return registers
}

func setr(paramA, paramB, paramC int, registers registerSet) registerSet {
	registers[paramC] = registers[paramA]
	return registers
}

func seti(paramA, paramB, paramC int, registers registerSet) registerSet {
	registers[paramC] = paramA
	return registers
}

func gtir(paramA, paramB, paramC int, registers registerSet) registerSet {
	if paramA > registers[paramB] {
		registers[paramC] = 1
	} else {
		registers[paramC] = 0
	}
	return registers
}

func gtri(paramA, paramB, paramC int, registers registerSet) registerSet {
	if registers[paramA] > paramB {
		registers[paramC] = 1
	} else {
		registers[paramC] = 0
	}
	return registers
}

func gtrr(paramA, paramB, paramC int, registers registerSet) registerSet {
	if registers[paramA] > registers[paramB] {
		registers[paramC] = 1
	} else {
		registers[paramC] = 0
	}
	return registers
}

func eqir(paramA, paramB, paramC int, registers registerSet) registerSet {
	if paramA == registers[paramB] {
		registers[paramC] = 1
	} else {
		registers[paramC] = 0
	}
	return registers
}

func eqri(paramA, paramB, paramC int, registers registerSet) registerSet {
	if registers[paramA] == paramB {
		registers[paramC] = 1
	} else {
		registers[paramC] = 0
	}
	return registers
}

func eqrr(paramA, paramB, paramC int, registers registerSet) registerSet {
	if registers[paramA] == registers[paramB] {
		registers[paramC] = 1
	} else {
		registers[paramC] = 0
	}
	return registers
}

var opFuncs = [...]opFunc{
	addr, addi,
	mulr, muli,
	banr, bani,
	borr, bori,
	setr, seti,
	gtir, gtri, gtrr,
	eqir, eqri, eqrr,
}

func getNextP1Iput(
	lines []string) ([]string, bool, registerSet, instruction, registerSet) {

	if lines[0] == "" {
		return lines[2:], false, registerSet{}, instruction{}, registerSet{}
	}

	var regsBefore, regsAfter registerSet
	var theOp instruction

	fmt.Sscanf(lines[0], "Before: [%d, %d, %d, %d]",
		&regsBefore[0], &regsBefore[1], &regsBefore[2], &regsBefore[3])

	fmt.Sscanf(lines[1], "%d %d %d %d",
		&theOp.opCode, &theOp.paramA, &theOp.paramB, &theOp.paramC)

	fmt.Sscanf(lines[2], "After:  [%d, %d, %d, %d]",
		&regsAfter[0], &regsAfter[1], &regsAfter[2], &regsAfter[3])

	return lines[4:], true, regsBefore, theOp, regsAfter
}

func main() {
	lines := utils.ReadLinesOrDie("input.txt")

	opCodeMatchCounts := map[int]int{}

	for {
		var regsBefore, regsAfter registerSet
		var theOp instruction
		var hasP1Input bool
		lines, hasP1Input, regsBefore, theOp, regsAfter = getNextP1Iput(lines)
		if !hasP1Input {
			break
		}

		matchCount := 0
		for _, opFn := range opFuncs {
			regsAfter2 := opFn(
				theOp.paramA, theOp.paramB, theOp.paramC, regsBefore)
			if regsAfter == regsAfter2 {
				matchCount++
			}
		}

		opCodeMatchCounts[matchCount]++
	}

	for i := 1; i <= numOpCodes; i++ {
		fmt.Printf("times a sample behaved like %d opcodes: %d\n",
			i, opCodeMatchCounts[i])
	}

	threeOrMore := 0
	for i := 3; i <= numOpCodes; i++ {
		threeOrMore += opCodeMatchCounts[i]
	}

	fmt.Printf("number behaving like three or more: %d\n", threeOrMore)
}
