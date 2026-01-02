package main

import (
	"fmt"
	"os"
)

func main() {
	res, err := run("./assets/input.txt")
	if err != nil {
		fmt.Printf("run failed: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("%d", res)
}
