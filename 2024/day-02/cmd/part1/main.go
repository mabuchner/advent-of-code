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

		if isSafe(record) {
			safeRecords = append(safeRecords, record)
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	fmt.Printf("%d", len(safeRecords))

	return nil
}

func isSafe(record []int64) bool {
	if len(record) <= 1 {
		return true
	}

	safeDiffMin := int64(1)
	safeDiffMax := int64(3)
	if record[0] > record[1] {
		safeDiffMin = -safeDiffMin
		safeDiffMax = -safeDiffMax

		tmp := safeDiffMin
		safeDiffMin = safeDiffMax
		safeDiffMax = tmp
	}

	for i := 1; i < len(record); i += 1 {
		last := record[i-1]
		current := record[i]
		diff := current - last
		if diff < safeDiffMin || diff > safeDiffMax {
			return false
		}
	}

	return true
}
