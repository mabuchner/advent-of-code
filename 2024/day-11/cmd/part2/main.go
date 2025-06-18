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

	memo := make(map[int64][]int64, len(nums))
	memoLen := make(map[int64]int64)

	length := int64(0)
	for i, num := range nums {
		fmt.Printf("%d: start (num=%d)\n", i, num)
		results := process25([]int64{num}, memo, memoLen)
		results = process25(results, memo, memoLen)
		length += process25len(results, memo, memoLen)
		fmt.Printf("%d: end (length=%d)\n", i, length)
	}
	fmt.Printf("%d", length)

	return nil
}

func process25(nums []int64, memo map[int64][]int64, memoLen map[int64]int64) []int64 {
	fmt.Printf("process25\n")
	results := []int64{}
	for i, num := range nums {
		if i%100000 == 0 {
			fmt.Printf("%d/%d\n", i, len(nums))
		}
		if mem, ok := memo[num]; ok {
			results = append(results, mem...)
			continue
		}

		result := []int64{num}
		for i := 0; i < 25; i += 1 {
			result = process(result)
		}
		memo[num] = result
		memoLen[num] = int64(len(result))
		results = append(results, result...)
	}
	return results
}

func process25len(nums []int64, memo map[int64][]int64, memoLen map[int64]int64) int64 {
	fmt.Printf("process25len\n")
	length := int64(0)
	for i, num := range nums {
		if i%100000 == 0 {
			fmt.Printf("%d/%d\n", i, len(nums))
		}

		if memLen, ok := memoLen[num]; ok {
			length += memLen
			continue
		}

		if mem, ok := memo[num]; ok {
			l := int64(len(mem))
			memoLen[num] = l
			length += l
			continue
		}

		result := []int64{num}
		for i := 0; i < 25; i += 1 {
			result = process(result)
		}
		length += int64(len(result))
	}
	return length
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
