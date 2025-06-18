package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	safeRecords := [][]int64{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strs := strings.Split(scanner.Text(), " ")

		record := make([]int64, 0, len(strs))
		for _, s := range strs {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return err
			}
			record = append(record, n)
		}

		for skipIndex := -1; skipIndex < len(record); skipIndex += 1 {
			if isSafe(record, skipIndex) {
				safeRecords = append(safeRecords, record)
				break
			}
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	fmt.Printf("%d", len(safeRecords))

	return nil
}

func isSafe(record []int64, skipIndex int) bool {
	firstIndex := 0
	secondIndex := 1
	if skipIndex == 0 {
		firstIndex += 1
		secondIndex += 1
	} else if skipIndex == 1 {
		secondIndex += 1
	}

	if secondIndex >= len(record) {
		return true
	}

	safeDiffMin := int64(1)
	safeDiffMax := int64(3)
	if record[firstIndex] > record[secondIndex] {
		safeDiffMin = -safeDiffMin
		safeDiffMax = -safeDiffMax

		tmp := safeDiffMin
		safeDiffMin = safeDiffMax
		safeDiffMax = tmp
	}

	for i := 0; i < len(record); i += 1 {
		if i == skipIndex {
			continue
		}

		lastIndex := i - 1
		if lastIndex == skipIndex {
			lastIndex -= 1
		}
		if lastIndex < 0 {
			continue
		}

		last := record[lastIndex]
		current := record[i]
		diff := current - last
		if diff < safeDiffMin || diff > safeDiffMax {
			return false
		}
	}

	return true
}
