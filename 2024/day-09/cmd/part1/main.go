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

	uncompressed := []int{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for i := range scanner.Text() {
			ch := scanner.Text()[i]
			repeat := int(ch - '0')
			empty := i%2 != 0

			content := -1
			if !empty {
				content = i / 2
			}

			for j := 0; j < repeat; j += 1 {
				uncompressed = append(uncompressed, content)
			}
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	defragmented := []int{}
	left := 0
	right := len(uncompressed) - 1
	for left < right {
		if uncompressed[left] != -1 {
			defragmented = append(defragmented, uncompressed[left])
			left += 1
			continue
		}

		if uncompressed[right] == -1 {
			right -= 1
			continue
		}

		defragmented = append(defragmented, uncompressed[right])
		left += 1
		right -= 1
	}

	// fmt.Printf("uncompressed: %v\n", uncompressed)
	// fmt.Printf("defrag: %v\n", defragmented)
	// fmt.Printf("defragStr: %v\n", defragmentedStr)

	checksum := int64(0)
	for i, n := range defragmented {
		checksum += int64(i) * int64(n)
	}

	fmt.Printf("%d", checksum)

	return nil
}
