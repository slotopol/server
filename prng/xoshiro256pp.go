package prng

import (
	"encoding/binary"
	"math/bits"
	mr "math/rand/v2"
)

// Xoshiro256PP implements the math/rand/v2.Source interface.
// It occupies 32 bytes (exactly half of a CPU cache line).
// It is not thread-safe - use a separate instance for each goroutine.
type Xoshiro256PP struct {
	s0, s1, s2, s3 uint64
}

var _ mr.Source = (*Xoshiro256PP)(nil)

// NewXoshiro256PP initializes the generator from a single 64-bit seed
// using the SplitMix64 algorithm (to protect against a zero state).
// The algorithm must strictly not be initialized with zeros
// (the state [0, 0, 0, 0] will get stuck in an infinite loop of zeros).
func NewXoshiro256PP(seed uint64) *Xoshiro256PP {
	var sm = seed
	var nextSplitMix = func() uint64 {
		sm += 0x9e3779b97f4a7c15
		var z = sm
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		return z ^ (z >> 31)
	}

	return &Xoshiro256PP{
		s0: nextSplitMix(),
		s1: nextSplitMix(),
		s2: nextSplitMix(),
		s3: nextSplitMix(),
	}
}

// Uint64 generates the next 64-bit number in ~1 ns.
// It compiles down to simple register operations
// (XOR, shift, and rotate-left/ROL).
func (x *Xoshiro256PP) Uint64() uint64 {
	var res = bits.RotateLeft64(x.s0+x.s3, 23) + x.s0

	var t = x.s1 << 17

	x.s2 ^= x.s0
	x.s3 ^= x.s1
	x.s1 ^= x.s2
	x.s0 ^= x.s3

	x.s2 ^= t
	x.s3 = bits.RotateLeft64(x.s3, 45)

	return res
}

// Jump is equivalent to 2^128 calls to Uint64().
// It allows a single sequence to be partitioned into guaranteed
// non-overlapping streams.
func (x *Xoshiro256PP) Jump() {
	var jumpTable = [4]uint64{
		0x180ec6d33cfd0aba,
		0xd5a61266f0c9392c,
		0xa9582618e03fc9aa,
		0x39abdc4529b1661c,
	}

	var s0, s1, s2, s3 uint64
	for _, jumpVal := range jumpTable {
		for b := 0; b < 64; b++ {
			if (jumpVal & (uint64(1) << b)) != 0 {
				s0 ^= x.s0
				s1 ^= x.s1
				s2 ^= x.s2
				s3 ^= x.s3
			}
			x.Uint64()
		}
	}

	x.s0, x.s1, x.s2, x.s3 = s0, s1, s2, s3
}

// Clone creates an exact, independent copy of the generator's current state.
func (x *Xoshiro256PP) Clone() *Xoshiro256PP {
	var copyState = *x
	return &copyState
}

// AppendBinary implements the [encoding.BinaryAppender] interface.
func (x *Xoshiro256PP) AppendBinary(b []byte) ([]byte, error) {
	b = append(b, "xo256pp:"...)
	b = binary.LittleEndian.AppendUint64(b, x.s0)
	b = binary.LittleEndian.AppendUint64(b, x.s1)
	b = binary.LittleEndian.AppendUint64(b, x.s2)
	b = binary.LittleEndian.AppendUint64(b, x.s3)
	return b, nil
}

// MarshalBinary implements the [encoding.BinaryMarshaler] interface.
func (x *Xoshiro256PP) MarshalBinary() ([]byte, error) {
	return x.AppendBinary(make([]byte, 0, 40))
}

// UnmarshalBinary implements the [encoding.BinaryUnmarshaler] interface.
func (x *Xoshiro256PP) UnmarshalBinary(data []byte) error {
	if len(data) != 40 || string(data[:8]) != "xo256pp:" {
		return ErrUnmarshalWyRand
	}
	x.s0 = binary.LittleEndian.Uint64(data[8:])
	x.s1 = binary.LittleEndian.Uint64(data[16:])
	x.s2 = binary.LittleEndian.Uint64(data[24:])
	x.s3 = binary.LittleEndian.Uint64(data[32:])
	return nil
}
