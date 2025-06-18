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

	nums := []int64{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strs := strings.Split(scanner.Text(), " ")

		for _, s := range strs {
			num, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return err
			}
			nums = append(nums, num)
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// fmt.Printf("%v\n", nums)
	for i := 0; i < 25; i += 1 {
		nums = process(nums)
		// fmt.Printf("%v\n", nums)
	}

	fmt.Printf("%d", len(nums))

	return nil
}

func process(nums []int64) []int64 {
	res := []int64{}

	for _, num := range nums {
		if num == 0 {
			res = append(res, 1)
			continue
		}

		numStr := strconv.FormatInt(num, 10)
		if len(numStr)%2 == 0 {
			length := len(numStr)
			leftStr := numStr[:length/2]
			rightStr := numStr[length/2:]

			leftNum, _ := strconv.ParseInt(leftStr, 10, 64)
			rightNum, _ := strconv.ParseInt(rightStr, 10, 64)

			res = append(res, leftNum, rightNum)

			continue
		}

		res = append(res, num*2024)
	}

	return res
}
