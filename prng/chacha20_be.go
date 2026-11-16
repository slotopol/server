//go:build mips || mips64 || ppc64 || s390x || armbe || arm64be || ppc || s390 || sparc || sparc64

package prng

import (
	"math/bits"
	"time"
	"unsafe"
)

// refill fills the buffer with 512 new numbers in one call to XORKeyStream.
func (s *ChaCha20) refill() {
	if s.readed >= s.thres || (s.period > 0 && time.Since(s.last) >= s.period) {
		s.reseed()
	}

	clear(s.words[:])
	var byteView = unsafe.Slice((*byte)(unsafe.Pointer(&s.words[0])), bufferWords*8)
	s.cipher.XORKeyStream(byteView, byteView)

	// Reverse the byte order in each word from Big Endian to the
	// canonical Little Endian to preserve bitwise determinism.
	for i := range s.words {
		s.words[i] = bits.ReverseBytes64(s.words[i])
	}

	s.readed += bufferWords
	s.index = 0
}
