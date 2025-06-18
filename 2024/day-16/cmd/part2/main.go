package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
)

type priorityItem struct {
	priority int
	pos      [2]int
	dirIndex int
	parent   *priorityItem
}

type priorityQueue []*priorityItem

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x any) {
	*pq = append(*pq, x.(*priorityItem))
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(*pq)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	return item
}

type path [][2]int

func main() {
	if err := run(); err != nil {
		fmt.Printf("run failed: %s", err.Error())
		os.Exit(1)
	}
}

func run() error {
	// file, err := os.Open("./assets/input_tiny.txt") // 45
	// file, err := os.Open("./assets/input_small.txt") // 64
	file, err := os.Open("./assets/input.txt")
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	m := [][]byte{}
	startPos := [2]int{}
	endPos := [2]int{}

	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		row := []byte{}
		for x := range scanner.Text() {
			ch := scanner.Text()[x]
			if ch == 'S' {
				startPos = [2]int{x, y}
				ch = '.'
			} else if ch == 'E' {
				endPos = [2]int{x, y}
				ch = '.'
			}
			row = append(row, ch)
		}
		m = append(m, row)
		y += 1
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// fmt.Printf("m=%v, startPos=%v, endPos=%v\n", m, startPos, endPos)

	bestPaths := minCostPaths(m, startPos, endPos)
	// fmt.Printf("bestPaths=%v\n", bestPaths)

	visitedMap := copyMap(m)
	for _, path := range bestPaths {
		for _, pos := range path {
			visitedMap[pos[1]][pos[0]] = 'O'
		}
	}

	// printMap(visitedMap)
	// fmt.Printf("\n")

	count := countVisited(visitedMap)
	fmt.Printf("%d", count)

	return nil
}

func minCostPaths(m [][]byte, startPos, endPos [2]int) []path {
	h := len(m)
	w := len(m[0])
	visited := make([][][4]int, h)
	for y := range visited {
		visited[y] = make([][4]int, w)
		for x := range visited[y] {
			visited[y][x][0] = math.MaxInt32
			visited[y][x][1] = math.MaxInt32
			visited[y][x][2] = math.MaxInt32
			visited[y][x][3] = math.MaxInt32
		}
	}

	dirs := [4][2]int{
		[2]int{1, 0},
		[2]int{0, 1},
		[2]int{-1, 0},
		[2]int{0, -1},
	}

	priorityQueue := make(priorityQueue, 0, 1024)
	priorityQueue = append(priorityQueue, &priorityItem{
		priority: 0,
		pos:      startPos,
		dirIndex: 0,
		parent:   nil,
	})

	bestItems := []*priorityItem{}
	for priorityQueue.Len() > 0 {
		item := heap.Pop(&priorityQueue).(*priorityItem)

		cost := item.priority
		pos := item.pos
		dirIndex := item.dirIndex

		if len(bestItems) > 0 && cost > bestItems[0].priority {
			break // No better solutions can be found
		}

		if pos == endPos {
			bestItems = append(bestItems, item)
			continue
		}

		if pos[0] < 0 || pos[1] >= w || pos[1] < 0 || pos[1] >= h {
			continue
		}

		if m[pos[1]][pos[0]] == '#' {
			continue
		}

		if cost > visited[pos[1]][pos[0]][dirIndex] {
			continue
		}

		visited[pos[1]][pos[0]][dirIndex] = cost

		walkPos := [2]int{
			pos[0] + dirs[dirIndex][0],
			pos[1] + dirs[dirIndex][1],
		}
		heap.Push(&priorityQueue, &priorityItem{
			priority: cost + 1,
			pos:      walkPos,
			dirIndex: dirIndex,
			parent:   item,
		})

		cwDirIndex := (dirIndex + 1) % 4
		heap.Push(&priorityQueue, &priorityItem{
			priority: cost + 1000,
			pos:      pos,
			dirIndex: cwDirIndex,
			parent:   item,
		})

		ccwDirIndex := (4 + dirIndex - 1) % 4
		heap.Push(&priorityQueue, &priorityItem{
			priority: cost + 1000,
			pos:      pos,
			dirIndex: ccwDirIndex,
			parent:   item,
		})
	}

	bestPaths := []path{}
	for _, item := range bestItems {
		bestPaths = append(bestPaths, getPathFromItem(item))
	}
	return bestPaths
}

func getPathFromItem(item *priorityItem) path {
	p := path{}
	node := item
	for node != nil {
		p = append(p, node.pos)
		node = node.parent
	}
	return p
}

func copyMap(m [][]byte) [][]byte {
	h := len(m)
	w := len(m[0])
	mm := make([][]byte, h)
	for i := 0; i < h; i += 1 {
		row := make([]byte, w)
		copy(row, m[i])
		mm[i] = row
	}
	return mm
}

func printMap(m [][]byte) {
	for _, row := range m {
		fmt.Printf("%s\n", string(row))
	}
}

func countVisited(m [][]byte) int {
	count := 0
	for _, row := range m {
		for _, v := range row {
			if v == 'O' {
				count += 1
			}
		}
	}
	return count
}
