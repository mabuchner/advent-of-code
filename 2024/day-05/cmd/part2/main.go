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

	sum := 0
	invalidUpdates := findInvalidUpdates(updates, orderings)
	// fmt.Printf("invalidUpdates=%v\n", invalidUpdates)
	fixedUpdates := fixInvalidUpdates(invalidUpdates, orderings)
	// fmt.Printf("fixedUpdates=%v\n", fixedUpdates)
	for _, update := range fixedUpdates {
		middleIndex := len(update) / 2
		sum += update[middleIndex]
	}
	fmt.Printf("%d", sum)

	return nil
}

func findInvalidUpdates(updates [][]int, orderings [][2]int) [][]int {
	invalidUpdates := [][]int{}

	befores := buildBeforeMap(orderings)

	for _, update := range updates {
		updateSet := buildUpdateSet(update)
		if !isValidUpdate(update, updateSet, befores) {
			invalidUpdates = append(invalidUpdates, update)
		}
	}

	return invalidUpdates
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

func fixInvalidUpdates(updates [][]int, orderings [][2]int) [][]int {
	fixedUpdates := [][]int{}
	for _, update := range updates {
		fixedUpdates = append(fixedUpdates, fixInvalidUpdate(update, orderings))
	}
	return fixedUpdates
}

func fixInvalidUpdate(update []int, orderings [][2]int) []int {
	// fmt.Printf("fixing: %v\n", update)

	fixedUpdate := make([]int, len(update))
	copy(fixedUpdate, update)

	indices := make(map[int]int, len(update))
	for i, n := range update {
		indices[n] = i
	}

	befores := buildBeforeMap(orderings)
	// for _, n := range fixedUpdate {
	// 	mustBefores := befores[n]
	// 	fmt.Printf("%d->", n)
	// 	for _, b := range mustBefores {
	// 		if _, ok := indices[b]; ok {
	// 			fmt.Printf(" %d", b)
	// 		}
	// 	}
	// 	fmt.Printf("\n")
	// }

	for i := 0; i < len(fixedUpdate); i += 1 {
		num := fixedUpdate[i]
		mustBefores := befores[num]
		maxBeforeIndex := findMaxIndex(indices, mustBefores)
		if i < maxBeforeIndex {
			beforeNum := fixedUpdate[maxBeforeIndex]
			fixedUpdate[i], fixedUpdate[maxBeforeIndex] = fixedUpdate[maxBeforeIndex], fixedUpdate[i]
			indices[beforeNum] = i
			indices[num] = maxBeforeIndex
			i -= 1 // Keep processing the current index with the name value
		}
	}

	return fixedUpdate
}

func findMaxIndex(indices map[int]int, nums []int) int {
	maxIndex := -1
	for _, n := range nums {
		if index, ok := indices[n]; ok && index > maxIndex {
			maxIndex = index
		}
	}
	return maxIndex
}

func buildBeforeMap(orderings [][2]int) map[int][]int {
	m := make(map[int][]int, len(orderings))
	for _, ordering := range orderings {
		m[ordering[1]] = append(m[ordering[1]], ordering[0])
	}
	return m
}

func buildUpdateSet(update []int) map[int]struct{} {
	set := make(map[int]struct{}, len(update))
	for _, n := range update {
		set[n] = struct{}{}
	}
	return set
}
