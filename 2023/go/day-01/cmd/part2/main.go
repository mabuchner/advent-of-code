package main

import (
	"bufio"
	"fmt"
	"math"
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

var numWords = [...]string{
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}

var numWordToNum = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

func process(lines []string) (int, error) {
	// Precondition: strings are ASCII encoded
	sum := 0
	for _, line := range lines {
		leftNum := 0
		rightNum := 0

		leftNumIndex := strings.IndexFunc(line, isNumber)
		rightNumIndex := strings.LastIndexFunc(line, isNumber)
		if leftNumIndex >= 0 {
			leftNum = int(line[leftNumIndex] - '0')
		}
		if rightNumIndex >= 0 {
			rightNum = int(line[rightNumIndex] - '0')
		}

		leftNumWordIndex, leftNumWord := indexAnySubString(line, numWords[:])
		rightNumWordIndex, rightNumWord := lastIndexAnySubString(line, numWords[:])
		if leftNumWordIndex >= 0 && (leftNumIndex == -1 || leftNumWordIndex < leftNumIndex) {
			leftNum = remapStr(leftNumWord, numWordToNum)
		}
		if rightNumWordIndex >= 0 && (rightNumIndex == -1 || rightNumWordIndex > rightNumIndex) {
			rightNum = remapStr(rightNumWord, numWordToNum)
		}

		num := leftNum * 10 + rightNum
		sum += num
	}
	return sum, nil
}

func isNumber(c rune) bool {
	return c >= '0' && c <= '9'
}

func indexAnySubString(s string, substrs []string) (int, string) {
	minIndex := math.MaxInt32
	minIndexSubstr := ""
	found := false
	for _, substr := range substrs {
		index := strings.Index(s, substr)
		if index >= 0 && index < minIndex {
			minIndex = index
			minIndexSubstr = substr
			found = true
		}
	}

	if !found {
		minIndex = -1
	}

	return minIndex, minIndexSubstr
}

func lastIndexAnySubString(s string, substrs []string) (int, string) {
	maxIndex := -1
	maxIndexSubstr := ""
	for _, substr := range substrs {
		index := strings.LastIndex(s, substr)
		if index >= 0 && index > maxIndex {
			maxIndex = index
			maxIndexSubstr = substr
		}
	}
	return maxIndex, maxIndexSubstr
}

func remapStr(s string, m map[string]int) int {
	return m[s]
}
