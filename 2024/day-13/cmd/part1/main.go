package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("run failed: %s\n", err.Error())
		os.Exit(1)
	}
}

type machine struct {
	dirA [2]int64
	dirB [2]int64
	goal [2]int64
}

func run() error {
	file, err := os.Open("./assets/input.txt")
	// file, err := os.Open("./assets/input_small.txt")
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	buttonCoordRegex, err := regexp.Compile(`Button [AB]: X\+(\d+), Y\+(\d+)`)
	if err != nil {
		return err
	}

	priceCoordRegex, err := regexp.Compile(`Prize: X=(\d+), Y=(\d+)`)
	if err != nil {
		return err
	}

	machines := []machine{}

	scanner := bufio.NewScanner(file)
	for {
		if !scanner.Scan() {
			break
		}

		if len(scanner.Text()) == 0 {
			continue
		}

		// TODO: Add error handling

		matchA := buttonCoordRegex.FindStringSubmatch(scanner.Text())
		ax, err := strconv.ParseInt(matchA[1], 10, 64)
		if err != nil {
			return err
		}
		ay, err := strconv.ParseInt(matchA[2], 10, 64)
		if err != nil {
			return err
		}

		if !scanner.Scan() {
			break
		}
		matchB := buttonCoordRegex.FindStringSubmatch(scanner.Text())
		bx, err := strconv.ParseInt(matchB[1], 10, 64)
		if err != nil {
			return err
		}
		by, err := strconv.ParseInt(matchB[2], 10, 64)
		if err != nil {
			return err
		}

		if !scanner.Scan() {
			break
		}
		priceMatch := priceCoordRegex.FindStringSubmatch(scanner.Text())
		px, err := strconv.ParseInt(priceMatch[1], 10, 64)
		if err != nil {
			return err
		}
		py, err := strconv.ParseInt(priceMatch[2], 10, 64)
		if err != nil {
			return err
		}

		machines = append(machines, machine{
			dirA: [2]int64{ax, ay},
			dirB: [2]int64{bx, by},
			goal: [2]int64{px, py},
		})
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// fmt.Printf("machines=%+v\n", machines)

	tokenCount := int64(0)

	for _, m := range machines {
	next_machine:
		for a := int64(0); a < 100; a += 1 {
			posA := [2]int64{
				a * m.dirA[0],
				a * m.dirA[1],
			}
			if posA[0] > m.goal[0] && posA[1] > m.goal[1] {
				break next_machine
			}

			for b := int64(0); b < 100; b += 1 {
				posB := [2]int64{
					b * m.dirB[0],
					b * m.dirB[1],
				}

				pos := [2]int64{
					posA[0] + posB[0],
					posA[1] + posB[1],
				}

				if pos[0] > m.goal[0] && pos[1] > m.goal[1] {
					break
				}

				if pos[0] == m.goal[0] && pos[1] == m.goal[1] {
					cost := a*3 + b
					tokenCount += cost
					// fmt.Printf("WON!!! a=%d, b=%d, cost=%d\n", a, b, cost)
					break next_machine
				}
			}
		}
	}

	fmt.Printf("%d", tokenCount)

	return nil
}
