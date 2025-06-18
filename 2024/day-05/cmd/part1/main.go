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

	orderings := [][2]int{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			break
		}

		before := 0
		after := 0
		_, err := fmt.Sscanf(scanner.Text(), "%d|%d", &before, &after)
		if err != nil {
			return err
		}

		orderings = append(orderings, [2]int{before, after})
	}

	updates := [][]int{}
	for scanner.Scan() {
		strs := strings.Split(scanner.Text(), ",")
		nums := []int{}
		for _, s := range strs {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return err
			}
			nums = append(nums, int(n))
		}
		updates = append(updates, nums)
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	// fmt.Printf("%v\n", orderings)
	// fmt.Printf("%v\n", updates)

	sum := 0
	validUpdates := findValidUpdates(updates, orderings)
	for _, update := range validUpdates {
		middleIndex := len(update) / 2
		sum += update[middleIndex]
	}
	fmt.Printf("%d", sum)

	return nil
}

func findValidUpdates(updates [][]int, orderings [][2]int) [][]int {
	validUpdates := [][]int{}

	befores := make(map[int][]int, len(orderings))
	for _, ordering := range orderings {
		befores[ordering[1]] = append(befores[ordering[1]], ordering[0])
	}
	// fmt.Printf("%v\n", befores)

	for _, update := range updates {
		updateSet := make(map[int]struct{}, len(update))
		for _, n := range update {
			updateSet[n] = struct{}{}
		}
		if isValidUpdate(update, updateSet, befores) {
			validUpdates = append(validUpdates, update)
		}
	}

	return validUpdates
}

func isValidUpdate(
	update []int,
	updateSet map[int]struct{},
	befores map[int][]int,
) bool {
	updated := map[int]struct{}{}
	for _, num := range update {
		mustBefores := befores[num]
		if !isAllUpdatedBefore(updateSet, updated, mustBefores) {
			return false
		}
		updated[num] = struct{}{}
	}
	return true
}

func isAllUpdatedBefore(updateSet map[int]struct{}, updated map[int]struct{}, mustBefores []int) bool {
	for _, mustBefore := range mustBefores {
		if _, ok := updateSet[mustBefore]; !ok {
			continue
		}

		if _, ok := updated[mustBefore]; !ok {
			return false
		}
	}
	return true
}
