package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
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
	h := len(m)
	if h <= 0 {
		return 0, nil
	}

	w := len(m[0])
	if w <= 0 {
		return 0, nil
	}

	startIndex := slices.IndexFunc(m[0], func(b byte) bool {
		return b == 'S'
	})
	if startIndex < 0 {
		return 0, errors.New("missing start")
	}

	splitCount := int64(0)

	beams := make([]bool, w)
	beams[startIndex] = true
	for y := 1; y < h; y += 1 {
		row := m[y]

		newBeams := make([]bool, w)
		for i, b := range beams {
			if !b {
				continue
			}

			if row[i] == '^' {
				newBeams[i-1] = true
				newBeams[i+1] = true
				splitCount += 1
			} else if row[i] == '.' {
				newBeams[i] = true
			} else {
				return 0, errors.New("unexpected input byte")
			}
		}
		beams = newBeams
	}

	return splitCount, nil
}
