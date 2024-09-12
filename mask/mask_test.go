package mask

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMask(t *testing.T) {
	assert.Equal(t, Last(0b0000_1111, I1), byte(0b0000_0001))
	assert.Equal(t, Last(0b0000_1111, I2), byte(0b0000_0011))
	assert.Equal(t, Last(0b0000_1111, I3), byte(0b0000_0111))
	assert.Equal(t, Last(0b0000_1111, I4), byte(0b0000_1111))

	assert.Equal(t, Last(0b1000_1111, I1), byte(0b0000_0001))
	assert.Equal(t, Last(0b1000_1111, I2), byte(0b0000_0011))
	assert.Equal(t, Last(0b1000_1111, I3), byte(0b0000_0111))
	assert.Equal(t, Last(0b1000_1111, I4), byte(0b0000_1111))

	assert.Equal(t, Last(0b0000_1010, I1), byte(0b0000_0000))
	assert.Equal(t, Last(0b0000_1010, I2), byte(0b0000_0010))
	assert.Equal(t, Last(0b0000_1010, I3), byte(0b0000_0010))
	assert.Equal(t, Last(0b0000_1010, I4), byte(0b0000_1010))

	assert.Equal(t, First(0b1111_1111, 1), byte(0b0000_0001))
	assert.Equal(t, First(0b1010_1111, 4), byte(0b0000_1010))

	assert.Equal(t, Range(0b1101_1000, I1, I2), byte(0b0000_0011))
	assert.Equal(t, Range(0b1101_1000, I2, I4), byte(0b0000_0101))
	assert.Equal(t, Range(0b1101_1000, I4, I5), byte(0b0000_0011))
	assert.Equal(t, Range(0b1101_1000, I5, I8), byte(0b0000_1000))

	assert.True(t, IsSet(0b1101_1000, 1))
	assert.True(t, IsSet(0b1101_1000, 2))
	assert.False(t, IsSet(0b1101_1000, 3))
	assert.True(t, IsSet(0b1101_1000, 4))

	assert.Equal(t, Set(0b0000_0000, 1, 0b0000_0010), byte(0b1000_0000))
	assert.Equal(t, Set(0b0000_0000, 1, 0b0000_0101), byte(0b1010_0000))
	assert.Equal(t, Set(0b0000_0000, 1, 0b0000_0111), byte(0b1110_0000))
	assert.Equal(t, Set(0b0000_0000, 2, 0b0000_0011), byte(0b0110_0000))
	assert.Equal(t, Set(0b0000_0000, 2, 0b0000_0111), byte(0b0111_0000))
	assert.Equal(t, Set(0b0000_0000, 5, 0b0000_1111), byte(0b0000_1111))
	assert.Equal(t, Set(0b0000_0000, 7, 0b0000_1000), byte(0b0000_0010))
	assert.Equal(t, Set(0b0000_0000, 7, 0b0000_1111), byte(0b0000_0011))
	assert.Equal(t, Set(0b1111_1111, 1, 0), byte(0b1111_1111))

	assert.Equal(t, Unset(0b1111_0000, 5, 8), byte(0b1111_0000))
	assert.Equal(t, Unset(0b1111_1111, 5, 8), byte(0b1111_0000))

	assert.Equal(t, Flip(0b1111_0000, 5, 5), byte(0b1111_1000))
	assert.Equal(t, Flip(0b1111_0000, 5, 8), byte(0b1111_1111))
	assert.Equal(t, Flip(0b1111_0000, 8, 8), byte(0b1111_0001))
	assert.Equal(t, Flip(0b1111_1111, 5, 8), byte(0b1111_0000))

	// assert.Panics(t, func() { _ = Last(byte(0), 10) })
	// assert.Panics(t, func() { _ = Range(byte(0), 0, 9) })
}

func BenchmarkLast(b *testing.B) {
	Last(0b1000_1111, 4)
}

func BenchmarkLastLoop(b *testing.B) {
	lastLoop(0b1000_1111, 4)
}

func BenchmarkFirst(b *testing.B) {
	First(0b1000_1111, 4)
}
