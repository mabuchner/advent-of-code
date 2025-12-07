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

	m := [][]byte{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := []byte{}
		for i := range len(line) {
			b := line[i]
			if b != '.' && b != '@' {
				return nil, errors.New("unexpected character in input")
			}

			row = append(row, b)
		}
		m = append(m, row)
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("scan error: %w", scanner.Err())
	}

	return m, nil
}

func process(m [][]byte) int64 {
	offsets := [8][2]int{
		{-1, 0},
		{-1, -1},
		{0, -1},
		{1, -1},
		{1, 0},
		{1, 1},
		{0, 1},
		{-1, 1},
	}

	rows := len(m)
	cols := len(m[0])

	// Count how many rolls of paper are surrounded by fewer than 4 rolls of
	// paper
	count := int64(0)
	for y := range rows {
		for x := range cols {
			if m[y][x] != '@' {
				continue
			}

			c := 0
			for _, offset := range offsets {
				px := x + offset[0]
				py := y + offset[1]
				if px >= 0 && px < cols &&
					py >= 0 && py < rows &&
					m[py][px] == '@' {
					c += 1
				}
			}

			if c < 4 {
				count += 1
			}
		}
	}
	return count
}
