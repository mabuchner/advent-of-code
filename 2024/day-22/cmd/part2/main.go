package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("run failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	// file, err := os.Open("./assets/input_small2.txt") // [-2 1 -1 3] =>  23
	file, err := os.Open("./assets/input.txt") // [-3 0 4 0] => 1570
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	initNums := []int64{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		num, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			return err
		}
		initNums = append(initNums, num)
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	// fmt.Printf("%v\n", initNums)

	allPrices := make([][]int, len(initNums))
	for i, num := range initNums {
		allPrices[i] = calcPricesN(num, 2000)
	}

	allDeltas := make([][]int, len(allPrices))
	for i, prices := range allPrices {
		allDeltas[i] = calcPriceDeltas(prices)
	}

	patternPrices := make(map[[4]int]int, 20000)
	for buyerIndex := range allDeltas {
		deltas := allDeltas[buyerIndex]
		prices := allPrices[buyerIndex]
		buyerPatternSet := make(map[[4]int]struct{}, 2000)
		for i := 3; i < len(deltas); i += 1 {
			pattern := [4]int{
				deltas[i-3],
				deltas[i-2],
				deltas[i-1],
				deltas[i],
			}

			// Sum the first occurrence of each pattern within each buyer
			if _, ok := buyerPatternSet[pattern]; !ok {
				patternPrices[pattern] += prices[i]
				buyerPatternSet[pattern] = struct{}{}
			}
		}
	}

	maxSum := 0
	// maxPattern := [4]int{}
	for _, sum := range patternPrices {
		if sum > maxSum {
			maxSum = sum
			// maxPattern = pattern
		}
	}
	// fmt.Printf("maxPattern=%v\n", maxPattern)
	fmt.Printf("%d", maxSum)

	return nil
}

func calcPriceDeltas(prices []int) []int {
	deltas := make([]int, len(prices))
	for i := 1; i < len(prices); i += 1 {
		deltas[i] = prices[i] - prices[i-1]
	}
	return deltas
}

func calcPricesN(num int64, n int) []int {
	prices := make([]int, n)
	for i := 0; i < n; i += 1 {
		prices[i] = calcPrice(num)
		num = calcNextNum(num)
	}
	return prices
}

func calcPrice(num int64) int {
	return int(num % 10)
}

func calcNextNum(num int64) int64 {
	num = prune((num * 64) ^ num)
	num = prune((num / 32) ^ num)
	num = prune((num * 2048) ^ num)
	return num
}

func prune(num int64) int64 {
	return num % 16777216
}
