package main

import (
	"fmt"

	"jblee.net/adventofcode2018/utils"
)

const numRegisters = 6

type registerSet [numRegisters]int

type instruction struct {
	opCode                 string
	paramA, paramB, paramC int
}

type program struct {
	ipRegister   int
	instructions []instruction
}

type virtMachine struct {
	instrPtr int
	regs     registerSet
}

type opFunc func(paramA, paramB, paramC int, registers *registerSet)

func addr(paramA, paramB, paramC int, registers *registerSet) {
	registers[paramC] = registers[paramA] + registers[paramB]
}

func addi(paramA, paramB, paramC int, registers *registerSet) {
	registers[paramC] = registers[paramA] + paramB
}

func mulr(paramA, paramB, paramC int, registers *registerSet) {
	registers[paramC] = registers[paramA] * registers[paramB]
}

func muli(paramA, paramB, paramC int, registers *registerSet) {
	registers[paramC] = registers[paramA] * paramB
}

func banr(paramA, paramB, paramC int, registers *registerSet) {
	registers[paramC] = registers[paramA] & registers[paramB]
}

func bani(paramA, paramB, paramC int, registers *registerSet) {
	registers[paramC] = registers[paramA] & paramB
}

func borr(paramA, paramB, paramC int, registers *registerSet) {
	registers[paramC] = registers[paramA] | registers[paramB]
}

func bori(paramA, paramB, paramC int, registers *registerSet) {
	registers[paramC] = registers[paramA] | paramB
}

func setr(paramA, paramB, paramC int, registers *registerSet) {
	registers[paramC] = registers[paramA]
}

func seti(paramA, paramB, paramC int, registers *registerSet) {
	registers[paramC] = paramA
}

func gtir(paramA, paramB, paramC int, registers *registerSet) {
	if paramA > registers[paramB] {
		registers[paramC] = 1
	} else {
		registers[paramC] = 0
	}
}

func gtri(paramA, paramB, paramC int, registers *registerSet) {
	if registers[paramA] > paramB {
		registers[paramC] = 1
	} else {
		registers[paramC] = 0
	}
}

func gtrr(paramA, paramB, paramC int, registers *registerSet) {
	if registers[paramA] > registers[paramB] {
		registers[paramC] = 1
	} else {
		registers[paramC] = 0
	}
}

func eqir(paramA, paramB, paramC int, registers *registerSet) {
	if paramA == registers[paramB] {
		registers[paramC] = 1
	} else {
		registers[paramC] = 0
	}
}

func eqri(paramA, paramB, paramC int, registers *registerSet) {
	if registers[paramA] == paramB {
		registers[paramC] = 1
	} else {
		registers[paramC] = 0
	}
}

func eqrr(paramA, paramB, paramC int, registers *registerSet) {
	if registers[paramA] == registers[paramB] {
		registers[paramC] = 1
	} else {
		registers[paramC] = 0
	}
}

var opFuncs = map[string]opFunc{
	"addr": addr,
	"addi": addi,
	"mulr": mulr,
	"muli": muli,
	"banr": banr,
	"bani": bani,
	"borr": borr,
	"bori": bori,
	"setr": setr,
	"seti": seti,
	"gtir": gtir,
	"gtri": gtri,
	"gtrr": gtrr,
	"eqir": eqir,
	"eqri": eqri,
	"eqrr": eqrr,
}

func parseProgram(lines []string) program {
	var prog program

	if len(lines) > 0 && lines[0][0] == '#' {
		fmt.Sscanf(lines[0], "#ip %d", &prog.ipRegister)
		lines = lines[1:]
	} else {
		prog.ipRegister = -1
	}

	prog.instructions = make([]instruction, len(lines))
	for i, line := range lines {
		instr := &prog.instructions[i]
		fmt.Sscanf(line, "%s %d %d %d",
			&instr.opCode, &instr.paramA, &instr.paramB, &instr.paramC)
	}

	return prog
}

func main() {
	lines := utils.ReadLinesOrDie("input.txt")
	// lines := utils.ReadLinesOrDie("sample_input1.txt")

	var vm virtMachine
	prog := parseProgram(lines)

	for {
		if vm.instrPtr < 0 || vm.instrPtr >= len(prog.instructions) {
			break
		}

		if prog.ipRegister > -1 {
			vm.regs[prog.ipRegister] = vm.instrPtr
		}

		instr := &prog.instructions[vm.instrPtr]

		// fmt.Printf("ip=%d %v %s %d %d %d", vm.instrPtr, vm.regs,
		// 	instr.opCode, instr.paramA, instr.paramB, instr.paramC)

		opFuncs[instr.opCode](
			instr.paramA, instr.paramB, instr.paramC, &vm.regs)

		// fmt.Printf(" %v\n", vm.regs)

		if prog.ipRegister > -1 {
			vm.instrPtr = vm.regs[prog.ipRegister]
		}

		vm.instrPtr++
	}

	fmt.Println(vm.regs)
}
