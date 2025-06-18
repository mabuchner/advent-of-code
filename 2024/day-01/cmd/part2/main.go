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
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	nums := []int{}
	counts := map[int]int{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := 0
		r := 0
		_, err := fmt.Sscan(scanner.Text(), &l, &r)
		if err != nil {
			return err
		}

		nums = append(nums, l)
		counts[r] = counts[r] + 1
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	score := 0
	for _, n := range nums {
		c := counts[n]
		score += n * c
	}

	fmt.Printf("%d", score)

	return nil
}
