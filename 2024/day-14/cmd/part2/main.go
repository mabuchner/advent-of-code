package main

import (
	"bufio"
	"fmt"
	"math"
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

	// m := getMap(robots, w, h)
	// printMap(m, w, h)

	minVar := math.MaxFloat64
	minVarStep := -1
	for i := 0; i < 10000; i += 1 {
		robots = simulate(robots, w, h)

		v := calcVariance(robots, w, h)
		if v < minVar {
			minVar = v
			minVarStep = i + 1
		}
		// fmt.Printf("%f,", v)

		// m := getMap(robots, w, h)
		// printMap(m, w, h)
		// fmt.Printf("%d,%f\n", i+1, v)
		// fmt.Scanln()
	}
	// fmt.Printf("%+v\n", robots)

	fmt.Printf("%d", minVarStep)

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

func getMap(robots []robot, w, h int64) [][]rune {
	m := make([][]rune, h)
	for y := int64(0); y < h; y += 1 {
		r := make([]rune, w)
		for x := int64(0); x < w; x += 1 {
			r[x] = '.'
		}
		m[y] = r
	}

	for _, r := range robots {
		m[r.pos[1]][r.pos[0]] = 'x'
	}

	return m
}

func printMap(m [][]rune, w, h int64) {
	for y := int64(0); y < h; y += 1 {
		fmt.Printf("%s\n", string(m[y]))
	}
}

func calcVariance(robots []robot, w, h int64) float64 {
	sumPos := [2]float64{}
	for _, r := range robots {
		sumPos[0] += float64(r.pos[0])
		sumPos[1] += float64(r.pos[1])
	}

	l := float64(len(robots))
	avgPos := [2]float64{
		sumPos[0] / l,
		sumPos[1] / l,
	}

	sumDist := 0.0
	for _, r := range robots {
		diff := [2]float64{
			float64(r.pos[0]) - avgPos[0],
			float64(r.pos[1]) - avgPos[1],
		}

		distSqr := diff[0]*diff[0] + diff[1]*diff[1]
		sumDist += distSqr
	}

	return sumDist / l
}
