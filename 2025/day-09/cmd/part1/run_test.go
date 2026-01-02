package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Run("small input", func(t *testing.T) {
		res, err := run("../../assets/input_small.txt")
		assert.NoError(t, err)
		assert.Equal(t, int64(50), res)
	})

	t.Run("normal input", func(t *testing.T) {
		res, err := run("../../assets/input.txt")
		assert.NoError(t, err)
		assert.Equal(t, int64(4745816424), res)
	})
}

func BenchmarkProcess(b *testing.B) {
	inputs, err := load("../../assets/input.txt")
	if err != nil {
		b.FailNow()
	}

	for b.Loop() {
		process(inputs)
	}
}
