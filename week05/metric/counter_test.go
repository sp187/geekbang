package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRing(t *testing.T) {
	r := newRing(3)
	r.add(1)
	r.next()
	r.add(2)
	r.next()
	r.add(3)
	assert.Equal(t, int64(6), r.sum)
	r.next()
	assert.Equal(t, int64(6-1), r.sum)
	r.next() // -2
	r.add(4) // + 4
	assert.Equal(t, int64(5-2+4), r.sum)
	// [3, 0, 4]
	r.next()
	// [0, 4, 0]
	r.next()
	// [4, 0, 0]
	assert.Equal(t, int64(4), r.sum)

}
