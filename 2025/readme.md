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

## Code Structure

Each day's solution follows a consistent file structure:

```
day-XX/
├── assets/
│   ├── input.txt          # Actual puzzle input
│   └── input_small.txt    # Sample input from puzzle description
└── cmd/
    ├── part1/
    │   ├── main.go
    │   ├── run.go
    │   └── run_test.go
    └── part2/
        ├── main.go
        ├── run.go
        └── run_test.go
```

### Function Separation

Each solution is organized into three main functions:

- `main()`: Entry point that calls run and handles errors
- `run(inputPath string)`: Orchestrates the solution by calling load and process
- `load(inputPath string)`: Reads and parses the input file into a structured format
- `process(input)`: Performs the actual computation on the parsed data

Separating loading from processing enables accurate benchmarking of the
algorithm's performance. By loading the input once before the benchmark loop,
we can measure only the processing time without including file I/O overhead.

### Tests

Tests are implemented in the `run_test.go` files within each part's directory.
Each test file includes `TestRun` which verifies correctness using both sample
and actual inputs, and `BenchmarkProcess` which measures processing
performance.

## Notes

### Day 9

The solution for day 9 part 2 contains code to generate a SVG file of the
puzzle input and the result. To enable the SVG generation, set the
`generateSVGEnabled` variable to `true`. The `process` function will then
create a file called `out.svg`.
