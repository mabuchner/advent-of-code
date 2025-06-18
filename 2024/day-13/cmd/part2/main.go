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
			goal: [2]int64{px + 10000000000000, py + 10000000000000},
		})
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// fmt.Printf("machines=%+v\n", machines)

	tokenCount := int64(0)

	for _, m := range machines {
		a, b := solve(m)

		valid := isSolution(m, a, b)
		// fmt.Printf("a=%d, b=%d, solution=%t\n", a, b, valid)

		if valid {
			cost := 3*a + b
			tokenCount += cost
		}
	}

	fmt.Printf("%d", tokenCount)

	return nil
}

func solve(m machine) (a, b int64) {
	equationX := [3]int64{
		m.dirA[0],
		m.dirB[0],
		m.goal[0],
	}

	equationY := [3]int64{
		m.dirA[1],
		m.dirB[1],
		m.goal[1],
	}

	tmpX := [3]int64{
		equationY[1] * equationX[0],
		equationY[1] * equationX[1],
		equationY[1] * equationX[2],
	}

	tmpY := [3]int64{
		equationX[1] * equationY[0],
		equationX[1] * equationY[1],
		equationX[1] * equationY[2],
	}

	tmp := [3]int64{
		tmpX[0] - tmpY[0],
		tmpX[1] - tmpY[1],
		tmpX[2] - tmpY[2],
	}

	a = tmp[2] / tmp[0]

	b = (equationX[2] - (equationX[0] * a)) / equationX[1]

	return a, b
}

func isSolution(m machine, a, b int64) bool {
	x := m.dirA[0]*a + m.dirB[0]*b
	y := m.dirA[1]*a + m.dirB[1]*b
	return x == m.goal[0] && y == m.goal[1]
}
