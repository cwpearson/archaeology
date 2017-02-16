package adler

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecurrence(t *testing.T) {
	data1 := []byte(`These are definitely some bytes`)
	data2 := []byte(`These are similar         bytes`)

	windowSize := 3
	a := NewSum(data1[0:windowSize], 0)
	b := NewSum(data2[0:windowSize], 0)
	for s := 0; s < len(data2)-windowSize; s++ {
		e := s + windowSize

		if bytes.Equal(data1[s:e], data2[s:e]) {
			assert.Equal(t, a.Current(), b.Current(), "Expected hashes of \""+string(data1[s:e])+"\" to match")
		}

		fmt.Printf("%s %s %X %X %t\n", string(data1[s:e]), string(data2[s:e]), a.Current(), b.Current(), a.Current() == b.Current())

		a.Roll(data1[e])
		b.Roll(data2[e])
	}

}
