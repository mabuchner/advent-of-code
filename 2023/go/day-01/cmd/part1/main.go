package main

import (
	"bufio"
	"errors"
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
	file, err := os.Open("./assets/input.txt")
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	sum, err := process(lines)
	if err != nil {
		return fmt.Errorf("process: %w", err)
	}

	fmt.Printf("%d", sum)

	return nil
}

func process(lines []string) (int, error) {
	// Precondition: strings are ASCII encoded
	sum := 0
	for _, line := range lines {
		leftNumIndex := strings.IndexFunc(line, isNumber)
		rightNumIndex := strings.LastIndexFunc(line, isNumber)
		if leftNumIndex < 0 || rightNumIndex < 0 {
			return -1, errors.New("no number found")
		}

		leftNum := int(line[leftNumIndex] - '0')
		rightNum := int(line[rightNumIndex] - '0')
		sum += leftNum * 10 + rightNum
	}
	return sum, nil
}

func isNumber(c rune) bool {
	return c >= '0' && c <= '9'
}
