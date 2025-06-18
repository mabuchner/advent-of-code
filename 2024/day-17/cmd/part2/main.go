package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("run failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	// file, err := os.Open("./assets/input_small2.txt") // 117440
	file, err := os.Open("./assets/input.txt")
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	cpu := cpu{}
	cpu.init()

	mode := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		switch mode {
		case 0, 1, 2: // Read register value
			numStr := scanner.Text()[12:]
			num, err := strconv.ParseInt(numStr, 10, 64)
			if err != nil {
				return err
			}
			cpu.registers[mode] = num
		case 3: // Skip empty line
		case 4: // Read program instructions
			instrStr := scanner.Text()[9:]
			instrStrs := strings.Split(instrStr, ",")
			for _, s := range instrStrs {
				num, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return err
				}
				cpu.instructions = append(cpu.instructions, num)
			}
		}

		mode += 1
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// fmt.Printf("cpu=%+v\n", cpu)
	// fmt.Println(disassemble(cpu.instructions))

	// Brute force last opcode
	solutions := []int64{0}
	for i := len(cpu.instructions) - 1; i >= 0; i -= 1 {
		targetInstruction := cpu.instructions[i]
		newSolutions := []int64{}
		for _, oldA := range solutions {
			for n := int64(0); n <= 7; n += 1 {
				newA := (oldA << 3) | n

				cpu.reset()
				cpu.setRegA(newA)
				cpu.run()

				isSolution := cpu.outBuf[0] == targetInstruction
				if isSolution {
					newSolutions = append(newSolutions, newA)
				}

				// fmt.Printf(
				// 	"n=%d (0b%s), newA=%d (0x%x, 0b%s), target=%d, outBuf=%v, isSolution=%t\n",
				// 	n,
				// 	strconv.FormatInt(n, 2),
				// 	newA,
				// 	newA,
				// 	strconv.FormatInt(newA, 2),
				// 	targetInstruction,
				// 	cpu.outBuf,
				// 	isSolution,
				// )
			}
		}

		solutions = newSolutions

		// fmt.Printf("solutions=%v\n", solutions)

		if len(solutions) == 0 {
			return fmt.Errorf("no solution found")
		}
	}

	fmt.Printf("%d", solutions[0])

	return nil
}

type instrFunc func(operand int64)

type cpu struct {
	registers    [3]int64
	pc           int64
	instructions []int64
	outBuf       []int64
	lookup       map[int64]instrFunc
}

func (c *cpu) init() {
	c.lookup = make(map[int64]instrFunc, 8)
	c.lookup[0] = c.adv
	c.lookup[1] = c.bxl
	c.lookup[2] = c.bst
	c.lookup[3] = c.jnz
	c.lookup[4] = c.bxc
	c.lookup[5] = c.out
	c.lookup[6] = c.bdv
	c.lookup[7] = c.cdv
}

func (c *cpu) reset() {
	for i := range c.registers {
		c.registers[i] = 0
	}
	c.pc = 0
	c.outBuf = nil
}

func (c *cpu) run() {
	for c.step() {
	}
}

func (c *cpu) step() bool {
	opcode := c.instructions[c.pc]
	operand := c.instructions[c.pc+1]
	c.lookup[opcode](operand)
	return c.pc < int64(len(c.instructions))
}

func (c *cpu) adv(operand int64) {
	numerator := c.regA()
	denominator := int64(1) << c.combo(operand)
	c.setRegA(numerator / denominator)
	c.pc += 2
}

func (c *cpu) bxl(operand int64) {
	c.setRegB(c.regB() ^ operand)
	c.pc += 2
}

func (c *cpu) bst(operand int64) {
	c.setRegB(c.combo(operand) % 8)
	c.pc += 2
}

func (c *cpu) jnz(operand int64) {
	if c.regA() != 0 {
		c.pc = operand
	} else {
		c.pc += 2
	}
}

func (c *cpu) bxc(operand int64) {
	c.setRegB(c.regB() ^ c.regC())
	c.pc += 2
}

func (c *cpu) out(operand int64) {
	c.outBuf = append(c.outBuf, c.combo(operand)%8)
	c.pc += 2
}

func (c *cpu) bdv(operand int64) {
	numerator := c.regA()
	denominator := int64(1) << c.combo(operand)
	c.setRegB(numerator / denominator)
	c.pc += 2
}

func (c *cpu) cdv(operand int64) {
	numerator := c.regA()
	denominator := int64(1) << c.combo(operand)
	c.setRegC(numerator / denominator)
	c.pc += 2
}

func (c *cpu) combo(operand int64) int64 {
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case 4:
		return c.regA()
	case 5:
		return c.regB()
	case 6:
		return c.regC()
	}
	panic("invalid combo operand value")
}

func (c *cpu) regA() int64 {
	return c.registers[0]
}

func (c *cpu) setRegA(val int64) {
	c.registers[0] = val
}

func (c *cpu) regB() int64 {
	return c.registers[1]
}

func (c *cpu) setRegB(val int64) {
	c.registers[1] = val
}

func (c *cpu) regC() int64 {
	return c.registers[2]
}

func (c *cpu) setRegC(val int64) {
	c.registers[2] = val
}

func disassemble(instructions []int64) string {
	buf := strings.Builder{}
	for pc := 0; pc < len(instructions); pc += 2 {
		opcode := instructions[pc]
		operand := instructions[pc+1]
		buf.WriteString(fmt.Sprintf("%5d:  %d,%d  %s\n", pc, opcode, operand, opString(opcode, operand)))
	}
	return buf.String()
}

func opString(opcode int64, operand int64) string {
	switch opcode {
	case 0:
		return fmt.Sprintf("adv %s   ; A=A>>%s", comboString(operand), comboString(operand))
	case 1:
		return fmt.Sprintf("bxl %d   ; B=B xor %d", operand, operand)
	case 2:
		return fmt.Sprintf("bst %s   ; B=%s&0b111", comboString(operand), comboString(operand))
	case 3:
		return fmt.Sprintf("jnz %-3d ; A!=0 ? PC=%d", operand, operand)
	case 4:
		return fmt.Sprintf("bxc     ; B=B xor C")
	case 5:
		return fmt.Sprintf("out %s&0b111", comboString(operand))
	case 6:
		return fmt.Sprintf("bdv %s   ; B=A>>%s", comboString(operand), comboString(operand))
	case 7:
		return fmt.Sprintf("cdv %s   ; C=A>>%s", comboString(operand), comboString(operand))
	}
	return ""
}

func comboString(operand int64) string {
	switch operand {
	case 0, 1, 2, 3:
		return strconv.FormatInt(operand, 10)
	case 4:
		return "A"
	case 5:
		return "B"
	case 6:
		return "C"
	}
	return "INVALID"
}
