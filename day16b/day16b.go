package main

import (
	"fmt"
	"sort"

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

type opFuncDesc struct {
	op   opFunc
	name string
}

var opFuncDescs = [...]opFuncDesc{
	opFuncDesc{addr, "addr"},
	opFuncDesc{addi, "addi"},
	opFuncDesc{mulr, "mulr"},
	opFuncDesc{muli, "muli"},
	opFuncDesc{banr, "banr"},
	opFuncDesc{bani, "bani"},
	opFuncDesc{borr, "borr"},
	opFuncDesc{bori, "bori"},
	opFuncDesc{setr, "setr"},
	opFuncDesc{seti, "seti"},
	opFuncDesc{gtir, "gtir"},
	opFuncDesc{gtri, "gtri"},
	opFuncDesc{gtrr, "gtrr"},
	opFuncDesc{eqir, "eqir"},
	opFuncDesc{eqri, "eqri"},
	opFuncDesc{eqrr, "eqrr"},
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

func doMatchingCycle(
	lines []string, opCodes map[int]int, knownOpFuncs map[int]bool) []string {

	opCodeMatchCounts := map[int]int{}

	for {
		var regsBefore, regsAfter registerSet
		var theOp instruction
		var hasP1Input bool
		lines, hasP1Input, regsBefore, theOp, regsAfter = getNextP1Iput(lines)
		if !hasP1Input {
			break
		}

		if _, exists := opCodes[theOp.opCode]; exists {
			continue
		}

		matchCount := 0
		var matchingOpFuncIdx int
		for i := range opFuncDescs {
			if knownOpFuncs[i] {
				continue
			}

			regsAfter2 := opFuncDescs[i].op(
				theOp.paramA, theOp.paramB, theOp.paramC, regsBefore)
			if regsAfter == regsAfter2 {
				matchCount++
				matchingOpFuncIdx = i
			}
		}

		// fmt.Printf("matchCount: %d\n", matchCount)

		opCodeMatchCounts[matchCount]++

		if matchCount == 1 {
			opCodes[theOp.opCode] = matchingOpFuncIdx
			knownOpFuncs[matchingOpFuncIdx] = true
			// fmt.Printf("  storing matchingOpFuncIdx: %d\n", matchingOpFuncIdx)
		}
	}

	fmt.Println()
	for i := 1; i <= numOpCodes; i++ {
		fmt.Printf("times a sample behaved like %d opcodes: %d\n",
			i, opCodeMatchCounts[i])
	}

	matchedOpCodes := make([]int, 0, len(opCodes))
	for key := range opCodes {
		matchedOpCodes = append(matchedOpCodes, key)
	}
	sort.Ints(matchedOpCodes)
	fmt.Printf("\nMatched opcodes so far: %v\n", matchedOpCodes)

	return lines
}

func checkMatching(
	lines []string, opCodes map[int]int, knownOpFuncs map[int]bool) []string {

	inputNum := 1

	for {
		var regsBefore, regsAfter registerSet
		var theOp instruction
		var hasP1Input bool
		lines, hasP1Input, regsBefore, theOp, regsAfter = getNextP1Iput(lines)
		if !hasP1Input {
			break
		}

		opCodeDescIdx := opCodes[theOp.opCode]
		regsAfter2 := opFuncDescs[opCodeDescIdx].op(
			theOp.paramA, theOp.paramB, theOp.paramC, regsBefore)
		if regsAfter != regsAfter2 {
			fmt.Printf("op DID NOT MATCH for input num: %d\n", inputNum)
		}

		inputNum++
	}

	return lines
}

func main() {
	lines := utils.ReadLinesOrDie("input.txt")

	opCodes := map[int]int{}
	knownOpFuncs := map[int]bool{}

	for len(opCodes) < numOpCodes {
		doMatchingCycle(lines, opCodes, knownOpFuncs)
	}

	remainingLines := checkMatching(lines, opCodes, knownOpFuncs)

	fmt.Println()
	for i := 0; i < len(knownOpFuncs); i++ {
		fmt.Printf("opcode %2d: %s\n", i, opFuncDescs[opCodes[i]].name)
	}

	var testRegisters registerSet
	for _, line := range remainingLines {
		var instr instruction
		fmt.Sscanf(line, "%d %d %d %d",
			&instr.opCode, &instr.paramA, &instr.paramB, &instr.paramC)

		opCodeDescIdx := opCodes[instr.opCode]
		testRegisters = opFuncDescs[opCodeDescIdx].op(
			instr.paramA, instr.paramB, instr.paramC, testRegisters)
	}

	fmt.Println()
	fmt.Printf("test registers: %v\n", testRegisters)
}
