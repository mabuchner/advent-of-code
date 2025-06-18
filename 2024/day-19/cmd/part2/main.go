package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("run failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	// file, err := os.Open("./assets/input_small.txt") // 16
	file, err := os.Open("./assets/input.txt") // 880877787214477
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	towels := []string{}
	designs := []string{}
	mode := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if mode == 0 { // Read available towels
			towels = strings.Split(scanner.Text(), ", ")
			mode += 1
		} else if mode == 1 { // Skip empty line
			mode += 1
		} else if mode == 2 { // Read target designs
			designs = append(designs, scanner.Text())
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// fmt.Printf("towels=%v, designs=%v\n", towels, designs)

	count := countPossibleArrangements(towels, designs)
	fmt.Printf("%d", count)

	return nil
}

func countPossibleArrangements(towels []string, designs []string) int {
	towelSet := make(map[string]struct{}, len(towels))
	for _, towel := range towels {
		towelSet[towel] = struct{}{}
	}

	memo := map[string]int{}
	memo[""] = 1

	var countArrangements func(design string) int

	countArrangements = func(design string) int {
		if combinations, ok := memo[design]; ok {
			return combinations
		}

		count := 0
		for i := 1; i <= len(design); i += 1 {
			prefix := design[0:i]
			if _, ok := towelSet[prefix]; ok {
				suffix := design[i:]
				count += countArrangements(suffix)
			}
		}
		memo[design] = count
		return count
	}

	count := 0
	for _, design := range designs {
		count += countArrangements(design)
	}
	return count
}
