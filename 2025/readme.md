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

When adding a new day, update the 'day' arrays in the
[2025-ci.yml](../.github/workflows/2025-ci.yml) file for each job.

Example

```yml
jobs:
  # ...
  matrix:
    day: ['01', '02', '03', '04'] # <--
```

## Notes

### Day 9

The solution for day 9 part 2 contains code to generate a SVG file of the
input and the result. To enable the SVG generation set the
`generateSVGEnabled` variable to `true`. The `process` function will then
create a file called `out.svg`.
