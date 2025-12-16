package main

import (
	"bufio"
	"fmt"
	"os"
)

func run(inputPath string) (int64, error) {
	input, err := load(inputPath)
	if err != nil {
		return 0, err
	}

	res, err := process(input)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func load(inputPath string) ([][]byte, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	m := [][]byte{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := []byte{}
		for i := range len(line) {
			row = append(row, line[i])
		}
		m = append(m, row)
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("scan error: %w", scanner.Err())
	}

	return m, nil
}

func process(m [][]byte) (int64, error) {
	if len(m) <= 0 {
		return 0, nil
	}

	rowCount := len(m)
	colCount := len(m[0])
	stack := []int64{}
	res := int64(0)
	for c := colCount - 1; c >= 0; c -= 1 {
		num := int64(0)
		digitCount := int64(0)
		for r := 0; r < rowCount; r += 1 {
			b := m[r][c]

			if b >= '0' && b <= '9' {
				num = num*10 + int64(b-'0')
				digitCount += 1
				continue
			}

			// Last row?
			if r == rowCount-1 {
				// Encountered at least one digit?
				if digitCount >= 1 {
					stack = append(stack, num)
				}

				if b == '+' {
					sum := int64(0)
					for _, n := range stack {
						sum += n
					}
					stack = stack[0:0]
					res += sum
					continue
				}

				if b == '*' {
					sum := int64(1)
					for _, n := range stack {
						sum *= n
					}
					stack = stack[0:0]
					res += sum
					continue
				}
			}
		}
	}

	return res, nil
}
