package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
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

func main() {
	if err := run(); err != nil {
		fmt.Printf("run failed: %s\n", err.Error())
		os.Exit(1)
	}
}

type path []string

func run() error {
	// file, err := os.Open("./assets/input_tiny.txt")
	// file, err := os.Open("./assets/input_small2.txt")
	// file, err := os.Open("./assets/input_small.txt")
	file, err := os.Open("./assets/input.txt")
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v", err)
	}
	defer file.Close()

	inputs := make(map[string]uint64, 1000)
	gates := make(map[string]gate, 1000)
	outputs := make([]string, 0, 100)
	mode := 0
	scanner := bufio.NewScanner(file)

	swaps := func() map[string]string {
		ss := []struct {
			a string
			b string
		}{
			{"z18", "hmt"},
			{"z27", "bfq"},
			{"z31", "hkh"},
			{"bng", "fjp"},
		}

		m := make(map[string]string, 2*len(ss))
		for _, s := range ss {
			m[s.a] = s.b
			m[s.b] = s.a
		}
		return m
	}()

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

			// To simplify the future processing, we ensure, that x inputs
			// always are the left input, while y inputs are always the right
			// input.
			if strings.HasPrefix(rhs, "x") || strings.HasPrefix(lhs, "y") {
				lhs, rhs = rhs, lhs
			}

			if swap, found := swaps[out]; found {
				out = swap
			}

			gate := gate{
				kind: kind,
				lhs:  lhs,
				rhs:  rhs,
				out:  out,
			}
			gates[out] = gate

			if strings.HasPrefix(out, "z") {
				outputs = append(outputs, out)
			}
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	slices.Sort(outputs)

	// fmt.Printf("inputs=%v\n", inputs)
	// fmt.Printf("gates=%v\n", gates)
	// fmt.Printf("outputs=%v\n", outputs)

	// if err := generateDiagrams(gates, inputs); err != nil {
	// 	return fmt.Errorf("generating diagrams: %w", err)
	// }

	// if err := generateExpressions(gates, inputs, outputs); err != nil {
	// 	return fmt.Errorf("generate expressions: %w", err)
	// }

	for i, output := range outputs {
		valid := isValidOutput(gates, inputs, output, i)
		fmt.Printf("%s valid: %v\n", output, valid)
	}

	invalidGates := make([]string, 0, 8)
	for s := range swaps {
		invalidGates = append(invalidGates, s)
	}
	slices.Sort(invalidGates)
	result := strings.Join(invalidGates, ",")
	fmt.Printf("%v\n", result)

	return nil
}

func generateDiagrams(
	gates map[string]gate,
	inputs map[string]uint64,
) error {
	deps := make(map[string][]string, len(gates))
	for _, g := range gates {
		deps[g.out] = append(deps[g.out], g.lhs)
		deps[g.out] = append(deps[g.out], g.rhs)
	}

	var dfs func(name string) []path
	dfs = func(name string) []path {
		// Input node?
		if _, found := inputs[name]; found {
			p := []string{name}
			return []path{p}
		}

		res := []path{}
		for _, d := range deps[name] {
			paths := dfs(d)
			for _, p := range paths {
				p = append(p, name)
				res = append(res, p)
			}
		}
		return res
	}
	outPaths := map[string][]path{}
	for _, g := range gates {
		if strings.HasPrefix(g.out, "z") {
			outPaths[g.out] = dfs(g.out)
		}
	}
	// fmt.Printf("outPaths=%v\n", outPaths)

	for out, paths := range outPaths {
		fmt.Printf("Generating Mermaid file for '%s' ...\n", out)

		mermaidBuf := strings.Builder{}
		mermaidBuf.WriteString("flowchart TD\n")
		drawn := map[string]struct{}{}
		for _, path := range paths {
			out := path[len(path)-1]
			for _, n := range path {
				if _, isGate := gates[n]; isGate {
					g := gates[n]
					mermaidBuf.WriteString(" --> ")
					gateName := fmt.Sprintf("g_%s_%s{{%s}}", out, n, g.kind)
					mermaidBuf.WriteString(gateName)

					if _, isDrawn := drawn[gateName]; isDrawn {
						break
					}
					drawn[gateName] = struct{}{}
					mermaidBuf.WriteString(" --> ")
				} else {
					mermaidBuf.WriteString("    ")
				}
				mermaidBuf.WriteString(n)
			}
			mermaidBuf.WriteString("\n")
		}

		err := func() error {
			f, err := os.Create(fmt.Sprintf("./output/%s.mmd", out))
			if err != nil {
				return fmt.Errorf("create file '%s': %w", out, err)
			}
			defer func() {
				if err := f.Close(); err != nil {
					fmt.Printf("close file '%s': %v", out, err)
				}
			}()

			_, err = f.WriteString(mermaidBuf.String())
			if err != nil {
				return fmt.Errorf("write string to file '%s': %w", out, err)
			}

			return nil
		}()
		if err != nil {
			return err
		}
	}

	for out := range outPaths {
		fmt.Printf("Generating PDF for '%s' ...\n", out)

		input := fmt.Sprintf("./output/%s.mmd", out)
		output := fmt.Sprintf("./output/%s.pdf", out)
		cmd := exec.Command(
			"mmdc",
			"--configFile",
			"./config.json",
			"--input",
			input,
			"--output",
			output,
		)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("command '%v': %w", cmd, err)
		}
	}

	return nil
}

