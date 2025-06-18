package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("run failed: %s", err.Error())
		os.Exit(1)
	}
}

func run() error {
	file, err := os.Open("./assets/input.txt")
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	nums := [][2]int64{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		l := len(s)

		for start := 0; start < l; {
			r := tryParse(s, start)
			if r.found {
				// fmt.Printf(
				// 	"%s,%s,%s\n",
				// 	s[r.start:r.end],
				// 	s[r.numStart1:r.numEnd1],
				// 	s[r.numStart2:r.numEnd2],
				// )

				num1, err := strconv.ParseInt(s[r.numStart1:r.numEnd1], 10, 64)
				if err != nil {
					return err
				}

				num2, err := strconv.ParseInt(s[r.numStart2:r.numEnd2], 10, 64)
				if err != nil {
					return err
				}

				nums = append(nums, [2]int64{num1, num2})

				start = r.end
			} else {
				start += 1
			}
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	sum := int64(0)
	for _, n := range nums {
		sum += n[0] * n[1]
	}

	fmt.Printf("%d", sum)

	return nil
}

type parseResult struct {
	found     bool
	start     int
	end       int
	numStart1 int
	numEnd1   int
	numStart2 int
	numEnd2   int
}

func tryParse(s string, start int) parseResult {
	if start+7 > len(s) { // mul(a,b)
		return parseResult{}
	}

	if s[start:start+4] != "mul(" {
		return parseResult{}
	}

	numStart1 := start + 4
	if !isDigit(s[numStart1]) {
		return parseResult{}
	}
	numEnd1 := numStart1 + 1
	for ; numEnd1 < numStart1+3; numEnd1 += 1 {
		if isComma(s[numEnd1]) {
			break
		}
		if !isDigit(s[numEnd1]) {
			return parseResult{}
		}
	}
	if !isComma(s[numEnd1]) {
		return parseResult{}
	}

	numStart2 := numEnd1 + 1
	if !isDigit(s[numStart2]) {
		return parseResult{}
	}
	numEnd2 := numStart2 + 1
	for ; numEnd2 < numStart2+3; numEnd2 += 1 {
		if isRBrace(s[numEnd2]) {
			break
		}
		if !isDigit(s[numEnd2]) {
			return parseResult{}
		}
	}
	if !isRBrace(s[numEnd2]) {
		return parseResult{}
	}

	return parseResult{
		found:     true,
		start:     start,
		end:       numEnd2 + 1,
		numStart1: numStart1,
		numEnd1:   numEnd1,
		numStart2: numStart2,
		numEnd2:   numEnd2,
	}
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func isComma(b byte) bool {
	return b == ','
}

func isRBrace(b byte) bool {
	return b == ')'
}
