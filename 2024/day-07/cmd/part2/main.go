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
	// file, err := os.Open("./assets/input_small.txt")
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	sum := int64(0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strs := strings.SplitN(scanner.Text(), ": ", 2)

		result, err := strconv.ParseInt(strs[0], 10, 64)
		if err != nil {
			return err
		}

		numStrs := strings.Split(strs[1], " ")
		nums := []int64{}
		for _, s := range numStrs {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return err
			}
			nums = append(nums, n)
		}

		if isPossible(result, nums) {
			sum += result
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	fmt.Printf("%d", sum)

	return nil
}

func isPossible(result int64, nums []int64) bool {
	var recusion func(acc int64, i int) bool

	recusion = func(acc int64, i int) bool {
		if i >= len(nums) {
			return acc == result
		}

		if recusion(acc+nums[i], i+1) || recusion(acc*nums[i], i+1) {
			return true
		}

		accStr := strconv.FormatInt(acc, 10)
		numStr := strconv.FormatInt(nums[i], 10)
		newAcc, _ := strconv.ParseInt(accStr+numStr, 10, 64)
		return recusion(newAcc, i+1)
	}

	return recusion(0, 0)
}
