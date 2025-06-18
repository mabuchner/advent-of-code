package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

	leftNums := []int{}
	rightNums := []int{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := 0
		r := 0
		_, err := fmt.Sscan(scanner.Text(), &l, &r)
		if err != nil {
			return err
		}

		leftNums = append(leftNums, l)
		rightNums = append(rightNums, r)
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	sort.Slice(leftNums, func(i, j int) bool {
		return leftNums[i] - leftNums[j] > 0
	})
	sort.Slice(rightNums, func(i, j int) bool {
		return rightNums[i] - rightNums[j] > 0
	})

	if len(leftNums) != len(rightNums) {
		panic("mismatching lengths")
	}
	diffs := make([]int, 0, len(leftNums))
	for i := range leftNums {
		diff := leftNums[i] - rightNums[i]
		absDiff := max(diff, -diff)
		diffs = append(diffs, absDiff)
	}

	sum := 0
	for _, d := range diffs {
		sum += d
	}

	fmt.Printf("%d", sum)

	return nil
}
