package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTimeInSeconds(tst *testing.T) {
	assert.Equal(tst, 10*time.Second, GetTimeInSeconds(10))
}

func TestGetTime(tst *testing.T) {
	time, err := getTime("America/Sao_Paulo", "Mon, 02 Jan 2006 15:04:05 MST")
	assert.NoError(tst, err)
	assert.GreaterOrEqual(tst, len(time), 10)
}

func TestToFixed(tst *testing.T) {
	num := 13.337333993
	assert.Equal(tst, 13.3373, ToFixed(num))
}
