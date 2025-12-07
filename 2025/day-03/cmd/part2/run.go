package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

const batterySelectionCount = 12

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
	total := int64(0)

	for _, bank := range banks {
		// Select the battery with the highest energy while leaving space
		// for selecting the remaining batteries

		n := batterySelectionCount // How many batteries left to process
		prevIndex := -1            // Index of the previously selected battery
		sum := int64(0)            // Current energy level
		for n > 0 {
			n -= 1
			biggest := byte(0)
			biggestIndex := 0
			for i := prevIndex + 1; i < len(bank)-n; i += 1 {
				if bank[i] > biggest {
					biggest = bank[i]
					biggestIndex = i
				}
			}
			prevIndex = biggestIndex
			sum = sum*10 + int64(biggest)
		}

		total += sum
	}

	return total
}
