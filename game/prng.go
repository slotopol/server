package game

import (
	"context"
	cr "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"math/bits"
	mr "math/rand/v2"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/crypto/chacha20"
)

const (
	bufferWords     = 512                  // 512 * 8 = 4096 bytes (64 blocks of ChaCha20)
	reseedThreshold = 64 * 1024 * 1024 / 8 // default reseeding every 64 MB
	reseedPeriod    = 15 * time.Minute     // default reseeding every 15 minutes
	cyclingPeriod   = 50 * time.Millisecond
)

type ChaCha20Source struct {
	mux    sync.Mutex
	cipher *chacha20.Cipher
	words  [bufferWords]uint64
	index  int
	readed uint64 // readed uint64 words
	thres  uint64 // ressed threshold, in uint64
	period time.Duration
	last   time.Time
}

var _ mr.Source = (*ChaCha20Source)(nil)

// Option configures the generator's parameters.
type Option func(*ChaCha20Source)

// WithReseedThreshold option sets the threshold of readed uint64 values
// to call reseed after that, default value is set on 64MB.
//
// ATTENTION: after 256GB readed bytes without reseed
// ChaCha20 will panic (uint32 blocks counter overflow).
func WithReseedThreshold(t uint64) Option {
	return func(s *ChaCha20Source) {
		s.thres = t
	}
}

// WithReseedPeriod option sets the period to reseed ChaCha20,
// default value is set on 15 minutes.
func WithReseedPeriod(d time.Duration) Option {
	return func(s *ChaCha20Source) {
		s.period = d
	}
}

// NewChaCha20Source creates a new generator with entropy from "crypto/rand".
// Usage sample:
//   var ctx, cancel = context.WithCancel(context.Background())
//   defer cancel()
//
//   var src, err = NewChaCha20Source()
//   if err != nil {
//   	panic(err)
//   }
//   src.StartBackgroundCycling(ctx, 0)
//
//   var rng = mr.New(src) // *Rand from "math/rand/v2"
//   fmt.Println("Stopping position:", rng.IntN(128))
func NewChaCha20Source(opts ...Option) (*ChaCha20Source, error) {
	var s = &ChaCha20Source{
		index:  bufferWords,
		thres:  reseedThreshold,
		period: reseedPeriod,
	}
	for _, opt := range opts {
		opt(s)
	}
	if err := s.reseed(); err != nil {
		return nil, err
	}
	return s, nil
}

// reseed collects new entropy, mixes it with the current state
// using SHA-256, and recreates the cipher.
func (s *ChaCha20Source) reseed() error {
	var entropy [64]byte // 32 bytes of new entropy + 32 bytes of the current stream (to preserve state)

	if _, err := io.ReadFull(cr.Reader, entropy[:32]); err != nil {
		return fmt.Errorf("failed to read OS entropy: %w", err)
	}

	// If the cipher has already been initialized, we mix its output into the new entropy
	// (a NIST/GLI requirement).
	if s.cipher != nil {
		s.cipher.XORKeyStream(entropy[32:], entropy[32:])
	} else {
		if _, err := io.ReadFull(cr.Reader, entropy[32:]); err != nil {
			return fmt.Errorf("failed to read initial extra entropy: %w", err)
		}
	}

	// Entropy conditioning via SHA-256 to obtain a uniform 32-byte key
	var key = sha256.Sum256(entropy[:])

	var nonce [12]byte
	if _, err := io.ReadFull(cr.Reader, nonce[:]); err != nil {
		return fmt.Errorf("failed to read nonce: %w", err)
	}

	var cipher, err = chacha20.NewUnauthenticatedCipher(key[:], nonce[:])
	if err != nil {
		return fmt.Errorf("failed to init chacha20: %w", err)
	}

	s.cipher = cipher
	s.readed = 0
	s.last = time.Now()
	s.index = bufferWords // Flush the buffer
	return nil
}

