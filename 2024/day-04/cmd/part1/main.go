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
	// file, err := os.Open("./assets/input_small.txt") // expected: 18
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	chars := [][]byte{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := []byte{}
		for _, c := range []byte(scanner.Text()) {
			line = append(line, c)
		}
		chars = append(chars, line)
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	count := countXmas(chars)
	fmt.Printf("%d", count)

	return nil
}

func countXmas(chars [][]byte) int {
	c := 0

	dirs := [8][2]int{
		[2]int{1, 0},
		[2]int{1, 1},
		[2]int{0, 1},
		[2]int{-1, 1},
		[2]int{-1, 0},
		[2]int{-1, -1},
		[2]int{0, -1},
		[2]int{1, -1},
	}

	for startY := range chars {
		for startX := range chars[startY] {
			for _, dir := range dirs {
				if isWord("XMAS", chars, startX, startY, dir) {
					c += 1
				}
			}
		}
	}

	return c
}

func isWord(word string, chars [][]byte, startX int, startY int, dir [2]int) bool {
	for offset := 0; offset < len(word); offset += 1 {
		y := startY + dir[1]*offset
		if y < 0 || y >= len(chars) {
			return false
		}

		x := startX + dir[0]*offset
		if x < 0 || x >= len(chars[y]) {
			return false
		}

		if chars[y][x] != word[offset] {
			return false
		}
	}

	return true
}
