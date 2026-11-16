//go:build !(mips || mips64 || ppc64 || s390x || armbe || arm64be || ppc || s390 || sparc || sparc64)

package prng

import (
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

	s.readed += bufferWords
	s.index = 0
}
