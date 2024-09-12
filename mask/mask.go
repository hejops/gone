// Package mask provides operations to extract and manipulate ranges of bits
// from a byte.
//
// All byte indices must be 1-indexed, and ranges must be inclusive.

package mask

import (
	_bits "math/bits"
)

// A byteIndex provides compile-time safety when indexing into a byte.
type byteIndex byte

const (
	I1 byteIndex = iota + 1
	I2
	I3
	I4
	I5
	I6
	I7
	I8
)

// https://pkg.go.dev/golang.org/x/text/internal/gen/bitfield
// https://cs.opensource.google/go/x/text/+/refs/tags/v0.18.0:internal/gen/bitfield/bitfield_test.go;l=16

// func checkByteIndex(n byteIndex) {
// 	// https://github.com/golang/go/issues/29649#issuecomment-454585328
// 	// https://github.com/golang/go/issues/29649#issuecomment-454820179
// 	//
// 	// Go does not allow us to model a constrained int with a type, hence
// 	// this helper func
// 	if n < 1 || n > 8 {
// 		panic("Invalid byte index provided -- must fall in the range [1,8].")
// 	}
// }

func checkByteRange(start byteIndex, end byteIndex) {
	if start > end {
		panic("Invalid range provided -- start must <= end.")
	}
}

// Last extracts the last n bits of b.
func Last(b byte, n byteIndex) byte {
	// this and lastLoop are about 0.0000015 ns/op, in the worst case

	// https://stackoverflow.com/a/15255834
	return b & ((1 << n) - 1)
}

func lastLoop(b byte, n byteIndex) byte {
	var last byte
	for bit := range n {
		last += (1 << bit)
	}
	return b & last
}

// First extracts the first n bits of b.
func First(b byte, n byteIndex) byte {
	// push the bits down, then apply the mask as usual
	return Last(b>>(8-n), n)
	// var first byte
	// for bit := range n {
	// 	first += (1 << bit)
	// }
	// return (b >> (8 - n)) & (first)
}

// Range extracts the inclusive range of bits [start:end] from b. Both start
// and end are 1-indexed.
func Range(b byte, start byteIndex, end byteIndex) byte {
	checkByteRange(start, end)
	// 0b1101_1000, 4, 5
	//      L_LLLL
	//      F_F
	tail := Last(b, 8-(start-1))
	return First(tail, end)
}

// IsSet reports whether the bit at pos is 1.
func IsSet(b byte, pos byteIndex) bool {
	return b&(1<<(8-pos)) != 0
}

// Set replaces the existing bits of b at pos (1-indexed) with new bits.
//
// If the new bits are zero, b is returned unchanged; Unset should be used to
// clear bits.
//
// If the new bits cannot fit at the desired pos, the new bits will be
// truncated.
func Set(b byte, pos byteIndex, bits byte) byte {
	if bits == 0 {
		return b
	}
	bitlen := byte(_bits.LeadingZeros8(bits))
	bits <<= bitlen
	bits >>= pos - 1
	return b | bits
}

// Unset clears the existing bits of b in the inclusive range [start:end].
func Unset(b byte, start byteIndex, end byteIndex) byte {
	checkByteRange(start, end)
	for ; start <= end; start++ {
		// hole := byte(math.MaxUint8 - 1<<(8-start))
		hole := byte(^(1 << byte(8-start))) // a full byte, with 1 bit unset
		b &= hole
	}
	return b
}

// Flip flips the existing bits of b in the inclusive range [start:end].
func Flip(b byte, start byteIndex, end byteIndex) byte {
	checkByteRange(start, end)
	for ; start <= end; start++ {
		b ^= (1 << (8 - start))
	}
	return b
}
