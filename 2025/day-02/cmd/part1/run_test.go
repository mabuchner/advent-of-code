package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Run("basics", func(t *testing.T) {
		assert.Equal(t, int64(6), decLen(123123))
		assert.Equal(t, int64(1000), pow10[3])
		assert.Equal(t, int64(123), int64(123123)/1000)
		assert.Equal(t, int64(123), int64(123123)%1000)
	})

	t.Run("small input", func(t *testing.T) {
		res, err := run("../../assets/input_small.txt")
		assert.NoError(t, err)
		assert.Equal(t, int64(1227775554), res)
	})

	t.Run("normal input", func(t *testing.T) {
		res, err := run("../../assets/input.txt")
		assert.NoError(t, err)
		assert.Equal(t, int64(20223751480), res)
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
