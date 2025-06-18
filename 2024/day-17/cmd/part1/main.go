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
	// file, err := os.Open("./assets/input_small.txt") // 4,6,3,5,6,3,5,2,1,0
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
			cpu.registers[mode] = int(num)
		case 3: // Skip empty line
		case 4: // Read program instructions
			instrStr := scanner.Text()[9:]
			instrStrs := strings.Split(instrStr, ",")
			for _, s := range instrStrs {
				num, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return err
				}
				cpu.instructions = append(cpu.instructions, int(num))
			}
		}

		mode += 1
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// fmt.Printf("cpu=%+v\n", cpu)

	for cpu.step() {
		// fmt.Printf("pc=%d\n", cpu.pc)
	}

	outStrs := make([]string, len(cpu.outBuf))
	for i, o := range cpu.outBuf {
		outStrs[i] = strconv.FormatInt(int64(o), 10)
	}
	fmt.Printf("%s", strings.Join(outStrs, ","))

	return nil
}

type instrFunc func(operand int)

type cpu struct {
	registers    [3]int
	pc           int
	instructions []int
	outBuf       []int
	lookup       map[int]instrFunc
}

func (c *cpu) init() {
	c.lookup = make(map[int]instrFunc, 8)
	c.lookup[0] = c.adv
	c.lookup[1] = c.bxl
	c.lookup[2] = c.bst
	c.lookup[3] = c.jnz
	c.lookup[4] = c.bxc
	c.lookup[5] = c.out
	c.lookup[6] = c.bdv
	c.lookup[7] = c.cdv
}

func (c *cpu) step() bool {
	opcode := c.instructions[c.pc]
	operand := c.instructions[c.pc+1]
	c.lookup[opcode](operand)
	return c.pc < len(c.instructions)
}

func (c *cpu) adv(operand int) {
	numerator := c.regA()
	denominator := 1 << c.combo(operand)
	c.setRegA(numerator / denominator)
	c.pc += 2
}

func (c *cpu) bxl(operand int) {
	c.setRegB(c.regB() ^ operand)
	c.pc += 2
}

func (c *cpu) bst(operand int) {
	c.setRegB(c.combo(operand) % 8)
	c.pc += 2
}

func (c *cpu) jnz(operand int) {
	if c.regA() != 0 {
		c.pc = operand
	} else {
		c.pc += 2
	}
}

func (c *cpu) bxc(operand int) {
	c.setRegB(c.regB() ^ c.regC())
	c.pc += 2
}

func (c *cpu) out(operand int) {
	c.outBuf = append(c.outBuf, c.combo(operand)%8)
	c.pc += 2
}

func (c *cpu) bdv(operand int) {
	numerator := c.regA()
	denominator := 1 << c.combo(operand)
	c.setRegB(numerator / denominator)
	c.pc += 2
}

func (c *cpu) cdv(operand int) {
	numerator := c.regA()
	denominator := 1 << c.combo(operand)
	c.setRegC(numerator / denominator)
	c.pc += 2
}

func (c *cpu) combo(operand int) int {
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

func (c *cpu) regA() int {
	return c.registers[0]
}

func (c *cpu) setRegA(val int) {
	c.registers[0] = val
}

func (c *cpu) regB() int {
	return c.registers[1]
}

func (c *cpu) setRegB(val int) {
	c.registers[1] = val
}

func (c *cpu) regC() int {
	return c.registers[2]
}

func (c *cpu) setRegC(val int) {
	c.registers[2] = val
}
