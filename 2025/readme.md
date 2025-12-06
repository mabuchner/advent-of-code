# Advent of Code 2025

My solutions for [Advent of Code 2025](https://adventofcode.com/2025) written in Go.

```sh
go test ./...
go test -bench=. ./...
go run ./cmd/part1 | pbcopy
go run ./cmd/part2 | pbcopy
```

Use day 1 as a template for a new day

```sh
cp -R day-01 day-02 && cd ./day-02 && sed -i '' -e 's/day-01/day-02/g' go.mod
```
