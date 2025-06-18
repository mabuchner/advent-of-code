package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("run failed: %s", err.Error())
		os.Exit(1)
	}
}

type BlockSpan struct {
	isUsed bool
	id     int
	length int
}

func run() error {
	file, err := os.Open("./assets/input.txt")
	// file, err := os.Open("./assets/input_small.txt")
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	disk := []BlockSpan{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for i := range scanner.Text() {
			ch := scanner.Text()[i]
			repeat := int(ch - '0')
			used := i%2 == 0
			disk = append(disk, BlockSpan{
				isUsed: used,
				id:     i / 2,
				length: repeat,
			})
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// fmt.Printf("%v\n", disk)

	for right := len(disk) - 1; right >= 0; right -= 1 {
		if !disk[right].isUsed {
			continue
		}
		requiredLength := disk[right].length

		for left := 0; left < right; left += 1 {
			if disk[left].isUsed || disk[left].length < requiredLength {
				continue
			}
			emptyLength := disk[left].length

			tmp := disk[right]

			disk[right].isUsed = false

			disk[left].length = emptyLength - requiredLength
			disk = slices.Insert(disk, left, tmp)

			break
		}
	}

	// fmt.Printf("%v\n", disk)

	diskStr := ""
	for _, b := range disk {
		data := "."
		if b.isUsed {
			data = strconv.FormatInt(int64(b.id), 10)
		}
		for i := 0; i < b.length; i += 1 {
			diskStr += data
		}
	}

	// fmt.Printf("%s\n", diskStr)

	index := 0
	checksum := 0
	for _, b := range disk {
		if !b.isUsed {
			index += b.length
			continue
		}

		for i := 0; i < b.length; i += 1 {
			checksum += index * b.id
			index += 1
		}
	}

	fmt.Printf("%d", checksum)

	return nil
}
