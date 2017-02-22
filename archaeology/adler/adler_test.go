package adler

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimilar(t *testing.T) {
	data1 := []byte(`These are definitely some bytes`)
	data2 := []byte(`These are similar         bytes`)

	windowSize := 3
	a := NewSum(data1[0:windowSize])
	b := NewSum(data2[0:windowSize])
	for s := 0; s < len(data2)-windowSize; s++ {
		e := s + windowSize

		if bytes.Equal(data1[s:e], data2[s:e]) {
			assert.Equal(t, a.Current(), b.Current(), "Expected hashes of \""+string(data1[s:e])+"\" to match")
		}

		// fmt.Printf("%s %s %X %X %t\n", string(data1[s:e]), string(data2[s:e]), a.Current(), b.Current(), a.Current() == b.Current())

		a.Roll(data1[e])
		b.Roll(data2[e])
	}

}

func TestRoll(t *testing.T) {
	data := []byte(`These are definitely some bytes`)

	windowSize := 3
	a := NewSum(data[3 : 3+windowSize])
	b := NewSum(data[0:windowSize])
	b.Roll(data[3])
	b.Roll(data[4])
	b.Roll(data[5])

	assert.Equal(t, a.Current(), b.Current(), "Expected hashes of \""+string(data[3:3+windowSize])+"\" to match")

}
