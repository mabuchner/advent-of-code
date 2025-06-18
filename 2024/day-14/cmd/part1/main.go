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

type robot struct {
	pos [2]int64
	vel [2]int64
}

func run() error {
	file, err := os.Open("./assets/input.txt")
	w := int64(101)
	h := int64(103)

	// file, err := os.Open("./assets/input_small.txt")
	// w := int64(11)
	// h := int64(7)
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	robots := []robot{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), " ", 2)

		posStr := parts[0][2:] // remove "p="
		posPart := strings.SplitN(posStr, ",", 2)
		px, err := strconv.ParseInt(posPart[0], 10, 64)
		if err != nil {
			return err
		}
		py, err := strconv.ParseInt(posPart[1], 10, 64)
		if err != nil {
			return err
		}

		velStr := parts[1][2:] // remove "v="
		velParts := strings.SplitN(velStr, ",", 2)
		vx, err := strconv.ParseInt(velParts[0], 10, 64)
		if err != nil {
			return err
		}
		vy, err := strconv.ParseInt(velParts[1], 10, 64)
		if err != nil {
			return err
		}

		robots = append(robots, robot{
			pos: [2]int64{px, py},
			vel: [2]int64{vx, vy},
		})
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	// fmt.Printf("%+v\n", robots)

	for i := 0; i < 100; i += 1 {
		robots = simulate(robots, w, h)
	}
	// fmt.Printf("%+v\n", robots)

	safetyFactor := calcSafetyFactor(robots, w, h)
	fmt.Printf("%d", safetyFactor)

	return nil
}

func simulate(robots []robot, w, h int64) []robot {
	result := make([]robot, 0, len(robots))
	for _, r := range robots {
		pos := [2]int64{
			((r.pos[0] + r.vel[0]) + w) % w,
			((r.pos[1] + r.vel[1]) + h) % h,
		}
		result = append(result, robot{
			pos: pos,
			vel: r.vel,
		})
	}
	return result
}

func calcSafetyFactor(robots []robot, w, h int64) int64 {
	cx := w / 2
	cy := h / 2

	counts := [2][2]int64{
		[2]int64{0, 0},
		[2]int64{0, 0},
	}

	for _, r := range robots {
		x := min(1, max(-1, r.pos[0]-cx))
		y := min(1, max(-1, r.pos[1]-cy))
		counts[(y+1)/2][(x+1)/2] += (x * x) * (y * y)
	}

	return counts[0][0] * counts[0][1] * counts[1][0] * counts[1][1]
}
