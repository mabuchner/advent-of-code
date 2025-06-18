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
		fmt.Printf("run failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	file, err := os.Open("./assets/input.txt") // 3059
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	sum, err := process(lines)
	if err != nil {
		return fmt.Errorf("process: %w", err)
	}

	fmt.Printf("%d", sum)

	return nil
}

type draw struct {
	red   int64
	green int64
	blue  int64
}

type game struct {
	id    int64
	draws []draw
}

func process(lines []string) (int64, error) {
	games := []game{}
	for _, line := range lines {
		game, err := gameFromStr(line)
		if err != nil {
			return 0, fmt.Errorf("gameFromStr: %w", err)
		}
		games = append(games, game)

	}

	sum := int64(0)
	for _, g := range games {
		if isPossible(&g, 12, 13, 14) {
			sum += g.id
		}
	}
	return sum, nil
}

func gameFromStr(s string) (game, error) {
	// <GAME_AND_DRAWS> := <GAME_STR>: <DRAWS_STR>
	//
	// <GAME_STR> := Game <ID_STR>
	//
	// <DRAWS_STR> := <DRAW_STR_0>; <DRAW_STR_1>; ...; <DRAW_STR_N>
	// <DRAW_STR_N> := [<GREEN_COUNT> green, ][<BLUE_COUNT> blue, ]<RED_COUNT> red
	gameAndDraws := strings.SplitN(s, ": ", 2)

	gameStr := gameAndDraws[0]
	idStr := strings.SplitN(gameStr, " ", 2)[1]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return game{}, fmt.Errorf("ParseInt game ID: %w", err)
	}

	drawsStr := gameAndDraws[1]
	drawStrs := strings.Split(drawsStr, "; ")
	draws := make([]draw, 0, len(drawStrs))
	for _, drawStr := range drawStrs {
		draw := draw{}
		countColorStrs := strings.Split(drawStr, ", ")
		for _, countColorStr := range countColorStrs {
			countColor := strings.SplitN(countColorStr, " ", 2)

			count, err := strconv.ParseInt(countColor[0], 10, 64)
			if err != nil {
				return game{}, fmt.Errorf("ParseInt color count: %w", err)
			}

			switch countColor[1] {
			case "red":
				draw.red += count
			case "green":
				draw.green += count
			case "blue":
				draw.blue += count
			}

		}
		draws = append(draws, draw)
	}

	return game{id: id, draws: draws}, nil
}

func isPossible(g *game, maxRed, maxGreen, maxBlue int64) bool {
	for _, draw := range g.draws {
		if draw.red > maxRed || draw.green > maxGreen || draw.blue > maxBlue {
			return false
		}
	}
	return true
}