func generateExpressions(
	gates map[string]gate,
	inputs map[string]uint64,
	outputs []string,
) error {
	for _, output := range outputs {
		expr := getExpression(gates, inputs, output)
		if err := os.WriteFile(
			fmt.Sprintf("./output/%s.txt", output),
			[]byte(expr),
			0644,
		); err != nil {
			return err
		}
	}
	return nil
}

func getExpression(
	gates map[string]gate,
	inputs map[string]uint64,
	root string,
) string {
	var dfs func(node string) string
	dfs = func(node string) string {
		if _, found := inputs[node]; found {
			return node
		}

		gate := gates[node]
		lhsExpr := dfs(gate.lhs)
		rhsExpr := dfs(gate.rhs)

		return fmt.Sprintf(
			"(%s %s %s)",
			lhsExpr,
			string(gate.kind),
			rhsExpr,
		)
	}

	return dfs(root)
}

func isValidOutput(
	gates map[string]gate,
	inputs map[string]uint64,
	output string,
	n int,
) bool {
	// Carry:
	// c_n = (x_n AND y_n) OR (c_{n-1} AND (x_n XOR y_n))
	// c_0 = x_0 AND y_0
	//
	// Sum:
	// z_end = c_{end-1}
	// z_n = c_{n-1} XOR (x_n XOR y_n)
	// z_0 = x_0 XOR y_0

	var isValidCarry func(root string, n int) bool
	isValidCarry = func(root string, n int) bool {
		g := gates[root]

		// c_0 = x_0 AND y_0
		if n == 0 {
			expectedLhs := fmt.Sprintf("x%02d", n)
			expectedRhs := fmt.Sprintf("y%02d", n)

			if g.kind != andGate {
				fmt.Printf(
					"invalid: %s (kind should be AND, but is %s)\n",
					root,
					string(g.kind),
				)
				return false
			}
			if g.lhs != expectedLhs {
				fmt.Printf(
					"invalid: %s (lhs should be '%s', but is %s)\n",
					root,
					expectedLhs,
					g.lhs,
				)
				return false
			}
			if g.rhs != expectedRhs {
				fmt.Printf(
					"invalid: %s (rhs should be '%s', but is %s)\n",
					root,
					expectedLhs,
					g.rhs,
				)
				return false
			}
			return true
		}

		// c_n = (x_n AND y_n) OR (c_{n-1} AND (x_n XOR y_n))

		if g.kind != orGate {
			fmt.Printf(
				"invalid: %s (kind should be OR, but is %s)\n",
				g.out,
				string(g.kind),
			)
			return false
		}
		lhsGate := gates[g.lhs]
		rhsGate := gates[g.rhs]
		if !strings.HasPrefix(lhsGate.lhs, "x") {
			lhsGate, rhsGate = rhsGate, lhsGate
		}

		// (x_n AND y_n)
		lhsValid := func(node gate) bool {
			expectedLhs := fmt.Sprintf("x%02d", n)
			expectedRhs := fmt.Sprintf("y%02d", n)

			if node.kind != andGate {
				fmt.Printf(
					"invalid: %s (kind should be AND, but is %s)\n",
					node.out,
					string(node.kind),
				)
				return false
			}
			if node.lhs != expectedLhs {
				fmt.Printf(
					"invalid: %s (lhs should be '%s', but is %s)\n",
					node.out,
					expectedLhs,
					node.lhs,
				)
				return false
			}
			if node.rhs != expectedRhs {
				fmt.Printf(
					"invalid: %s (rhs should be '%s', but is %s)\n",
					node.out,
					expectedLhs,
					node.rhs,
				)
				return false
			}
			return true
		}(lhsGate)
		if !lhsValid {
			return false
		}

		// (c_{n-1} AND (x_n XOR y_n))
		rhsValid := func(node gate) bool {
			lhsGate := gates[node.lhs]
			rhsGate := gates[node.rhs]
			if !strings.HasPrefix(rhsGate.lhs, "x") {
				lhsGate, rhsGate = rhsGate, lhsGate
			}

			if node.kind != andGate {
				fmt.Printf(
					"invalid: %s (kind should be AND, but is %s)\n",
					node.out,
					string(node.kind),
				)
				return false
			}
			if rhsGate.kind != xorGate {
				fmt.Printf(
					"invalid: %s (kind should be XOR, but is %s)\n",
					rhsGate.out,
					string(rhsGate.kind),
				)
				return false
			}

			expectedLhs := fmt.Sprintf("x%02d", n)
			expectedRhs := fmt.Sprintf("y%02d", n)
			if rhsGate.lhs != expectedLhs {
				fmt.Printf(
					"invalid: %s (lhs should be '%s', but is %s)\n",
					rhsGate.out,
					expectedLhs,
					rhsGate.lhs,
				)
				return false
			}
			if rhsGate.rhs != expectedRhs {
				fmt.Printf(
					"invalid: %s (rhs should be '%s', but is %s)\n",
					rhsGate,
					expectedLhs,
					rhsGate.rhs,
				)
				return false
			}
			return isValidCarry(lhsGate.out, n-1)
		}(rhsGate)
		if !rhsValid {
			return false
		}

		return true
	}

	g := gates[output]

	// z_0 = x_0 XOR y_0
	if n == 0 {
		return g.kind == xorGate &&
			g.lhs == fmt.Sprintf("x%02d", n) &&
			g.rhs == fmt.Sprintf("y%02d", n)
	}

	// z_end = c_{end-1}
	end := 45
	if n >= end {
		return isValidCarry(g.out, end-1)
	}

	// z_n = c_{n-1} XOR (x_n XOR y_n)
	lhsGate := gates[g.lhs]
	rhsGate := gates[g.rhs]
	if !strings.HasPrefix(rhsGate.lhs, "x") {
		lhsGate, rhsGate = rhsGate, lhsGate
	}
	if g.kind != xorGate {
		fmt.Printf(
			"invalid: %s (kind should be XOR, but is %s)\n",
			g.out,
			string(g.kind),
		)
		return false
	}
	expectedLhs := fmt.Sprintf("x%02d", n)
	expectedRhs := fmt.Sprintf("y%02d", n)
	if rhsGate.lhs != expectedLhs {
		fmt.Printf(
			"invalid: %s (lhs should be '%s', but is '%s')\n",
			rhsGate.out,
			expectedLhs,
			rhsGate.lhs,
		)
		return false
	}
	if rhsGate.rhs != expectedRhs {
		fmt.Printf(
			"invalid: %s (rhs should be '%s', but is '%s')\n",
			rhsGate.out,
			expectedRhs,
			rhsGate.rhs,
		)
		return false
	}
	return isValidCarry(lhsGate.out, n-1)
}
