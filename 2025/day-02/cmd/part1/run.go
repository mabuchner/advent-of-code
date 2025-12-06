package main

import (
	"bufio"
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
	return process(input), nil
}

func load(inputPath string) ([][2]int64, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	ranges := make([][2]int64, 0, 5000)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue // Skip empty lines
		}

		for startWithEnd := range strings.SplitSeq(line, ",") {
			if len(startWithEnd) == 0 {
				continue // Skip empty range string (e.g. comma at the end of line)
			}
			s := strings.SplitN(startWithEnd, "-", 2)
			if len(s) != 2 {
				return nil, fmt.Errorf("malformed range '%s'", startWithEnd)
			}

			start, err := strconv.ParseInt(s[0], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse start '%s': %w", s[0], err)
			}

			end, err := strconv.ParseInt(s[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse end '%s': %w", s[1], err)
			}

			ranges = append(ranges, [2]int64{start, end})
		}
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("scan error: %w", scanner.Err())
	}

	return ranges, nil
}

func process(ranges [][2]int64) int64 {
	// Find invalid IDs (made only of some sequence of digits repeated twice)
	// and add them up

	invalidIDs := make([]int64, 0, 1024)
	for _, r := range ranges {
		for id := r[0]; id <= r[1]; id += 1 {
			if isInvalidID(id) {
				invalidIDs = append(invalidIDs, id)
			}
		}
	}

	sum := int64(0)
	for _, id := range invalidIDs {
		sum += id
	}
	return sum
}

func isInvalidID(id int64) bool {
	l := decLen(id)
	if l <= 1 {
		return false
	}

	fac := pow10[l/2]
	prefix := id / fac
	suffix := id % fac
	return prefix == suffix
}

const maxLen = 30

var pow10 [maxLen + 1]int64 = func() [maxLen + 1]int64 {
	res := [maxLen + 1]int64{}
	for i := range int64(maxLen + 1) {
		res[i] = powInt64(10, i)
	}
	return res
}()

func powInt64(n int64, exp int64) int64 {
	if exp == 0 {
		return int64(1)
	}

	res := n
	for exp > 1 {
		res *= n
		exp -= 1
	}
	return res
}

func decLen(n int64) int64 {
	l := int64(1)
	for n >= 10 {
		l += 1
		n /= 10
	}
	return l
}
