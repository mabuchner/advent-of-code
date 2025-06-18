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
	// file, err := os.Open("./assets/input_small.txt")
	file, err := os.Open("./assets/input.txt")
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

	sum := int64(0)
	for _, num := range initNums {
		sum += calcNextNumN(num, 2000)
	}
	fmt.Printf("%d", sum)

	return nil
}

func calcNextNumN(num int64, n int) int64 {
	for i := 0; i < n; i += 1 {
		num = calcNextNum(num)
	}
	return num
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
