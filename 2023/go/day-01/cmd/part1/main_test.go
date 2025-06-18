package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	lines := []string{
		"1abc2",       // 12
		"pqr3stu8vwx", // 38
		"a1b2c3d4e5f", // 15
		"treb7uchet",  // 77
	}

	sum, err := process(lines)

	assert.NoError(t, err)
	assert.Equal(t, 142, sum)
}
