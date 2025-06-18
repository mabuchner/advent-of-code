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
	// file, err := os.Open("./assets/input_small.txt")
	// w := 7
	// h := 7
	// t := 12
	file, err := os.Open("./assets/input.txt")
	w := 71
	h := 71
	t := 1024
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	m := make([][]byte, h)
	for y := range m {
		row := make([]byte, w)
		for x := 0; x < w; x += 1 {
			row[x] = '.'
		}
		m[y] = row
	}

	coords := [][2]int64{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strs := strings.SplitN(scanner.Text(), ",", 2)

		x, err := strconv.ParseInt(strs[0], 10, 64)
		if err != nil {
			return err
		}

		y, err := strconv.ParseInt(strs[1], 10, 64)
		if err != nil {
			return err
		}

		coords = append(coords, [2]int64{x, y})
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// fmt.Printf("coords=%v\n", coords)

	for i := 0; i < t; i += 1 {
		pos := coords[i]
		m[pos[1]][pos[0]] = '#'
	}
	// printMap(m)

	minDist := findMinDist(m)
	fmt.Printf("%d", minDist)

	return nil
}

func printMap(m [][]byte) {
	for _, row := range m {
		fmt.Printf("%s\n", string(row))
	}
}

type distPos struct {
	dist int
	pos  [2]int64
}

func findMinDist(m [][]byte) int {
	h := int64(len(m))
	w := int64(len(m[0]))

	dirs := [4][2]int64{
		[2]int64{1, 0},
		[2]int64{0, 1},
		[2]int64{-1, 0},
		[2]int64{0, -1},
	}

	start := [2]int64{0, 0}
	end := [2]int64{w - 1, h - 1}

	visited := make([][]bool, h)
	for y := range visited {
		visited[y] = make([]bool, w)
	}

	queue := []distPos{distPos{0, start}}
	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		pos := item.pos
		dist := item.dist
		if pos == end {
			return dist
		}

		if pos[0] < 0 || pos[0] >= w || pos[1] < 0 || pos[1] >= h {
			continue
		}
		if m[pos[1]][pos[0]] == '#' {
			continue
		}
		if visited[pos[1]][pos[0]] {
			continue
		}

		visited[pos[1]][pos[0]] = true

		for _, dir := range dirs {
			newPos := [2]int64{
				pos[0] + dir[0],
				pos[1] + dir[1],
			}
			queue = append(queue, distPos{dist + 1, newPos})
		}
	}

	return -1
}
