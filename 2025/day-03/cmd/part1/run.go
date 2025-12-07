package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func run(inputPath string) (int64, error) {
	input, err := load(inputPath)
	if err != nil {
		return 0, err
	}
	return process(input), nil
}

func load(inputPath string) ([][]byte, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	banks := [][]byte{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		bank := []byte{}
		for i := range len(line) {
			if line[i] < '0' || line[i] > '9' {
				return nil, errors.New("unexpected character in input")
			}
			n := line[i] - '0'
			bank = append(bank, n)
		}
		banks = append(banks, bank)
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("scan error: %w", scanner.Err())
	}

	return banks, nil
}

func process(banks [][]byte) int64 {
	sum := int64(0)

	for _, bank := range banks {
		// Find the largest number, leaving at least on number on the right
		l := byte(0)
		lIndex := 0
		for i := 0; i < len(bank)-1; i += 1 {
			if bank[i] > l {
				lIndex = i
				l = bank[i]
			}
		}

		// Find the largest number after the first number
		r := byte(0)
		for i := lIndex + 1; i < len(bank); i += 1 {
			r = max(r, bank[i])
		}

		sum += int64(10*l + r)
	}

	return sum
}
