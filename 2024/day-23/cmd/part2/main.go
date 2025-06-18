package main

import (
	"bufio"
	"fmt"
	"maps"
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
	// file, err := os.Open("./assets/input_small.txt") // co,de,ka,ta
	file, err := os.Open("./assets/input.txt")
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

	adj := make(map[string]map[string]struct{}, len(connections)*2)
	for _, c := range connections {
		if adj[c[0]] == nil {
			adj[c[0]] = map[string]struct{}{}
		}
		if adj[c[1]] == nil {
			adj[c[1]] = map[string]struct{}{}
		}
		adj[c[0]][c[1]] = struct{}{}
		adj[c[1]][c[0]] = struct{}{}
	}
	// fmt.Printf("adj=%v\n", adj)

	largestCluster := findLargestClique(adj)
	// fmt.Printf("largestCluster=%v (len=%d)\n", largestCluster, len(largestCluster))

	nodes := make([]string, 0, len(largestCluster))
	for n := range largestCluster {
		nodes = append(nodes, n)
	}
	sort.Strings(nodes)

	password := strings.Join(nodes, ",")
	fmt.Print(password)

	return nil
}

func findLargestClique(adj map[string]map[string]struct{}) map[string]struct{} {
	var largestClique map[string]struct{}

	// R (Current Clique): A set representing the current clique being
	//                     constructed.
	// P (Potential Candidates): A set of vertices that can be added to R to
	//                           form a larger clique.
	// X (Excluded Vertices): A set of vertices that have already been
	//                        processed and should not be reconsidered for the
	//                        current clique.
	var bronKerbosch func(R, P, X map[string]struct{})
	bronKerbosch = func(R, P, X map[string]struct{}) {
		if len(P) == 0 && len(X) == 0 {
			if len(R) > len(largestClique) {
				largestClique = maps.Clone(R)
			}
			return
		}

		for v := range P {
			// R ⋃ {v}
			newR := maps.Clone(R)
			newR[v] = struct{}{}

			// Neighbours of v
			Nv := adj[v]

			// P ⋂ N(v)
			newP := make(map[string]struct{}, len(Nv))
			for n := range Nv {
				if _, ok := P[n]; ok {
					newP[n] = struct{}{}
				}
			}

			// X ⋂ N(v)
			newX := make(map[string]struct{}, len(Nv))
			for n := range Nv {
				if _, ok := X[n]; ok {
					newX[n] = struct{}{}
				}
			}

			bronKerbosch(newR, newP, newX)

			delete(P, v)      // P := P \ {v}
			X[v] = struct{}{} // X := X ⋃ {v}
		}
	}

	R := make(map[string]struct{}, len(adj))

	P := make(map[string]struct{}, len(adj))
	for v := range adj {
		P[v] = struct{}{}
	}

	X := make(map[string]struct{}, len(adj))

	bronKerbosch(R, P, X)

	return largestClique
}
