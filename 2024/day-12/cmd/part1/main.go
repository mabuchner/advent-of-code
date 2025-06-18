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

	m := [][]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := []string{}
		for _, c := range scanner.Text() {
			row = append(row, string(c))
		}
		m = append(m, row)
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// fmt.Printf("%v\n", m)

	h := len(m)
	w := len(m[0])
	visited := make([][]bool, h)
	for y := 0; y < h; y += 1 {
		visited[y] = make([]bool, w)
	}

	sum := 0
	for y := 0; y < h; y += 1 {
		for x := 0; x < w; x += 1 {
			if !visited[y][x] {
				area, perimeter := visit(m, x, y, visited)
				sum += area * perimeter
			}
		}
	}

	fmt.Printf("%d", sum)

	return nil
}

func visit(m [][]string, startX, startY int, vis [][]bool) (area int, perimeter int) {
	h := len(m)
	w := len(m[0])
	visited := make([][]bool, h)
	for y := 0; y < h; y += 1 {
		visited[y] = make([]bool, w)
	}

	s := m[startY][startX]

	dirs := [4][2]int{
		[2]int{1, 0},
		[2]int{0, 1},
		[2]int{-1, 0},
		[2]int{0, -1},
	}

	var recusion func(x, y int) int

	recusion = func(x, y int) int {
		if x < 0 || x >= w || y < 0 || y >= h {
			return 0
		}

		if m[y][x] != s {
			return 0
		}

		if visited[y][x] {
			return 0
		}

		visited[y][x] = true
		vis[y][x] = true

		a := 1
		for _, dir := range dirs {
			a += recusion(x+dir[0], y+dir[1])
		}
		return a
	}

	area = recusion(startX, startY)

	perimeter = 0
	for y := 0; y < h; y += 1 {
		for x := 0; x < w; x += 1 {
			if !visited[y][x] {
				continue
			}

			p := 4
			for _, dir := range dirs {
				nx := x + dir[0]
				ny := y + dir[1]
				if nx >= 0 && nx < w && ny >= 0 && ny < h && visited[ny][nx] {
					p -= 1
				}
			}
			perimeter += p
		}
	}

	return area, perimeter
}
