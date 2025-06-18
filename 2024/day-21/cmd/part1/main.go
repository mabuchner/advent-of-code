package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var (
	numpad = [][]byte{
		{'7', '8', '9'},
		{'4', '5', '7'},
		{'1', '2', '3'},
		{'G', '0', 'A'},
	}
	numpadButtonPositions = map[byte][2]int{
		'7': {0, 0},
		'8': {1, 0},
		'9': {2, 0},

		'4': {0, 1},
		'5': {1, 1},
		'6': {2, 1},

		'1': {0, 2},
		'2': {1, 2},
		'3': {2, 2},

		'G': {0, 3},
		'0': {1, 3},
		'A': {2, 3},
	}
	numpadInputs = computeBestPadInputs(numpad)

	dirpad = [][]byte{
		{'G', '^', 'A'},
		{'<', 'v', '>'},
	}
	dirpadButtonPositions = map[byte][2]int{
		'G': {0, 0},
		'^': {1, 0},
		'A': {2, 0},

		'<': {0, 1},
		'v': {1, 1},
		'>': {2, 1},
	}
	dirpadInputs = computeBestPadInputs(dirpad)
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("run failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	// file, err := os.Open("./assets/input_small.txt") // 126384
	file, err := os.Open("./assets/input.txt") // 205160
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	codes := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		codes = append(codes, scanner.Text())
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	// fmt.Printf("codes=%v\n", codes)

	// fmt.Printf("numpadInputs=%v\n", numpadInputs)
	// fmt.Printf("dirpadInputs=%v\n", dirpadInputs)

	complexity := 0
	for _, code := range codes {
		num, err := getCodeNum(code)
		if err != nil {
			return err
		}

		humanInputLen := getHumanInputLen(code, 2)

		complexity += humanInputLen * int(num)
	}

	fmt.Printf("%d", complexity)

	return nil
}

func computeBestPadInputs(pad [][]byte) map[[2][2]int]string {
	size := [2]int{len(pad[0]), len(pad)}

	inputs := map[[2][2]int]string{}
	for sy := 0; sy < size[1]; sy += 1 {
		for sx := 0; sx < size[0]; sx += 1 {
			start := [2]int{sx, sy}
			for ey := 0; ey < size[1]; ey += 1 {
				for ex := 0; ex < size[0]; ex += 1 {
					end := [2]int{ex, ey}
					inputs[[2][2]int{start, end}] = computeBestPadInput(start, end, pad)
				}
			}
		}
	}

	return inputs
}

func computeBestPadInput(start, end [2]int, pad [][]byte) string {
	// Ignore starting or ending on a gap
	if pad[start[1]][start[0]] == 'G' ||
		pad[end[1]][end[0]] == 'G' {
		return ""
	}

	leri := func(start [2]int) ([2]int, string, bool) {
		pos := start
		input := ""
		gap := false

		// Left
		for pos[0] > end[0] {
			input += "<"
			pos[0] -= 1
			gap = gap || pad[pos[1]][pos[0]] == 'G'
		}

		// Right
		for pos[0] < end[0] {
			input += ">"
			pos[0] += 1
			gap = gap || pad[pos[1]][pos[0]] == 'G'
		}

		return pos, input, gap
	}

	updo := func(start [2]int) ([2]int, string, bool) {
		pos := start
		input := ""
		gap := false

		// Up
		for pos[1] > end[1] {
			input += "^"
			pos[1] -= 1
			gap = gap || pad[pos[1]][pos[0]] == 'G'
		}

		// Down
		for pos[1] < end[1] {
			input += "v"
			pos[1] += 1
			gap = gap || pad[pos[1]][pos[0]] == 'G'
		}

		return pos, input, gap
	}

	// Going left?
	if start[0] > end[0] {
		// Go left first, then up or down
		midPos, input1, gap1 := leri(start)
		finalPos, input2, gap2 := updo(midPos)

		// If we hit a gap, then try the other order
		if gap1 || gap2 {
			midPos, input1, gap1 = updo(start)
			finalPos, input2, gap2 = leri(midPos)

			if gap1 || gap2 {
				fmt.Printf(
					"Hit gap both ways (start=%v, end=%v, pad=%v)\n",
					start,
					end,
					pad,
				)
			}
			if finalPos != end {
				fmt.Printf(
					"Missed end (start=%v, end=%v, finalPos=%v, pad=%v)\n",
					start,
					end,
					finalPos,
					pad,
				)
			}
		}

		return input1 + input2 + "A"
	} else {
		// Go up or down, then left or right
		midPos, input1, gap1 := updo(start)
		finalPos, input2, gap2 := leri(midPos)

		// If we hit a gap, then try the other order
		if gap1 || gap2 {
			midPos, input1, gap1 = leri(start)
			finalPos, input2, gap2 = updo(midPos)

			if gap1 || gap2 {
				fmt.Printf(
					"Hit gap both ways (start=%v, end=%v, pad=%v)\n",
					start,
					end,
					pad,
				)
			}
			if finalPos != end {
				fmt.Printf(
					"Missed end (start=%v, end=%v, finalPos=%v, pad=%v)\n",
					start,
					end,
					finalPos,
					pad,
				)
			}
		}

		return input1 + input2 + "A"
	}
}

func getHumanInputLen(code string, robotDirpadCount int) int {
	// fmt.Printf("[c] code=%s\n", code)
	input := getNextInput(numpad, numpadButtonPositions, numpadInputs, code)
	// fmt.Printf("[%d] input=%s\n", robotDirpadCount, input)
	for robotDirpadCount > 0 {
		input = getNextInput(dirpad, dirpadButtonPositions, dirpadInputs, input)
		robotDirpadCount -= 1
		// fmt.Printf("[%d] input=%s\n", robotDirpadCount, input)
	}
	return len(input)
}

func getNextInput(
	pad [][]byte,
	buttonPositions map[byte][2]int,
	allInputs map[[2][2]int]string,
	code string,
) string {
	pos := buttonPositions['A']
	input := ""
	for i := range code {
		c := code[i]
		newPos := buttonPositions[c]
		input += allInputs[[2][2]int{pos, newPos}]
		pos = newPos
	}
	return input
}

func getCodeNum(code string) (int64, error) {
	numStr := code[0:3]
	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to parse code number from '%s': %w",
			code,
			err,
		)
	}
	return num, err
}
