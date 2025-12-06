package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

const digitCount = 100

func run(inputPath string) (int64, error) {
	input, err := load(inputPath)
	if err != nil {
		return 0, err
	}
	return process(input), nil
}

func load(inputPath string) ([]int64, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	dists := make([]int64, 0, 5000)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 2 {
			return nil, errors.New("line too short")
		}

		dir := line[0]
		if dir != 'L' && dir != 'R' {
			return nil, errors.New("unexpected direction")
		}

		dist, err := strconv.ParseInt(line[1:], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse distance: %w", err)
		}

		if dir == 'L' {
			dist = -dist
		}

		dists = append(dists, dist)
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("scan error: %w", scanner.Err())
	}

	return dists, nil
}

func process(dists []int64) int64 {
	pos := int64(digitCount / 2)
	count := int64(0)

	for _, dist := range dists {
		pos += dist

		// Adding digitCount to make sure, that "underflow" is correctly handled
		pos += digitCount

		// Modulo digitCount to make sure, that "overflow" is correctly handled
		pos %= digitCount

		if pos == 0 {
			count += 1
		}
	}

	return count
}
