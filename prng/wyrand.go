package prng

import (
	cr "crypto/rand"
	"encoding/binary"
	"math/bits"
	mr "math/rand/v2"
	"sync/atomic"
)

// WyRand implements math/rand/v2.Source based on the ultra-fast WyRand algorithm.
// The structure occupies only 8 bytes and fits entirely within a CPU register.
//
// ATTENTION: Not thread-safe - create a separate instance for each goroutine.
// Usage sample (for Monte-Carlo simulations and tests):
//   var src = NewRandomWyRand()
//   var rng = mr.New(src)
//   fmt.Println("Stopping position:", rng.IntN(128))
type WyRand struct {
	state uint64
}

var _ mr.Source = (*WyRand)(nil)

// NewWyRand creates a generator with a given initial state (seed).
func NewWyRand(seed uint64) *WyRand {
	return &WyRand{state: seed}
}

// NewRandomWyRand initializes the generator with a random seed from "crypto/rand".
func NewRandomWyRand() *WyRand {
	var b [8]byte
	cr.Read(b[:])
	return &WyRand{state: binary.LittleEndian.Uint64(b[:])}
}

// GetWyRand returns fast WyRand random from given seed.
func GetWyRand(seed uint64) uint64 {
	seed += 0xa0761d6478bd642f
	var hi, lo = bits.Mul64(seed, seed^0xe7037ed1a0b428db)
	return hi ^ lo
}

// Seed resets the WyRand to behave the same way as NewWyRand(seed).
func (s *WyRand) Seed(seed uint64) {
	s.state = seed
}

// Uint64 generates the next random 64-bit number.
// The Go compiler automatically inlines this method into the calling code.
func (s *WyRand) Uint64() uint64 {
	s.state += 0xa0761d6478bd642f
	var hi, lo = bits.Mul64(s.state, s.state^0xe7037ed1a0b428db)
	return hi ^ lo
}

// AppendBinary implements the [encoding.BinaryAppender] interface.
func (s *WyRand) AppendBinary(b []byte) ([]byte, error) {
	b = append(b, "wyr:"...)
	b = binary.LittleEndian.AppendUint64(b, s.state)
	return b, nil
}

// MarshalBinary implements the [encoding.BinaryMarshaler] interface.
func (s *WyRand) MarshalBinary() ([]byte, error) {
	return s.AppendBinary(make([]byte, 0, 12))
}

// UnmarshalBinary implements the [encoding.BinaryUnmarshaler] interface.
func (s *WyRand) UnmarshalBinary(data []byte) error {
	if len(data) != 12 || string(data[:4]) != "wyr:" {
		return ErrUnmarshalWyRand
	}
	s.state = binary.LittleEndian.Uint64(data[4:])
	return nil
}

// WyRandSafe implements math/rand/v2.Source based on the ultra-fast WyRand algorithm.
// The structure occupies only 8 bytes and fits entirely within a CPU register.
//
// Algorithm is thread-safe, and oriented on cross-goroutines usage.
type WyRandSafe struct {
	state atomic.Uint64
}

var _ mr.Source = (*WyRandSafe)(nil)

// NewWyRandSafe creates a generator with a given initial state (seed).
func NewWyRandSafe(seed uint64) *WyRandSafe {
	var s WyRandSafe
	s.state.Store(seed)
	return &s
}

// NewRandomWyRandSafe initializes the generator with a random seed from "crypto/rand".
func NewRandomWyRandSafe() *WyRandSafe {
	var b [8]byte
	cr.Read(b[:])
	var s WyRandSafe
	s.state.Store(binary.LittleEndian.Uint64(b[:]))
	return &s
}

// Seed resets the WyRandSafe to behave the same way as NewWyRandSafe(seed).
func (s *WyRandSafe) Seed(seed uint64) {
	s.state.Store(seed)
}

// Uint64 generates the next random 64-bit number.
// This is thread-safe version of WyRand implementation.
func (s *WyRandSafe) Uint64() uint64 {
	var next = s.state.Add(0xa0761d6478bd642f)
	hi, lo := bits.Mul64(next, next^0xe7037ed1a0b428db)
	return hi ^ lo
}

// AppendBinary implements the [encoding.BinaryAppender] interface.
func (s *WyRandSafe) AppendBinary(b []byte) ([]byte, error) {
	b = append(b, "wyr:"...)
	b = binary.LittleEndian.AppendUint64(b, s.state.Load())
	return b, nil
}

// MarshalBinary implements the [encoding.BinaryMarshaler] interface.
func (s *WyRandSafe) MarshalBinary() ([]byte, error) {
	return s.AppendBinary(make([]byte, 0, 12))
}

// UnmarshalBinary implements the [encoding.BinaryUnmarshaler] interface.
func (s *WyRandSafe) UnmarshalBinary(data []byte) error {
	if len(data) != 12 || string(data[:4]) != "wyr:" {
		return ErrUnmarshalWyRand
	}
	s.state.Store(binary.LittleEndian.Uint64(data[4:]))
	return nil
}
