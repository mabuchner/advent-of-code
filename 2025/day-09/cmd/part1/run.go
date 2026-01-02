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
	return process(input), nil
}

func load(inputPath string) ([][2]int64, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	tiles := [][2]int64{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		s := strings.SplitN(line, ",", 2)
		if len(s) != 2 {
			return nil, errors.New("unexpected input")
		}

		x, err := strconv.ParseInt(s[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse x: %w", err)
		}

		y, err := strconv.ParseInt(s[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse y: %w", err)
		}

		tiles = append(tiles, [2]int64{x, y})
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("scan error: %w", scanner.Err())
	}

	return tiles, nil
}

func process(tiles [][2]int64) int64 {
	maxArea := int64(0)
	for i := 0; i < len(tiles); i += 1 {
		ti := tiles[i]
		for j := i + 1; j < len(tiles); j += 1 {
			tj := tiles[j]
			a := calcArea(ti, tj)
			maxArea = max(maxArea, a)
		}
	}
	return maxArea
}

func calcArea(pa, pb [2]int64) int64 {
	minX := min(pa[0], pb[0])
	minY := min(pa[1], pb[1])
	maxX := max(pa[0], pb[0])
	maxY := max(pa[1], pb[1])
	return (maxX - minX + 1) * (maxY - minY + 1)
}
