package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("run failed: %s\n", err.Error())
		os.Exit(1)
	}
}

type readMode int

const (
	readModeInit readMode = iota
	readModeLock readMode = iota
	readModeKey  readMode = iota
)

func run() error {
	// file, err := os.Open("./assets/input_small.txt")
	file, err := os.Open("./assets/input.txt")
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	locks := make([][5]string, 0, 100)
	keys := make([][5]string, 0, 100)

	mode := readModeInit
	buf := [5]string{}
	rowIndex := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if mode == readModeInit {
			if len(line) > 0 {
				if line[0] == '#' {
					mode = readModeLock
				} else {
					mode = readModeKey
				}
			}
		} else if mode == readModeLock {
			if rowIndex < 5 {
				buf[rowIndex] = line
				rowIndex += 1
			} else {
				locks = append(locks, buf)
				buf = [5]string{}
				rowIndex = 0
				mode = readModeInit
			}
		} else if mode == readModeKey {
			if rowIndex < 5 {
				buf[rowIndex] = line
				rowIndex += 1
			} else {
				keys = append(keys, buf)
				buf = [5]string{}
				rowIndex = 0
				mode = readModeInit
			}
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// fmt.Printf("locks=%v\n", locks)
	// fmt.Printf("keys=%v\n", keys)

	lockHeights := calcAllHeight(locks)
	keyHeights := calcAllHeight(keys)
	// fmt.Printf("lockHeights=%v\n", lockHeights)
	// fmt.Printf("keyHeights=%v\n", keyHeights)

	result := 0
	for _, lock := range lockHeights {
		for _, key := range keyHeights {
			if doesFit(lock, key) {
				result += 1
			}
		}
	}
	fmt.Printf("%d", result)

	return nil
}

func calcAllHeight(bufs [][5]string) [][5]int {
	allHeights := [][5]int{}
	for _, buf := range bufs {
		heights := calcHeights(buf)
		allHeights = append(allHeights, heights)
	}
	return allHeights
}

func calcHeights(buf [5]string) [5]int {
	heights := [5]int{}
	for _, line := range buf {
		for col, ch := range line {
			if ch == '#' {
				heights[col] += 1
			}
		}
	}
	return heights
}

func doesFit(lockHeights, keyHeights [5]int) bool {
	for col := 0; col < 5; col += 1 {
		sum := lockHeights[col] + keyHeights[col]
		if sum > 5 {
			return false
		}
	}
	return true
}
