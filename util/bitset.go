package util

import (
	"iter"
	"math/bits"
)

// Bitset64 is bitset on 64 bites.
type Bitset64 uint64

// MakeBitset creates bits set from slice of integer indexes.
func MakeBitset64(indexes ...int) Bitset64 {
	var bs Bitset64
	for _, n := range indexes {
		bs |= 1 << n
	}
	return bs
}

// MakeBitNum creates bits set with first num bits.
func MakeBitNum64(count, from int) Bitset64 {
	return (1<<count - 1) << from
}

// Num returns number of one bits in set.
func (bs Bitset64) Num() int {
	return bits.OnesCount64(uint64(bs))
}

// Bits iterates over ones in bitset.
func (bs Bitset64) Bits() iter.Seq[int] {
	return func(yield func(int) bool) {
		var bit Bitset64 = 1
		for pos := range 64 {
			if bs&bit > 0 && !yield(pos) {
				return
			}
			bit <<= 1
		}
	}
}

// Is checks that bit with given number is set.
func (bs Bitset64) Is(n int) bool {
	return bs&(1<<n) > 0
}

// Set bit with given number.
func (bs *Bitset64) Set(n int) *Bitset64 {
	*bs |= 1 << n
	return bs
}

// Res resets bit with given number.
func (bs *Bitset64) Res(n int) *Bitset64 {
	*bs &^= 1 << n
	return bs
}

// Toggle bit with given number.
func (bs *Bitset64) Toggle(n int) *Bitset64 {
	*bs ^= 1 << n
	return bs
}

// Sets first n bits.
func (bs *Bitset64) SetNum(count, from int) *Bitset64 {
	*bs = (1<<count - 1) << from
	return bs
}

func (bs *Bitset64) And(mask Bitset64) *Bitset64 {
	*bs &= mask
	return bs
}

func (bs *Bitset64) Or(mask Bitset64) *Bitset64 {
	*bs |= mask
	return bs
}

func (bs *Bitset64) AndNot(mask Bitset64) *Bitset64 {
	*bs &^= mask
	return bs
}

func (bs *Bitset64) Xor(mask Bitset64) *Bitset64 {
	*bs ^= mask
	return bs
}

func (bs Bitset64) IsZero() bool {
	return bs == 0
}

// Bitset128 is bitset on 128 bites.
type Bitset128 [2]uint64

// MakeBitset128 creates bits set from slice of integer indexes.
func MakeBitset128(indexes ...int) Bitset128 {
	var bs Bitset128
	for _, n := range indexes {
		bs[n/64] |= 1 << (n % 64)
	}
	return bs
}

// MakeBitNum128 creates bits set with first num bits.
func MakeBitNum128(count, from int) (bs Bitset128) {
	bs.SetNum(count, from)
	return
}

// Num returns number of one bits in set.
func (bs *Bitset128) Num() (count int) {
	for _, u := range bs {
		count += bits.OnesCount64(u)
	}
	return
}

// Bits iterates over ones in bitset.
func (bs *Bitset128) Bits() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i, u := range bs {
			var bit uint64 = 1
			for pos := range 64 {
				if u&bit > 0 && !yield(i*64+pos) {
					return
				}
				bit <<= 1
			}
		}
	}
}

// Is checks that bit with given number is set.
func (bs *Bitset128) Is(n int) bool {
	return bs[n/64]&(1<<(n%64)) > 0
}

// Set bit with given number.
func (bs *Bitset128) Set(n int) *Bitset128 {
	bs[n/64] |= 1 << (n % 64)
	return bs
}

// Res resets bit with given number.
func (bs *Bitset128) Res(n int) *Bitset128 {
	bs[n/64] &^= 1 << (n % 64)
	return bs
}

// Toggle bit with given number.
func (bs *Bitset128) Toggle(n int) *Bitset128 {
	bs[n/64] ^= 1 << (n % 64)
	return bs
}

// LShift implements left shift of bitset.
func (bs *Bitset128) LShift(count int) *Bitset128 {
	var c uint64
	for i, u := range bs {
		bs[i] = (u << count) | c
		c = u >> (64 - count)
	}
	return bs
}

// Sets first n bits.
func (bs *Bitset128) SetNum(count, from int) *Bitset128 {
	var i int
	for i = 0; i < count/64; i++ {
		bs[i] = 0xffffffffffffffff
	}
	bs[i] = 1<<(count%64) - 1
	if from > 0 {
		bs.LShift(from)
	}
	return bs
}

func (bs *Bitset128) And(mask Bitset128) *Bitset128 {
	for i, u := range mask {
		bs[i] &= u
	}
	return bs
}

func (bs *Bitset128) Or(mask Bitset128) *Bitset128 {
	for i, u := range mask {
		bs[i] |= u
	}
	return bs
}

func (bs *Bitset128) AndNot(mask Bitset128) *Bitset128 {
	for i, u := range mask {
		bs[i] &^= u
	}
	return bs
}

func (bs *Bitset128) Xor(mask Bitset128) *Bitset128 {
	for i, u := range mask {
		bs[i] ^= u
	}
	return bs
}

func (bs *Bitset128) IsZero() bool {
	for _, u := range bs {
		if u != 0 {
			return false
		}
	}
	return true
}