// refill fills the buffer with 512 new numbers in one call to XORKeyStream.
func (s *ChaCha20Source) refill() {
	if s.readed >= s.thres || (s.period > 0 && time.Since(s.last) >= s.period) {
		if err := s.reseed(); err != nil {
			panic(fmt.Sprintf("critical RNG failure during reseed: %v", err))
		}
	}

	clear(s.words[:])
	var byteView = unsafe.Slice((*byte)(unsafe.Pointer(&s.words[0])), bufferWords*8)
	s.cipher.XORKeyStream(byteView, byteView)

	s.readed += bufferWords
	s.index = 0
}

// Uint64 returns a random 64-bit number.
func (s *ChaCha20Source) Uint64() uint64 {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.index >= bufferWords {
		s.refill()
	}
	var val = s.words[s.index]
	s.index++
	return val
}

// StartBackgroundCycling runs a background goroutine to shift
// the state in accordance with GLI-19 requirements.
func (s *ChaCha20Source) StartBackgroundCycling(ctx context.Context, d time.Duration) {
	if d == 0 {
		d = cyclingPeriod
	}
	go func() {
		var ticker = time.NewTicker(d)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = s.Uint64()
			}
		}
	}()
}

// WyRandSource implements math/rand/v2.Source based on the ultra-fast WyRand algorithm.
// The structure occupies only 8 bytes and fits entirely within a CPU register.
//
// ATTENTION: Not thread-safe - create a separate instance for each goroutine.
// Usage sample (for Monte-Carlo simulations and tests):
//   var src = NewRandomWyRandSource()
//   var rng = mr.New(src)
//   fmt.Println("Stopping position:", rng.IntN(128))
type WyRandSource struct {
	state uint64
}

var _ mr.Source = (*WyRandSource)(nil)

// NewWyRandSource creates a generator with a given initial state (seed).
func NewWyRandSource(seed uint64) *WyRandSource {
	return &WyRandSource{state: seed}
}

// NewRandomWyRandSource initializes the generator with a random seed from "crypto/rand".
func NewRandomWyRandSource() *WyRandSource {
	var b [8]byte
	if _, err := cr.Read(b[:]); err != nil {
		panic(err)
	}
	return &WyRandSource{state: binary.LittleEndian.Uint64(b[:])}
}

// Uint64 generates the next random 64-bit number.
// The Go compiler automatically inlines this method into the calling code.
func (s *WyRandSource) Uint64() uint64 {
	s.state += 0xa0761d6478bd642f
	var hi, lo = bits.Mul64(s.state, s.state^0xe7037ed1a0b428db)
	return hi ^ lo
}

// WyRandSafeSource implements math/rand/v2.Source based on the ultra-fast WyRand algorithm.
// The structure occupies only 8 bytes and fits entirely within a CPU register.
//
// Algorithm is thread-safe, and oriented on cross-goroutines usage.
type WyRandSafeSource struct {
	state uint64
}

var _ mr.Source = (*WyRandSafeSource)(nil)

// NewWyRandSafeSource creates a generator with a given initial state (seed).
func NewWyRandSafeSource(seed uint64) *WyRandSafeSource {
	return &WyRandSafeSource{state: seed}
}

// NewRandomWyRandSafeSource initializes the generator with a random seed from "crypto/rand".
func NewRandomWyRandSafeSource() *WyRandSafeSource {
	var b [8]byte
	if _, err := cr.Read(b[:]); err != nil {
		panic(err)
	}
	return &WyRandSafeSource{state: binary.LittleEndian.Uint64(b[:])}
}

// Uint64 generates the next random 64-bit number.
// This is thread-safe version of WyRand implementation.
func (s *WyRandSafeSource) Uint64() uint64 {
	var next = atomic.AddUint64(&s.state, 0xa0761d6478bd642f)
	hi, lo := bits.Mul64(next, next^0xe7037ed1a0b428db)
	return hi ^ lo
}
