package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpack(tst *testing.T) {
	var t1, t2 string
	input := []string{"1", "2"}
	Unpack(input, &t1, &t2)

	assert.Equal(tst, "1", t1)
	assert.Equal(tst, "2", t2)
}

func TestUnpackArraySmallerThanInput(tst *testing.T) {
	var t1, t2, t3 string
	input := []string{"1", "2"}
	Unpack(input, &t1, &t2, &t3)

	assert.Equal(tst, "1", t1)
	assert.Equal(tst, "2", t2)
	assert.Equal(tst, "", t3)
}

func TestUnpackArrayBiggerThanInput(tst *testing.T) {
	var t1, t2, t3 string
	input := []string{"1", "2", "3", "4"}
	Unpack(input, &t1, &t2, &t3)

	assert.Equal(tst, "1", t1)
	assert.Equal(tst, "2", t2)
	assert.Equal(tst, "3", t3)
}
