package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type gateKind string

const (
	andGate gateKind = "AND"
	orGate  gateKind = "OR"
	xorGate gateKind = "XOR"
)

type gate struct {
	kind gateKind
	lhs  string
	rhs  string
	out  string
}

func (kind gateKind) process(lhs, rhs uint64) uint64 {
	switch kind {
	case "AND":
		return lhs & rhs
	case "OR":
		return lhs | rhs
	case "XOR":
		return lhs ^ rhs
	}
	panic(fmt.Errorf("unexpected gate kind '%s'", kind))
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("run failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	// file, err := os.Open("./assets/input_tiny.txt") // 4
	// file, err := os.Open("./assets/input_small.txt") // 2024
	file, err := os.Open("./assets/input.txt") // 51745744348272
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	inputs := map[string]uint64{}
	gatesMap := map[string]gate{}
	mode := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Reading inputs
		if mode == 0 {
			if len(scanner.Text()) == 0 {
				mode = 1
				continue
			}

			nameAndValue := strings.SplitN(scanner.Text(), ": ", 2)
			if nameAndValue[1] == "0" {
				inputs[nameAndValue[0]] = 0
			} else {
				inputs[nameAndValue[0]] = 1
			}
		} else {
			gateDef := strings.SplitN(scanner.Text(), " ", 5)
			lhs := gateDef[0]
			kindStr := gateDef[1]
			rhs := gateDef[2]
			out := gateDef[4]

			var kind gateKind
			if kindStr == string(andGate) {
				kind = andGate
			} else if kindStr == string(orGate) {
				kind = orGate
			} else if kindStr == string(xorGate) {
				kind = xorGate
			} else {
				return fmt.Errorf("unexpected gate kind '%s'", kindStr)
			}

			gate := gate{
				kind: kind,
				lhs:  lhs,
				rhs:  rhs,
				out:  out,
			}
			gatesMap[out] = gate
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// fmt.Printf("inputs=%v\n", inputs)
	// fmt.Printf("gatesMap=%v\n", gatesMap)

	var dfs func(out string) uint64

	dfs = func(out string) uint64 {
		v, ok := inputs[out]
		if ok {
			return v
		}

		gate := gatesMap[out]
		lhsVal := dfs(gate.lhs)
		rhsVal := dfs(gate.rhs)

		res := gate.kind.process(lhsVal, rhsVal)

		inputs[gate.out] = res

		return res
	}

	result := uint64(0)
	for out := range gatesMap {
		if !strings.HasPrefix(out, "z") {
			continue
		}

		val := dfs(out)

		idx, err := strconv.ParseUint(out[1:], 10, 64)
		if err != nil {
			return err
		}

		result |= val << idx
	}

	fmt.Printf("%d", result)

	return nil
}
