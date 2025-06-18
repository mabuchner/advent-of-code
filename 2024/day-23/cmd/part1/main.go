package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("run failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	// file, err := os.Open("./assets/input_small.txt") // 7
	file, err := os.Open("./assets/input.txt") // 1054
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	connections := [][2]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		c := strings.SplitN(scanner.Text(), "-", 2)
		connections = append(connections, [2]string{c[0], c[1]})
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	// fmt.Printf("connections=%v\n", connections)

	adj := make(map[string][]string, len(connections)*2)
	for _, c := range connections {
		adj[c[0]] = append(adj[c[0]], c[1])
		adj[c[1]] = append(adj[c[1]], c[0])
	}
	// fmt.Printf("adj=%v\n", adj)

	triplets := make(map[[3]string]struct{}, len(adj))
	for c1, connections := range adj {
		if !strings.HasPrefix(c1, "t") {
			continue
		}
		for _, c2 := range connections {
			connections2 := adj[c2]
			for _, c3 := range connections2 {
				connections4 := adj[c3]
				for _, c4 := range connections4 {
					if c1 == c4 {
						triplet := [3]string{c1, c2, c3}
						sort.Strings(triplet[:])
						triplets[triplet] = struct{}{}
					}
				}
			}
		}
	}
	// fmt.Printf("triplets=%v (len=%d)\n", triplets, len(triplets))
	fmt.Printf("%d", len(triplets))

	return nil
}
