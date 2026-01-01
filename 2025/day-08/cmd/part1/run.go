package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func run(inputPath string) (int64, error) {
	input, err := load(inputPath)
	if err != nil {
		return 0, err
	}
	return process(input), nil
}

func load(inputPath string) ([][3]float64, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	boxes := [][3]float64{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		pos := [3]float64{}
		idx := 0
		for s := range strings.SplitSeq(line, ",") {
			n, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return nil, nil
			}

			pos[idx] = n

			idx += 1
		}

		boxes = append(boxes, pos)
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("scan error: %w", scanner.Err())
	}

	return boxes, nil
}

type connection struct {
	boxI int
	boxJ int
	dist float64
}

func process(boxes [][3]float64) int64 {
	cableCount := 1000
	if len(boxes) <= 20 {
		// Use fewer connections for the small example input
		cableCount = 10
	}

	// Compute the distance between all possible connections (with j>i)
	connections := make([]connection, 0, len(boxes)*len(boxes))
	for i := 0; i < len(boxes); i += 1 {
		p0 := boxes[i]

		for j := i + 1; j < len(boxes); j += 1 {
			p1 := boxes[j]
			d := distSqr(p0, p1)
			connections = append(connections, connection{
				boxI: i,
				boxJ: j,
				dist: d,
			})
		}
	}

	// Sort the connections in ascending order by distance
	// (This could get optimized using quick select, since we are only
	// interested in the N shortest distances)
	slices.SortFunc(connections, func(a, b connection) int {
		if a.dist < b.dist {
			return -1
		}
		if a.dist > b.dist {
			return 1
		}
		return 0
	})

	// Keep the top-N shortest distances
	connections = connections[:cableCount]

	// Map from box index to cluster ID and from cluster ID to box indices.
	// Initially each box is its own cluster.
	boxToCluster := make(map[int]int, len(boxes))
	clusterToBoxes := make(map[int][]int, len(boxes))
	for i := range boxes {
		boxToCluster[i] = i
		clusterToBoxes[i] = []int{i}
	}

	for _, c := range connections {
		clusterI := boxToCluster[c.boxI]
		clusterJ := boxToCluster[c.boxJ]
		mergeClusters(boxToCluster, clusterToBoxes, clusterI, clusterJ)
	}

	// Calculate the size of each cluster
	clusterSizes := make(map[int]int, len(boxToCluster))
	for _, cluster := range boxToCluster {
		clusterSizes[cluster] += 1
	}

	// Sort the cluster sizes to find the 3 largest
	// (This could get optimized using quick select, since we are only
	// interested in the 3 largest clusters)
	sortedSizes := make([]int, 0, len(clusterSizes))
	for _, size := range clusterSizes {
		sortedSizes = append(sortedSizes, size)
	}
	slices.SortFunc(sortedSizes, func(a, b int) int {
		return b - a
	})

	res := int64(1)
	for i := 0; i < 3; i += 1 {
		res *= int64(sortedSizes[i])
	}
	return res
}

func distSqr(p0, p1 [3]float64) float64 {
	dx := p0[0] - p1[0]
	dy := p0[1] - p1[1]
	dz := p0[2] - p1[2]
	return dx*dx + dy*dy + dz*dz
}

func mergeClusters(
	boxToCluster map[int]int,
	clusterToBoxes map[int][]int,
	clusterI int,
	clusterJ int,
) {
	newClusterID := min(clusterI, clusterJ)
	oldClusterID := max(clusterI, clusterJ)

	// Skip identical clusters
	if newClusterID == oldClusterID {
		return
	}

	// Move boxes to the new cluster
	oldClusterBoxes := clusterToBoxes[oldClusterID]
	for _, box := range oldClusterBoxes {
		boxToCluster[box] = newClusterID
	}
	clusterToBoxes[newClusterID] = append(
		clusterToBoxes[newClusterID],
		oldClusterBoxes...,
	)

	// "Delete" the old cluster
	delete(clusterToBoxes, oldClusterID)
}
