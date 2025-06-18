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
	// file, err := os.Open("./assets/input_tiny_tiny.txt")
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
					pos = [2]int{x * 2, y}
					ch = '.'
				}

				if ch == 'O' {
					row = append(row, '[')
					row = append(row, ']')
					continue
				}

				row = append(row, ch)
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

	// fmt.Printf("pos=%v, moves=%v\n", pos, moves)
	// printMap(m)

	for _, move := range moves {
		// tmp := m[pos[1]][pos[0]]
		// m[pos[1]][pos[0]] = '@'
		// fmt.Printf("pos=%v, move=%s\n", pos, string(move))
		// printMap(m)
		// m[pos[1]][pos[0]] = tmp

		pos = step(m, pos, move)

		// tmp = m[pos[1]][pos[0]]
		// m[pos[1]][pos[0]] = '@'
		// printMap(m)
		// m[pos[1]][pos[0]] = tmp
		// fmt.Println()
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

	if !canMove(m, newPos, dir, pos) {
		return pos
	}

	moveBox(m, newPos, dir, pos)
	m[newPos[1]][newPos[0]] = '.' // This might not be needed

	return newPos
}

func canMove(m [][]byte, pos [2]int, dir [2]int, pushPos [2]int) bool {
	tile := m[pos[1]][pos[0]]
	if tile == '#' {
		return false
	}
	if tile == '.' {
		return true
	}

	// tile=='[' || tile==']'

	newPos := [2]int{
		pos[0] + dir[0],
		pos[1] + dir[1],
	}

	if !canMove(m, newPos, dir, pos) {
		return false
	}

	// Also check other half of the box for vertical movement
	if dir[1] != 0 {
		newPos2 := pos
		if tile == '[' {
			newPos2[0] += 1
		} else { // tile == ']'
			newPos2[0] -= 1
		}

		// Do not go back to the previous tile
		if newPos2 != pushPos {
			if !canMove(m, newPos2, dir, pos) {
				return false
			}
		}
	}

	return true
}

func moveBox(m [][]byte, pos [2]int, dir [2]int, pushPos [2]int) {
	tile := m[pos[1]][pos[0]]
	if tile == '#' || tile == '.' {
		return
	}

	newPos := [2]int{
		pos[0] + dir[0],
		pos[1] + dir[1],
	}
	moveBox(m, newPos, dir, pos)

	if dir[1] != 0 {
		newPos2 := pos
		if tile == '[' {
			newPos2[0] += 1
		} else { // tile == ']'
			newPos2[0] -= 1
		}

		if newPos2 != pushPos {
			moveBox(m, newPos2, dir, pos)
		}
	}

	m[pos[1]][pos[0]] = '.'
	m[newPos[1]][newPos[0]] = tile
}

func calcScore(m [][]byte) int {
	h := len(m)
	w := len(m[0])

	score := 0
	for y := 0; y < h; y += 1 {
		for x := 0; x < w; x += 1 {
			if m[y][x] == '[' {
				score += 100*y + x
			}
		}
	}

	return score
}

func printMap(m [][]byte) {
	for _, row := range m {
		fmt.Printf("%s\n", string(row))
	}
}
