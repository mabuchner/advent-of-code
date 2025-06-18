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

	m := [][]int{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := []int{}
		s := scanner.Text()
		for i := range s {
			row = append(row, int(s[i]-'0'))
		}
		m = append(m, row)
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// fmt.Printf("%v\n", m)

	score := 0
	h := len(m)
	w := len(m[0])
	for y := 0; y < h; y += 1 {
		for x := 0; x < w; x += 1 {
			if m[y][x] == 0 {
				score += visit(m, x, y)
			}
		}
	}

	fmt.Printf("%d", score)

	return nil
}

func visit(m [][]int, startX int, startY int) int {
	var recursion func(x, y int) int

	h := len(m)
	w := len(m[0])

	dirs := [4][2]int{
		[2]int{1, 0},
		[2]int{0, 1},
		[2]int{-1, 0},
		[2]int{0, -1},
	}

	visited := [][]bool{}
	for y := 0; y < h; y += 1 {
		visited = append(visited, make([]bool, w))
	}

	recursion = func(x, y int) int {
		if m[y][x] == 9 {
			if !visited[y][x] {
				visited[y][x] = true
				return 1
			}
			return 0
		}

		score := 0
		for _, dir := range dirs {
			newX := x + dir[0]
			newY := y + dir[1]
			if newX >= 0 && newX < w && newY >= 0 && newY < h && m[newY][newX]-m[y][x] == 1 {
				score += recursion(newX, newY)
			}
		}
		return score
	}

	return recursion(startX, startY)
}
