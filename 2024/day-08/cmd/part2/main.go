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
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := []byte{}
		for i := range scanner.Text() {
			row = append(row, scanner.Text()[i])
		}
		m = append(m, row)
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	// for _, row := range m {
	// 	fmt.Printf("%v\n", row)
	// }

	h := len(m)
	w := len(m[0])

	antinodes := [][]int{}
	for y := 0; y < h; y += 1 {
		antinodes = append(antinodes, make([]int, len(m[y])))
	}

	for y := 0; y < h; y += 1 {
		for x := 0; x < w; x += 1 {
			if m[y][x] == '.' {
				continue
			}
			for yy := 0; yy < h; yy += 1 {
				for xx := 0; xx < w; xx += 1 {
					if x == xx && y == yy {
						continue
					}

					if m[y][x] == m[yy][xx] {
						dX := xx - x
						dY := yy - y
						xxx := x + dX
						yyy := y + dY
						for xxx >= 0 && xxx < w && yyy >= 0 && yyy < h {
							antinodes[yyy][xxx] = 1
							xxx += dX
							yyy += dY
						}
					}
				}
			}
		}
	}

	// for _, row := range antinodes {
	// 	fmt.Printf("%v\n", row)
	// }

	sum := 0
	for _, row := range antinodes {
		for _, n := range row {
			sum += n
		}
	}
	fmt.Printf("%d", sum)

	return nil
}
