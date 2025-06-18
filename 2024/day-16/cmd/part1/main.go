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

func main() {
	if err := run(); err != nil {
		fmt.Printf("run failed: %s", err.Error())
		os.Exit(1)
	}
}

func run() error {
	// file, err := os.Open("./assets/input_tiny.txt") // 7036
	// file, err := os.Open("./assets/input_small.txt") // 11048
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

	cost := minCost(m, startPos, endPos)
	fmt.Printf("%d", cost)

	return nil
}

func minCost(m [][]byte, startPos, endPos [2]int) int {
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
	})

	// i := 0

	for priorityQueue.Len() > 0 {
		item := heap.Pop(&priorityQueue).(*priorityItem)

		cost := item.priority
		pos := item.pos
		dirIndex := item.dirIndex

		// if i%10000000 == 0 {
		// 	fmt.Printf("pos=%v, cost=%d\n", pos, cost)
		// }
		// i += 1

		if pos[0] < 0 || pos[1] >= w || pos[1] < 0 || pos[1] >= h {
			continue
		}

		if m[pos[1]][pos[0]] == '#' {
			continue
		}

		if cost > visited[pos[1]][pos[0]][dirIndex] {
			continue
		}

		if pos == endPos {
			return cost
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
		})

		cwDirIndex := (dirIndex + 1) % 4
		heap.Push(&priorityQueue, &priorityItem{
			priority: cost + 1000,
			pos:      pos,
			dirIndex: cwDirIndex,
		})

		ccwDirIndex := (4 + dirIndex - 1) % 4
		heap.Push(&priorityQueue, &priorityItem{
			priority: cost + 1000,
			pos:      pos,
			dirIndex: ccwDirIndex,
		})
	}

	return -1
}
