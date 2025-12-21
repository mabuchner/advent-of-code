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

	// Use memoization to avoid repeating the same timeline calculation more
	// than once
	memo := make([][]int64, h)
	for y := range h {
		r := make([]int64, w)
		for x := range w {
			r[x] = -1 // Use -1 to indicate an undefined number of timelines
		}
		memo[y] = r
	}

	// Whenever the beam reaches the last row, it's one timeline
	last := memo[h-1]
	for x := range w {
		last[x] = 1
	}

	var visit func(x, y int) int64
	visit = func(x, y int) int64 {
		if x < 0 || x >= w || y < 0 || y >= h {
			return 0
		}

		// Position already visited (or last row)?
		if memo[y][x] != -1 {
			return memo[y][x]
		}

		t := int64(0)
		b := m[y][x]
		if b == '.' {
			t = visit(x, y+1)
		} else if b == '^' {
			t = visit(x-1, y+1) + visit(x+1, y+1)
		}
		memo[y][x] = t

		return t
	}

	return visit(startIndex, 1), nil
}
