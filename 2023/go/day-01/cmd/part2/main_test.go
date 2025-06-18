package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	lines := []string{
		"two1nine",         // 29
		"eightwothree",     // 83
		"abcone2threexyz",  // 13
		"xtwone3four",      // 24
		"4nineeightseven2", // 42
		"zoneight234",      // 14
		"7pqrstsixteen",    // 76
	}

	sum, err := process(lines)

	assert.NoError(t, err)
	assert.Equal(t, 281, sum)

	sum, err = process([]string{"three2fiveonexrllxsvfive"})
	assert.NoError(t, err)
	assert.Equal(t, 35, sum)
}
