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
	// file, err := os.Open("./assets/input_small.txt") // expected: 9
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

	templates := [][]string{
		[]string{
			"M.S",
			".A.",
			"M.S",
		},
		[]string{
			"S.S",
			".A.",
			"M.M",
		},
		[]string{
			"S.M",
			".A.",
			"S.M",
		},
		[]string{
			"M.M",
			".A.",
			"S.S",
		},
	}

	for startY := range chars {
		for startX := range chars[startY] {
			for _, template := range templates {
				if matches(chars, startX, startY, template) {
					c += 1
				}
			}
		}
	}

	return c
}

func matches(chars [][]byte, startX int, startY int, template []string) bool {
	th := len(template)
	if th <= 0 {
		panic("unexpected template height")
	}

	tw := len(template[0])
	if tw <= 0 {
		panic("unexpected template width")
	}

	if startY+th > len(chars) {
		return false
	}
	if startX+tw > len(chars[0]) {
		return false
	}

	for offsetY := range template {
		for offsetX := range template[offsetY] {
			t := template[offsetY][offsetX]
			c := chars[startY+offsetY][startX+offsetX]
			if t != '.' && c != t {
				return false
			}
		}
	}

	return true
}
