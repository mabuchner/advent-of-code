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

const (
	bitmaskUp = 1 << iota
	bitmaskRight
	bitmaskDown
	bitmaskLeft
	bitmaskObstacle
)

const (
	bitmaskAllDirs = bitmaskUp | bitmaskRight | bitmaskDown | bitmaskLeft
)

func run() error {
	file, err := os.Open("./assets/input.txt")
	// file, err := os.Open("./assets/input_small.txt")
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	m := [][]int{}
	start := [2]int{-1, -1}
	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		row := []int{}
		line := scanner.Text()
		for x := range line {
			b := line[x]
			if b == '^' {
				start[0] = x
				start[1] = y
				row = append(row, bitmaskUp)
			} else if b == '#' {
				row = append(row, bitmaskObstacle)
			} else {
				row = append(row, 0)
			}
		}
		m = append(m, row)
		y += 1
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	visitedPositions := getVisited(m, start)

	loopCount := 0
	for _, pos := range visitedPositions {
		if hasLoop(m, start, pos) {
			loopCount += 1
		}
	}

	fmt.Printf("%v\n", loopCount)

	return nil
}

func getVisited(m [][]int, start [2]int) [][2]int {
	newMap := make([][]int, 0, len(m))
	for _, row := range m {
		newRow := make([]int, 0, len(row))
		for _, n := range row {
			newRow = append(newRow, n)
		}
		newMap = append(newMap, newRow)
	}

	dirs := [4][2]int{
		[2]int{0, -1}, // Up
		[2]int{1, 0},  // Right
		[2]int{0, 1},  // Down
		[2]int{-1, 0}, // Left
	}

	dirBitmasks := [4]int{
		bitmaskUp,
		bitmaskRight,
		bitmaskDown,
		bitmaskLeft,
	}

	dirIndex := 0

	pos := [2]int{start[0], start[1]}

	h := len(newMap)
	w := len(newMap[0])
	for {
		dir := dirs[dirIndex]

		newPos := [2]int{
			pos[0] + dir[0],
			pos[1] + dir[1],
		}

		if newPos[0] < 0 || newPos[1] < 0 || newPos[0] >= w || newPos[1] >= h {
			break
		}

		if newMap[newPos[1]][newPos[0]] != bitmaskObstacle {
			newMap[newPos[1]][newPos[0]] |= dirBitmasks[dirIndex]
			pos = newPos
		} else {
			dirIndex = (dirIndex + 1) % len(dirs)
		}
	}

	// fmt.Printf("newMap=%v\n", newMap)

	visited := [][2]int{}
	for y, row := range newMap {
		for x, n := range row {
			if (n & bitmaskAllDirs) != 0 {
				visited = append(visited, [2]int{x, y})
			}
		}
	}
	return visited
}

func hasLoop(m [][]int, start [2]int, obstaclePos [2]int) bool {
	newMap := make([][]int, 0, len(m))
	for _, row := range m {
		newRow := make([]int, 0, len(row))
		for _, n := range row {
			newRow = append(newRow, n)
		}
		newMap = append(newMap, newRow)
	}

	newMap[obstaclePos[1]][obstaclePos[0]] = bitmaskObstacle
	// fmt.Printf("map=%v\n", m)

	dirs := [4][2]int{
		[2]int{0, -1}, // Up
		[2]int{1, 0},  // Right
		[2]int{0, 1},  // Down
		[2]int{-1, 0}, // Left
	}

	dirBitmasks := [4]int{
		bitmaskUp,
		bitmaskRight,
		bitmaskDown,
		bitmaskLeft,
	}

	dirIndex := 0

	pos := [2]int{start[0], start[1]}

	h := len(newMap)
	w := len(newMap[0])
	for {
		dir := dirs[dirIndex]

		newPos := [2]int{
			pos[0] + dir[0],
			pos[1] + dir[1],
		}

		if newPos[0] < 0 || newPos[1] < 0 || newPos[0] >= w || newPos[1] >= h {
			break
		}

		n := newMap[newPos[1]][newPos[0]]
		if newMap[newPos[1]][newPos[0]] != bitmaskObstacle {
			dirBitmask := dirBitmasks[dirIndex]
			if (n & dirBitmask) != 0 { // Already visited with same direction
				return true
			}
			newMap[newPos[1]][newPos[0]] |= dirBitmask
			pos = newPos
		} else {
			dirIndex = (dirIndex + 1) % len(dirs)
			newMap[pos[1]][pos[0]] |= dirBitmasks[dirIndex]
		}
	}

	// fmt.Printf("newMap=%v\n", newMap)

	return false
}
