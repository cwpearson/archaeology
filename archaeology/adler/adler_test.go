package adler

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestRecurrence(t *testing.T) {
	data := []byte(`These are definitely some bytes`)

	a := NewSum(data[3:5], 3)
	b := NewSum(data[2:4], 2)
	b.Recurrence(data[4], data[2])

	assert.Equal(t, a.Current(), b.Current(), "Should be equal")

}
