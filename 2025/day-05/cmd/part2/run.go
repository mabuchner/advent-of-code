package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func run(inputPath string) (int64, error) {
	ranges, _, err := load(inputPath)
	if err != nil {
		return 0, err
	}
	return process(ranges), nil
}

func load(inputPath string) ([][2]int64, []int64, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	readRanges := true
	ranges := [][2]int64{}
	ids := []int64{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			readRanges = false
			continue
		}

		if readRanges {
			startAndEnd := strings.SplitN(line, "-", 2)
			if len(startAndEnd) != 2 {
				return nil, nil, fmt.Errorf("malformed range '%s'", line)
			}

			start, errStart := strconv.ParseInt(startAndEnd[0], 10, 64)
			if errStart != nil {
				return nil, nil, fmt.Errorf("parse start '%s': %w", startAndEnd[0], err)
			}
			end, errEnd := strconv.ParseInt(startAndEnd[1], 10, 64)
			if errEnd != nil {
				return nil, nil, fmt.Errorf("parse end '%s': %w", startAndEnd[1], err)
			}

			ranges = append(ranges, [2]int64{start, end})
		} else {
			id, err := strconv.ParseInt(line, 10, 64)
			if err != nil {
				return nil, nil, fmt.Errorf("parse ID '%s': %w", line, err)
			}
			ids = append(ids, id)
		}
	}
	if scanner.Err() != nil {
		return nil, nil, fmt.Errorf("scan error: %w", scanner.Err())
	}

	return ranges, ids, nil
}

func process(ranges [][2]int64) int64 {
	merged := mergeRanges(ranges)
	freshCount := int64(0)
	for _, r := range merged {
		freshCount += r[1] - r[0] + 1
	}
	return freshCount
}

func mergeRanges(ranges [][2]int64) [][2]int64 {
	if len(ranges) == 0 {
		return nil
	}

	sorted := make([][2]int64, len(ranges))
	copy(sorted, ranges)
	sort.Slice(sorted, func(i int, j int) bool {
		return sorted[i][0] < sorted[j][0]
	})

	merged := make([][2]int64, 0, len(sorted))
	merged = append(merged, sorted[0])
	for i := 1; i < len(sorted); i += 1 {
		last := &merged[len(merged)-1]
		current := sorted[i]

		if current[0] <= last[1] {
			last[1] = max(last[1], current[1])
		} else {
			merged = append(merged, current)
		}
	}

	return merged
}
