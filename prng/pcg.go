package prng

import (
	cr "crypto/rand"
	"encoding/binary"
	mr "math/rand/v2"
)

type Rand = mr.Rand // alias to math/rand/v2.Rand

var New = mr.New // alias to math/rand/v2.New

func NewRandomPCG() *mr.PCG {
	var b [16]byte
	cr.Read(b[:])
	return mr.NewPCG(
		binary.LittleEndian.Uint64(b[:8]),
		binary.LittleEndian.Uint64(b[8:]))
}
