package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("run failed: %s", err.Error())
		os.Exit(1)
	}
}

func run() error {
	file, err := os.Open("./assets/input.txt")
	// file, err := os.Open("./assets/input_small.txt")
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	m := [][]byte{}
	start := [2]int{-1, -1}
	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		row := []byte{}
		line := scanner.Text()
		for x := range line {
			b := line[x]
			if b != '^' {
				row = append(row, b)
			} else {
				start[0] = x
				start[1] = y
				row = append(row, 'X')
			}
		}
		m = append(m, row)
		y += 1
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// fmt.Printf("map=%v\n", m)
	// fmt.Printf("start=%v\n", start)

	dirs := [4][2]int{
		[2]int{0, -1}, // Up
		[2]int{1, 0},  // Right
		[2]int{0, 1},  // Down
		[2]int{-1, 0}, // Left
	}

	dirIndex := 0

	pos := [2]int{start[0], start[1]}

	h := len(m)
	w := len(m[0])
	for {
		dir := dirs[dirIndex]

		newPos := [2]int{
			pos[0] + dir[0],
			pos[1] + dir[1],
		}

		if newPos[0] < 0 || newPos[1] < 0 || newPos[0] >= w || newPos[1] >= h {
			break
		}

		if m[newPos[1]][newPos[0]] != '#' {
			m[newPos[1]][newPos[0]] = 'X'
			pos = newPos
		} else {
			dirIndex = (dirIndex + 1) % len(dirs)
		}
	}

	count := 0
	for _, row := range m {
		for _, b := range row {
			if b == 'X' {
				count += 1
			}
		}
	}
	fmt.Printf("%d", count)

	return nil
}
