package main

import (
	"bufio"
	"errors"
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
	// file, err := os.Open("./assets/input_tiny.txt")
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	m := [][]byte{}
	pos := [2]int{}
	moves := []byte{}

	scanner := bufio.NewScanner(file)
	readMoves := false
	y := 0
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			readMoves = true
			continue
		}

		if !readMoves {
			row := []byte{}
			for x := range scanner.Text() {
				ch := scanner.Text()[x]
				if ch == '@' {
					pos = [2]int{x, y}
					ch = '.'
				}
				row = append(row, ch)
			}
			m = append(m, row)
			y += 1

			continue
		}

		for i := range scanner.Text() {
			ch := scanner.Text()[i]
			moves = append(moves, ch)
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// fmt.Printf("m=%v, pos=%v, moves=%v\n", m, pos, moves)

	for _, move := range moves {
		pos = step(m, pos, move)
	}

	fmt.Printf("%d", calcScore(m))

	return nil
}

func step(m [][]byte, pos [2]int, move byte) [2]int {
	var dir [2]int
	if move == '>' {
		dir = [2]int{1, 0}
	} else if move == 'v' {
		dir = [2]int{0, 1}
	} else if move == '<' {
		dir = [2]int{-1, 0}
	} else if move == '^' {
		dir = [2]int{0, -1}
	} else {
		panic(errors.New("unexpected move"))
	}

	newPos := [2]int{
		pos[0] + dir[0],
		pos[1] + dir[1],
	}

	newTile := m[newPos[1]][newPos[0]]
	if newTile == '.' {
		return newPos
	}

	if newTile == '#' {
		return pos
	}

	// newTile == 'O'
	checkPos := newPos
	for { // This assumes that there are walls all around the map
		checkPos[0] += dir[0]
		checkPos[1] += dir[1]

		checkTile := m[checkPos[1]][checkPos[0]]
		if checkTile == '#' {
			return pos
		}

		if checkTile == '.' {
			break
		}
	}

	m[newPos[1]][newPos[0]] = '.'
	m[checkPos[1]][checkPos[0]] = 'O'

	return newPos
}

func calcScore(m [][]byte) int {
	h := len(m)
	w := len(m[0])

	score := 0
	for y := 0; y < h; y += 1 {
		for x := 0; x < w; x += 1 {
			if m[y][x] == 'O' {
				score += 100*y + x
			}
		}
	}

	return score
}
