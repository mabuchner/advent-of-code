package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

var generateSVGEnabled bool = false

func run(inputPath string) (int64, error) {
	input, err := load(inputPath)
	if err != nil {
		return 0, err
	}
	return process(input), nil
}

func load(inputPath string) ([][2]int64, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	tiles := [][2]int64{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		s := strings.SplitN(line, ",", 2)
		if len(s) != 2 {
			return nil, errors.New("unexpected input")
		}

		x, err := strconv.ParseInt(s[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse x: %w", err)
		}

		y, err := strconv.ParseInt(s[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse y: %w", err)
		}

		tiles = append(tiles, [2]int64{x, y})
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("scan error: %w", scanner.Err())
	}

	return tiles, nil
}

type candidate struct {
	i    int
	j    int
	area int64
}

func process(tiles [][2]int64) int64 {
	minX := int64(math.MaxInt64)
	maxX := int64(0)
	minY := int64(math.MaxInt64)
	maxY := int64(0)
	for _, tile := range tiles {
		minX = min(minX, tile[0])
		maxX = max(maxX, tile[0])
		minY = min(minY, tile[1])
		maxY = max(maxY, tile[1])
	}
	// fmt.Printf("minX=%d, minY=%d, maxX=%d, maxY=%d\n", minX, minY, maxX, maxY)

	// Find all candidates sorted by area
	candidates := []candidate{}
	for i := 0; i < len(tiles); i += 1 {
		ti := tiles[i]
		for j := i + 1; j < len(tiles); j += 1 {
			tj := tiles[j]
			a := calcArea(ti, tj)
			candidates = append(candidates, candidate{
				i:    i,
				j:    j,
				area: a,
			})
		}
	}
	slices.SortFunc(candidates, func(a, b candidate) int {
		return int(min(1, max(-1, b.area-a.area)))
	})

	// Process candidates until we found a valid one
	var bestCandidate *candidate
	for i, candidate := range candidates {
		fmt.Printf("Processing candidate %d ...\n", i)
		if isValid(tiles, candidate) {
			bestCandidate = &candidate
			break
		}
	}
	fmt.Printf("Done!\n")

	if generateSVGEnabled {
		svg := generateSVG(minX, minY, maxX, maxY, tiles, bestCandidate)
		err := os.WriteFile("out.svg", []byte(svg), 0644)
		if err != nil {
			fmt.Println("Failed to write SVG file")
		}
	}

	if bestCandidate == nil {
		return -1
	}
	return bestCandidate.area
}

func calcArea(pa, pb [2]int64) int64 {
	minX := min(pa[0], pb[0])
	minY := min(pa[1], pb[1])
	maxX := max(pa[0], pb[0])
	maxY := max(pa[1], pb[1])
	return (maxX - minX + 1) * (maxY - minY + 1)
}

func isValid(tiles [][2]int64, c candidate) bool {
	pa := tiles[c.i]
	pb := tiles[c.j]

	minX := min(pa[0], pb[0])
	minY := min(pa[1], pb[1])
	maxX := max(pa[0], pb[0])
	maxY := max(pa[1], pb[1])

	rectVertices := [4][2]int64{
		{minX, minY}, // Top-left
		{maxX, minY}, // Top-right
		{maxX, maxY}, // Bottom-right
		{minX, maxY}, // Bottom-left
	}

	// Check all corners before we do a more exhaustive check
	if !isPointInPolygon(tiles, rectVertices[0]) ||
		!isPointInPolygon(tiles, rectVertices[1]) ||
		!isPointInPolygon(tiles, rectVertices[2]) ||
		!isPointInPolygon(tiles, rectVertices[3]) {
		return false
	}

	return isEdgeInPolygon(tiles, rectVertices[0], rectVertices[1]) &&
		isEdgeInPolygon(tiles, rectVertices[1], rectVertices[2]) &&
		isEdgeInPolygon(tiles, rectVertices[2], rectVertices[3]) &&
		isEdgeInPolygon(tiles, rectVertices[3], rectVertices[0])
}

func isEdgeInPolygon(vertices [][2]int64, start, end [2]int64) bool {
	delta := [2]int64{
		min(1, max(-1, end[0]-start[0])),
		min(1, max(-1, end[1]-start[1])),
	}

	p := start
	for {
		if !isPointInPolygon(vertices, p) {
			return false
		}

		if p == end {
			break
		}

		p[0] += delta[0]
		p[1] += delta[1]
	}
	return true
}

func isPointInPolygon(vertices [][2]int64, p [2]int64) bool {
	n := len(vertices)
	crossings := 0
	for i := 0; i < n; i += 1 {
		v0 := vertices[i]
		v1 := vertices[(i+1)%n]

		// Horizontal edge?
		if v0[1] == v1[1] {
			// Special case: Point on horizontal edge?
			minX := min(v0[0], v1[0])
			maxX := max(v0[0], v1[0])
			if p[1] == v0[1] && p[0] >= minX && p[0] <= maxX {
				return true
			}

			// Otherwise, skip horizontal edge
			continue
		}

		// Special case: Point on vertical edge?
		minY := min(v0[1], v1[1])
		maxY := max(v0[1], v1[1])
		if v0[0] == v1[0] {
			if p[0] == v0[0] && p[1] >= minY && p[1] <= maxY {
				return true
			}
		}

		// Skip edges on the right of the point
		if v0[0] > p[0] && v1[0] > p[0] {
			continue
		}

		// Count the number of horizontal intersections
		// (< maxY to avoid counting the same edge point twice)
		if p[1] >= minY && p[1] < maxY {
			crossings += 1
		}
	}

	inPolygon := (crossings % 2) == 1
	return inPolygon
}

func generateSVG(
	minX int64,
	minY int64,
	maxX int64,
	maxY int64,
	points [][2]int64,
	cdd *candidate,
) string {
	var sb strings.Builder

	// Padding as percentage
	width := maxX - minX
	height := maxY - minY
	paddingPercent := 0.02
	maxDim := max(height, width)
	padding := max(1, int64(float64(maxDim)*paddingPercent))

	// Specify viewport
	viewBoxMinX := minX - padding
	viewBoxMinY := minY - padding
	viewBoxWidth := width + 2*padding
	viewBoxHeight := height + 2*padding

	// SVG start
	sb.WriteString(
		fmt.Sprintf(
			`<svg xmlns="http://www.w3.org/2000/svg" with="800" height="600" preserveAspectRatio="xMidYMid meet" viewBox="%d %d %d %d">`,
			viewBoxMinX,
			viewBoxMinY,
			viewBoxWidth,
			viewBoxHeight,
		),
	)

	// Polygons
	sb.WriteString(`<polygon fill="none" stroke="black" stroke-width="1" vector-effect="non-scaling-stroke" points="`)
	for i, p := range points {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d,%d", p[0], p[1]))
	}
	sb.WriteString(`"/>`)

	// Highlight area
	if cdd != nil {
		p0 := points[cdd.i]
		p1 := points[cdd.j]
		minX := min(p0[0], p1[0])
		maxX := max(p0[0], p1[0])
		minY := min(p0[1], p1[1])
		maxY := max(p0[1], p1[1])

		x := minX
		y := minY
		width := maxX - minX
		height := maxY - minY

		sb.WriteString(
			fmt.Sprintf(
				`<rect x="%d" y="%d" width="%d" height="%d" fill="red" opacity="0.3" stroke="none" stroke-width="1" vector-effect="non-scaling-stroke"/>`,
				x,
				y,
				width,
				height,
			),
		)
	}

	// SVG end
	sb.WriteString(`</svg>`)

	return sb.String()
}
