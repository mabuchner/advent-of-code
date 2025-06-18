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
	// file, err := os.Open("./assets/input_small.txt") // 6
	file, err := os.Open("./assets/input.txt") // 353
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

	count := countPossibleDesigns(towels, designs)
	fmt.Printf("%d", count)

	return nil
}

func countPossibleDesigns(towels []string, designs []string) int {
	memo := map[string]bool{}
	memo[""] = true
	for _, towel := range towels {
		memo[towel] = true
	}

	var isDesignPossible func(design string) bool

	isDesignPossible = func(design string) bool {
		if possible, ok := memo[design]; ok {
			return possible
		}

		for i := 1; i <= len(design); i += 1 {
			prefix := design[0:i]
			suffix := design[i:]
			if memo[prefix] && isDesignPossible(suffix) {
				memo[design] = true
				return true
			}
			memo[design] = false
		}

		return false
	}

	count := 0
	for _, design := range designs {
		if isDesignPossible(design) {
			count += 1
		}
	}
	return count
}
