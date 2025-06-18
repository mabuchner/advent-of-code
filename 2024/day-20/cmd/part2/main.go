package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("run failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	// file, err := os.Open("./assets/input_small.txt") // savingCounts (with saving>=50)=map[50:32 52:31 54:29 56:39 58:25 60:23 62:20 64:19 66:12 68:14 70:12 72:22 74:4 76:3]
	file, err := os.Open("./assets/input.txt") // 999556
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	m := [][]byte{}
	startPos := [2]int{}
	endPos := [2]int{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := []byte{}
		for x := range scanner.Text() {
			ch := scanner.Text()[x]
			if ch == 'S' {
				startPos[0] = x
				startPos[1] = len(m)
				ch = '.'
			} else if ch == 'E' {
				endPos[0] = x
				endPos[1] = len(m)
				ch = '.'
			}
			row = append(row, ch)
		}
		m = append(m, row)
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// printMap(m)
	// fmt.Printf("startPos=%v, endPos=%v\n", startPos, endPos)

	path, dists := findPathDists(m, startPos, endPos)
	// printDistsMap(dists)
	// fmt.Printf("%v\n", path)

	savings := findCheatSavings(path, dists, 20)
	// fmt.Printf("savings=%v\n", savings)

	bigSavingCount := 0
	savingCounts := map[int]int{}
	for _, saving := range savings {
		if saving >= 50 {
			savingCounts[saving] += 1
		}
		if saving >= 100 {
			bigSavingCount += 1
		}
	}
	// fmt.Printf("savingCounts=%v\n", savingCounts)
	fmt.Printf("%v", bigSavingCount)

	return nil
}

func printMap(m [][]byte) {
	for _, row := range m {
		fmt.Printf("%s\n", string(row))
	}
}

func printDistsMap(m [][]int) {
	buf := strings.Builder{}
	for _, row := range m {
		for _, d := range row {
			buf.WriteString(fmt.Sprintf("%4d ", d))
		}
		fmt.Printf("%s\n", buf.String())
		buf.Reset()
	}
}

func findPathDists(m [][]byte, startPos [2]int, endPos [2]int) ([][2]int, [][]int) {
	dirs := [4][2]int{
		[2]int{1, 0},
		[2]int{0, 1},
		[2]int{-1, 0},
		[2]int{0, -1},
	}

	size := [2]int{len(m[0]), len(m)}

	dists := make([][]int, size[1])
	for y := range dists {
		dists[y] = make([]int, size[0])
		for x := range dists[y] {
			dists[y][x] = -1
		}
	}

	path := [][2]int{startPos}
	pos := startPos
	dists[startPos[1]][startPos[0]] = 0
	for pos != endPos {
		var newPos [2]int
		for _, dir := range dirs {
			newPos = [2]int{
				pos[0] + dir[0],
				pos[1] + dir[1],
			}

			if dists[newPos[1]][newPos[0]] == -1 &&
				mapAt(m, size, newPos) == '.' {
				break
			}
		}
		dists[newPos[1]][newPos[0]] = len(path)
		path = append(path, newPos)
		pos = newPos
	}
	return path, dists
}

func findCheatSavings(
	path [][2]int,
	dists [][]int,
	maxCheatDelta int,
) map[[2][2]int]int {
	savings := map[[2][2]int]int{}

	size := [2]int{len(dists[0]), len(dists)}

	for _, pos := range path {
		posDist := dists[pos[1]][pos[0]]
		for dy := -maxCheatDelta; dy <= maxCheatDelta; dy += 1 {
			for dx := -maxCheatDelta; dx <= maxCheatDelta; dx += 1 {
				cheatDelta := abs(dx) + abs(dy)
				if cheatDelta > 0 && cheatDelta <= maxCheatDelta {
					cheatPos := [2]int{
						pos[0] + dx,
						pos[1] + dy,
					}

					if cheatPos[0] < 0 || cheatPos[0] >= size[0] ||
						cheatPos[1] < 0 || cheatPos[1] >= size[1] {
						continue
					}

					// Not on a walkable tile?
					cheatPosDist := dists[cheatPos[1]][cheatPos[0]]
					if cheatPosDist == -1 {
						continue
					}

					saving := cheatPosDist - posDist - cheatDelta
					if saving > 0 {
						savings[[2][2]int{pos, cheatPos}] = saving
					}
				}
			}
		}
	}

	return savings
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func mapAt(m [][]byte, size [2]int, pos [2]int) byte {
	if pos[0] < 0 || pos[0] >= size[0] ||
		pos[1] < 0 || pos[1] >= size[1] {
		return 0
	}
	return m[pos[1]][pos[0]]
}
