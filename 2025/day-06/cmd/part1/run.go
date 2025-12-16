package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func run(inputPath string) (int64, error) {
	input, err := load(inputPath)
	if err != nil {
		return 0, err
	}

	res, err := process(input)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func load(inputPath string) ([][]string, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	rows := [][]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		cols := strings.Fields(line)
		if len(rows) > 0 && len(rows[0]) != len(cols) {
			return nil, errors.New("inconsistent column count")
		}

		rows = append(rows, cols)
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("scan error: %w", scanner.Err())
	}

	return rows, nil
}

func process(rows [][]string) (int64, error) {
	if len(rows) <= 1 {
		return 0, nil
	}

	rowCount := len(rows)
	colCount := len(rows[0])
	res := int64(0)
	for c := range colCount {
		op := rows[len(rows)-1][c]
		if op == "+" {
			sum := int64(0)
			for r := 0; r < rowCount-1; r += 1 {
				num, err := strconv.ParseInt(rows[r][c], 10, 64)
				if err != nil {
					return 0, fmt.Errorf("parse number at %d,%d: %w", r, c, err)
				}
				sum += num
			}
			res += sum
		} else if op == "*" {
			sum := int64(1)
			for r := 0; r < rowCount-1; r += 1 {
				num, err := strconv.ParseInt(rows[r][c], 10, 64)
				if err != nil {
					return 0, fmt.Errorf("parse number at %d,%d: %w", r, c, err)
				}
				sum *= num
			}
			res += sum
		} else {
			return 0, fmt.Errorf("unexpected operation '%s'", op)
		}
	}

	return res, nil
}
